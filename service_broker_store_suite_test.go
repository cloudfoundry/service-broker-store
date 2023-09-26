package service_broker_store_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestServiceBrokerStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ServiceBrokerStore Suite")
}

