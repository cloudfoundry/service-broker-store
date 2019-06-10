package brokerstore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

type SqlStore struct {
	logger    lager.Logger
	StoreType string
	Database  SqlConnection
}

func NewSqlStore(logger lager.Logger, dbDriver, username, password, host, port, dbName, caCert string, skipHostnameValidation bool) (*SqlStore, error) {
	var err error
	var toDatabase SqlVariant
	switch dbDriver {
	case "mysql":
		toDatabase = NewMySqlVariant(username, password, host, port, dbName, caCert, skipHostnameValidation)
	case "postgres":
		toDatabase = NewPostgresVariant(username, password, host, port, dbName, caCert)
	default:
		err = fmt.Errorf("Unrecognized Driver: %s", dbDriver)
		logger.Error("db-driver-unrecognized", err)
		return nil, err
	}

	store, err := NewSqlStoreWithVariant(logger, toDatabase)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func NewSqlStoreWithVariant(logger lager.Logger, toDatabase SqlVariant) (*SqlStore, error) {
	database := NewSqlConnection(toDatabase)

	err := initialize(logger, database)

	if err != nil {
		logger.Error("sql-failed-to-initialize-database", err)
		return nil, err
	}

	return &SqlStore{
		logger:   logger,
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

func (s *SqlStore) Retire() error {
	_, err := s.Database.Exec("INSERT INTO service_instances (id, value) VALUES (?, ?)", "migrated-to-credhub", "true")
	if err != nil {
		return err
	}

	return nil
}

func (s *SqlStore) IsRetired() (bool, error) {
	var id, value string

	if result := s.Database.QueryRow("SELECT id, value FROM service_instances WHERE id = ?", "migrated-to-credhub").Scan(&id, &value); result == nil {
		return true, nil
	} else if result == sql.ErrNoRows {
		return false, nil
	} else {
		return false, result
	}
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
	logger := s.logger.Session("create-instance-details")
	logger.Info("start")
	defer logger.Info("end")

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
	logger := s.logger.Session("retrieve-instance-details")
	logger.Info("start")
	defer logger.Info("end")

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
	logger := s.logger.Session("retrieve-binding-details")
	logger.Info("start")
	defer logger.Info("end")

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

func (s *SqlStore) RetrieveAllInstanceDetails() (map[string]ServiceInstance, error) {
	logger := s.logger.Session("retrieve-all-instance-details")
	logger.Info("start")
	defer logger.Info("end")

	serviceInstances := map[string]ServiceInstance{}

	rows, err := s.Database.Query("SELECT id, value FROM service_instances")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var serviceInstance ServiceInstance
		var id string
		var jsonValue []byte
		err = rows.Scan(&id, &jsonValue)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonValue, &serviceInstance)
		if err != nil {
			return nil, err
		}
		serviceInstances[id] = serviceInstance
	}

	return serviceInstances, nil
}

func (s *SqlStore) RetrieveAllBindingDetails() (map[string]brokerapi.BindDetails, error) {
	logger := s.logger.Session("retrieve-all-binding-details")
	logger.Info("start")
	defer logger.Info("end")

	bindingDetails := map[string]brokerapi.BindDetails{}

	rows, err := s.Database.Query("SELECT id, value FROM service_bindings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var bindDetails brokerapi.BindDetails
		var id string
		var jsonValue []byte
		err = rows.Scan(&id, &jsonValue)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonValue, &bindDetails)
		if err != nil {
			return nil, err
		}
		bindingDetails[id] = bindDetails
	}

	return bindingDetails, nil
}

func (s *SqlStore) CreateBindingDetails(id string, details brokerapi.BindDetails) error {
	logger := s.logger.Session("create-binding-details")
	logger.Info("start")
	defer logger.Info("end")

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
	logger := s.logger.Session("delete-instance-details")
	logger.Info("start")
	defer logger.Info("end")

	_, err := s.Database.Exec("DELETE FROM service_instances WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlStore) DeleteBindingDetails(id string) error {
	logger := s.logger.Session("delete-binding-details")
	logger.Info("start")
	defer logger.Info("end")

	_, err := s.Database.Exec("DELETE FROM service_bindings WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlStore) IsInstanceConflict(id string, details ServiceInstance) bool {
	return isInstanceConflict(s, id, details)
}

func (s *SqlStore) IsBindingConflict(id string, details brokerapi.BindDetails) bool {
	return isBindingConflict(s, id, details)
}
