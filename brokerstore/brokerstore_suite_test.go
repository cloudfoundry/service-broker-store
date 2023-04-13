package brokerstore_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBrokerstore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Brokerstore Suite")
}
