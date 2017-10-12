package brokerstore

import (
	"fmt"
	"strings"

	"crypto/x509"

	"code.cloudfoundry.org/goshims/ioutilshim"
	"code.cloudfoundry.org/goshims/osshim"
	"code.cloudfoundry.org/goshims/sqlshim"
	"code.cloudfoundry.org/lager"
)

type postgresVariant struct {
	sql                sqlshim.Sql
	ioutil             ioutilshim.Ioutil
	os                 osshim.Os
	dbConnectionString string
	caCert             string
	dbName             string
}

func NewPostgresVariant(username, password, host, port, dbName, caCert string) SqlVariant {
	return NewPostgresVariantWithShims(username, password, host, port, dbName, caCert, &sqlshim.SqlShim{}, &ioutilshim.IoutilShim{}, &osshim.OsShim{})
}

func NewPostgresVariantWithShims(username, password, host, port, dbName, caCert string, sql sqlshim.Sql, ioutil ioutilshim.Ioutil, os osshim.Os) SqlVariant {
	return &postgresVariant{
		sql:                sql,
		os:                 os,
		ioutil:             ioutil,
		dbConnectionString: fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName),
		caCert:             caCert,
		dbName:             dbName,
	}
}

func (c *postgresVariant) Connect(logger lager.Logger) (sqlshim.SqlDB, error) {
	logger = logger.Session("postgres-connection-connect")
	logger.Info("start")
	defer logger.Info("end")

	if c.caCert == "" {
		c.dbConnectionString = fmt.Sprintf("%s?sslmode=disable", c.dbConnectionString)
	} else {
		certBytes := []byte(c.caCert)

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(certBytes); !ok {
			err := fmt.Errorf("Invalid CA Cert for %s", c.dbName)
			logger.Error("failed-to-parse-sql-ca", err)
			return nil, err
		}

		tmpFile, err := c.ioutil.TempFile("", "postgress-ca-cert")
		if err != nil {
			logger.Error("tempfile-create-failed", err)
			return nil, err
		}

		if _, err := tmpFile.Write([]byte(c.caCert)); err != nil {
			logger.Error("tempfile-write-failed", err)
			return nil, err
		}
		if err := tmpFile.Close(); err != nil {
			logger.Error("tempfile-close-failed", err)
			return nil, err
		}

		c.caCert = tmpFile.Name()
		c.dbConnectionString = fmt.Sprintf("%s?sslmode=verify-ca&sslrootcert=%s", c.dbConnectionString, c.caCert)
	}

	sqlDB, err := c.sql.Open("postgres", c.dbConnectionString)
	return sqlDB, err
}

func (c *postgresVariant) Flavorify(query string) string {
	strParts := strings.Split(query, "?")
	for i := 1; i < len(strParts); i++ {
		strParts[i-1] = fmt.Sprintf("%s$%d", strParts[i-1], i)
	}
	return strings.Join(strParts, "")
}

func (c *postgresVariant) Close() error {
	if c.caCert != "" {
		return c.os.Remove(c.caCert)
	}
	return nil
}
