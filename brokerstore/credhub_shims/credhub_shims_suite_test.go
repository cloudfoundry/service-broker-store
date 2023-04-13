package credhub_shims_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCredhubShims(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CredhubShims Suite")
}
