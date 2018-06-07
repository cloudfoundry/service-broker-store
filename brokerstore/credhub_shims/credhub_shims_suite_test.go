package credhub_shims_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCredhubShims(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CredhubShims Suite")
}
