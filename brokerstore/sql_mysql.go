package brokerstore

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"strings"
	"time"

	"code.cloudfoundry.org/goshims/mysqlshim"
	"code.cloudfoundry.org/goshims/sqlshim"
	"code.cloudfoundry.org/lager"
)

type mysqlVariant struct {
	sql                    sqlshim.Sql
	mysql                  mysqlshim.MySQL
	dbConnectionString     string
	caCert                 string
	dbName                 string
	skipHostnameValidation bool
}

func NewMySqlVariant(username, password, host, port, dbName, caCert string, skipHostnameValidation bool) SqlVariant {
	return NewMySqlVariantWithSqlObject(username, password, host, port, dbName, caCert, skipHostnameValidation, &sqlshim.SqlShim{}, &mysqlshim.MySQLShim{})
}

func NewMySqlVariantWithSqlObject(username, password, host, port, dbName, caCert string, skipHostnameValidation bool, sql sqlshim.Sql, mysql mysqlshim.MySQL) SqlVariant {
	return &mysqlVariant{
		sql:                    sql,
		mysql:                  mysql,
		dbConnectionString:     fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName),
		caCert:                 caCert,
		dbName:                 dbName,
		skipHostnameValidation: skipHostnameValidation,
	}
}

func (c *mysqlVariant) Connect(logger lager.Logger) (sqlshim.SqlDB, error) {
	logger = logger.Session("mysql-connection-connect")
	logger.Info("start")
	defer logger.Info("end")

	if c.caCert != "" {
		cfg, err := c.mysql.ParseDSN(c.dbConnectionString)
		if err != nil {
			logger.Fatal("invalid-db-connection-string", err, lager.Data{"connection-string": c.dbConnectionString})
		}

		logger.Debug("secure-mysql")
		// parse off any leading whitespace from the ca certificate
		cert := strings.Replace(c.caCert, "  ", "", -1)
		cert = strings.Replace(cert, "\n ", "\n", -1)
		cert = strings.Replace(cert, "\t", "", -1)
		cert = strings.TrimLeft(cert, " \t")

		certBytes := []byte(cert)

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(certBytes); !ok {
			err := fmt.Errorf("Invalid CA Cert for %s", c.dbName)
			logger.Error("failed-to-parse-sql-ca", err)
			return nil, err

		}

		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			RootCAs:            caCertPool,
		}

		if c.skipHostnameValidation {
			tlsConfig.InsecureSkipVerify = true

			tlsConfig.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				return VerifyCertificatesIgnoreHostname(rawCerts, caCertPool)
			}
		}

		ourKey := "nfs-tls"
		c.mysql.RegisterTLSConfig(ourKey, tlsConfig)
		cfg.TLSConfig = ourKey
		cfg.Timeout = 10 * time.Minute
		cfg.ReadTimeout = 10 * time.Minute
		cfg.WriteTimeout = 10 * time.Minute
		c.dbConnectionString = cfg.FormatDSN()
	}

	sqlDB, err := c.sql.Open("mysql", c.dbConnectionString)
	return sqlDB, err
}

func (c *mysqlVariant) Flavorify(query string) string {
	return query
}

func (c *mysqlVariant) Close() error {
	return nil
}

func VerifyCertificatesIgnoreHostname(rawCerts [][]byte, caCertPool *x509.CertPool) error {
	certs := make([]*x509.Certificate, len(rawCerts))

	for i, asn1Data := range rawCerts {
		cert, err := x509.ParseCertificate(asn1Data)
		if err != nil {
			return errors.New("tls: failed to parse certificate from server: " + err.Error())
		}

		certs[i] = cert
	}

	opts := x509.VerifyOptions{
		Roots:         caCertPool,
		CurrentTime:   time.Now(),
		Intermediates: x509.NewCertPool(),
	}

	for i, cert := range certs {
		if i == 0 {
			continue
		}

		opts.Intermediates.AddCert(cert)
	}

	_, err := certs[0].Verify(opts)
	return err
}
