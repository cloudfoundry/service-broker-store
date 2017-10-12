package brokerstore

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"code.cloudfoundry.org/goshims/sqlshim"
	"code.cloudfoundry.org/lager"
)

//go:generate counterfeiter -o ../brokerstorefakes/fake_sql_variant.go . SqlVariant
type SqlVariant interface {
	Connect(logger lager.Logger) (sqlshim.SqlDB, error)
	Flavorify(query string) string
	Close() error
}

//go:generate counterfeiter -o ../brokerstorefakes/fake_sql_connection.go . SqlConnection
type SqlConnection interface {
	Connect(logger lager.Logger) error
	sqlshim.SqlDB
}

type sqlConnection struct {
	sqlDB sqlshim.SqlDB
	leaf  SqlVariant
}

func NewSqlConnection(variant SqlVariant) SqlConnection {
	if variant == nil {
		panic("variant cannot be nil")
	}
	return &sqlConnection{
		leaf: variant,
	}
}

func (c *sqlConnection) flavorify(query string) string {
	return c.leaf.Flavorify(query)
}

func (c *sqlConnection) Connect(logger lager.Logger) error {
	sqlDB, err := c.leaf.Connect(logger)
	if err != nil {
		return err
	}

	c.sqlDB = sqlDB

	err = c.Ping()
	return err
}

func (c *sqlConnection) Ping() error {
	return c.sqlDB.Ping()
}
func (c *sqlConnection) Close() error {
	defer c.leaf.Close()
	return c.sqlDB.Close()
}
func (c *sqlConnection) SetMaxIdleConns(n int) {
	c.sqlDB.SetMaxIdleConns(n)
}
func (c *sqlConnection) SetMaxOpenConns(n int) {
	c.sqlDB.SetMaxOpenConns(n)
}
func (c *sqlConnection) SetConnMaxLifetime(d time.Duration) {
	c.sqlDB.SetConnMaxLifetime(d)
}
func (c *sqlConnection) Stats() sql.DBStats {
	return c.sqlDB.Stats()
}
func (c *sqlConnection) Prepare(query string) (*sql.Stmt, error) {
	return c.sqlDB.Prepare(c.flavorify(query))
}
func (c *sqlConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.sqlDB.Exec(c.flavorify(query), args...)
}
func (c *sqlConnection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.sqlDB.Query(c.flavorify(query), args...)
}
func (c *sqlConnection) QueryRow(query string, args ...interface{}) *sql.Row {
	return c.sqlDB.QueryRow(c.flavorify(query), args...)
}
func (c *sqlConnection) Begin() (*sql.Tx, error) {
	return c.sqlDB.Begin()
}
func (c *sqlConnection) Driver() driver.Driver {
	return c.sqlDB.Driver()
}
