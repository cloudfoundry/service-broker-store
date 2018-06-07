package brokerstore

import (
	"fmt"

	"database/sql"

	"encoding/json"
	"errors"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

type SqlStore struct {
	StoreType string
	Database  SqlConnection
}

func NewSqlStore(logger lager.Logger, dbDriver, username, password, host, port, dbName, caCert string) (Store, error) {

	var err error
	var toDatabase SqlVariant
	switch dbDriver {
	case "mysql":
		toDatabase = NewMySqlVariant(username, password, host, port, dbName, caCert)
	case "postgres":
		toDatabase = NewPostgresVariant(username, password, host, port, dbName, caCert)
	default:
		err = fmt.Errorf("Unrecognized Driver: %s", dbDriver)
		logger.Error("db-driver-unrecognized", err)
		return nil, err
	}
	return NewSqlStoreWithVariant(logger, toDatabase)
}

func NewSqlStoreWithVariant(logger lager.Logger, toDatabase SqlVariant) (Store, error) {
	database := NewSqlConnection(toDatabase)

	err := initialize(logger, database)

	if err != nil {
		logger.Error("sql-failed-to-initialize-database", err)
		return nil, err
	}

	return NewSqlStoreWithDatabase(logger, database)
}

func NewSqlStoreWithDatabase(logger lager.Logger, database SqlConnection) (Store, error) {
	return &SqlStore{
		Database: database,
	}, nil
}

func initialize(logger lager.Logger, db SqlConnection) error {
	logger = logger.Session("initialize-database")
	logger.Info("start")
	defer logger.Info("end")

	var err error
	err = db.Connect(logger)
	if err != nil {
		logger.Error("sql-failed-to-connect", err)
		return err
	}

	// TODO: uniquify table names?
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS service_instances(
				id VARCHAR(255) PRIMARY KEY,
				value VARCHAR(4096)
			)
		`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS service_bindings(
				id VARCHAR(255) PRIMARY KEY,
				value VARCHAR(4096)
			)
		`)
	return err
}

func (s *SqlStore) Restore(logger lager.Logger) error {
	return nil
}

func (s *SqlStore) Save(logger lager.Logger) error {
	return nil
}

func (s *SqlStore) Cleanup() error {
	return nil
}

func (s *SqlStore) CreateInstanceDetails(id string, details ServiceInstance) error {
	jsonData, err := json.Marshal(details)
	if err != nil {
		return err
	}

	if passwordCheck(jsonData) {
		return errors.New("passwords are not allowed in service instance configuration")
	}

	_, err = s.Database.Exec("INSERT INTO service_instances (id, value) VALUES (?, ?)", id, jsonData)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlStore) RetrieveInstanceDetails(id string) (ServiceInstance, error) {
	var serviceID string
	var value []byte
	var serviceInstance ServiceInstance
	if err := s.Database.QueryRow("SELECT id, value FROM service_instances WHERE id = ?", id).Scan(&serviceID, &value); err == nil {
		err = json.Unmarshal(value, &serviceInstance)
		if err != nil {
			return ServiceInstance{}, err
		}
		return serviceInstance, nil
	} else if err == sql.ErrNoRows {
		return ServiceInstance{}, brokerapi.ErrInstanceDoesNotExist
	} else {
		return ServiceInstance{}, err
	}
}

func (s *SqlStore) RetrieveBindingDetails(id string) (brokerapi.BindDetails, error) {
	var bindingID string
	var value []byte
	bindDetails := brokerapi.BindDetails{}
	if err := s.Database.QueryRow("SELECT id, value FROM service_bindings WHERE id = ?", id).Scan(&bindingID, &value); err == nil {
		err = json.Unmarshal(value, &bindDetails)
		if err != nil {
			return brokerapi.BindDetails{}, err
		}
		return bindDetails, nil
	} else if err == sql.ErrNoRows {
		return brokerapi.BindDetails{}, brokerapi.ErrInstanceDoesNotExist
	} else {
		return brokerapi.BindDetails{}, err
	}
}

func (s *SqlStore) CreateBindingDetails(id string, details brokerapi.BindDetails) error {
	storeDetails, err := redactBindingDetails(details)

	jsonData, err := json.Marshal(storeDetails)
	if err != nil {
		return err
	}
	_, err = s.Database.Exec("INSERT INTO service_bindings (id, value) VALUES (?, ?)", id, jsonData)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlStore) DeleteInstanceDetails(id string) error {
	_, err := s.Database.Exec("DELETE FROM service_instances WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlStore) DeleteBindingDetails(id string) error {
	_, err := s.Database.Exec("DELETE FROM service_bindings WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlStore) keyValueInTable(logger lager.Logger, key, value, table string) (error, bool) {
	var queriedServiceID string
	query := fmt.Sprintf(`SELECT %s.%s FROM %s WHERE %s.%s = ?`, table, key, table, table, key)
	row := s.Database.QueryRow(query, value)
	if row == nil {
		err := fmt.Errorf("Row error!")
		logger.Error("failed-query", err)
		return err, true
	}
	err := row.Scan(&queriedServiceID)
	if err == nil {
		return nil, true
	} else if err == sql.ErrNoRows {
		return nil, false
	}

	logger.Debug("failed-query", lager.Data{"Query": query})
	logger.Error("failed-query", err)
	return err, true
}

func (s *SqlStore) IsInstanceConflict(id string, details ServiceInstance) bool {
	return isInstanceConflict(s, id, details)
}

func (s *SqlStore) IsBindingConflict(id string, details brokerapi.BindDetails) bool {
	return isBindingConflict(s, id, details)
}
