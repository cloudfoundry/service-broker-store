package brokerstore

import (
	"fmt"

	"crypto/tls"
	"crypto/x509"
	"time"

	"code.cloudfoundry.org/goshims/sqlshim"
	"code.cloudfoundry.org/lager"
	"github.com/go-sql-driver/mysql"
)

type mysqlVariant struct {
	sql                sqlshim.Sql
	dbConnectionString string
	caCert             string
	dbName             string
}

func NewMySqlVariant(username, password, host, port, dbName, caCert string) SqlVariant {
	return NewMySqlVariantWithSqlObject(username, password, host, port, dbName, caCert, &sqlshim.SqlShim{})
}

func NewMySqlVariantWithSqlObject(username, password, host, port, dbName, caCert string, sql sqlshim.Sql) SqlVariant {
	return &mysqlVariant{
		sql:                sql,
		dbConnectionString: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName),
		caCert:             caCert,
		dbName:             dbName,
	}
}

func (c *mysqlVariant) Connect(logger lager.Logger) (sqlshim.SqlDB, error) {
	logger = logger.Session("mysql-connection-connect")
	logger.Info("start")
	defer logger.Info("end")

	if c.caCert != "" {
		cfg, err := mysql.ParseDSN(c.dbConnectionString)
		if err != nil {
			logger.Fatal("invalid-db-connection-string", err, lager.Data{"connection-string": c.dbConnectionString})
		}

		logger.Debug("secure-mysql")
		certBytes := []byte(c.caCert)

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
		ourKey := "nfs-tls"
		mysql.RegisterTLSConfig(ourKey, tlsConfig)
		cfg.TLSConfig = ourKey
		cfg.Timeout = 10 * time.Minute
		cfg.ReadTimeout = 10 * time.Minute
		cfg.WriteTimeout = 10 * time.Minute
		c.dbConnectionString = cfg.FormatDSN()
	}

	logger.Info("db-string", lager.Data{"value": c.dbConnectionString})
	sqlDB, err := c.sql.Open("mysql", c.dbConnectionString)
	return sqlDB, err
}

func (c *mysqlVariant) Flavorify(query string) string {
	return query
}

func (c *mysqlVariant) Close() error {
	return nil
}
