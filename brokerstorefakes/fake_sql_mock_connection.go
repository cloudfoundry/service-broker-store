package brokerstorefakes

import (
	"code.cloudfoundry.org/goshims/sqlshim"
	"code.cloudfoundry.org/lager"
)

type FakeSQLMockConnection struct {
	sqlshim.SqlDB
}

func (fake FakeSQLMockConnection) Connect(logger lager.Logger) error {
	return nil
}

func (fake FakeSQLMockConnection) Migrated() (bool, error) {
	return false, nil
}
