package service_broker_store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestServiceBrokerStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ServiceBrokerStore Suite")
}
