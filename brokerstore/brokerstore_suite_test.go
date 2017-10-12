package brokerstore_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBrokerstore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Brokerstore Suite")
}
