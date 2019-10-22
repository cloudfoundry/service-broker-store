package brokerstore_test

import (
	"code.cloudfoundry.org/lager/lagertest"
	"code.cloudfoundry.org/service-broker-store/brokerstore"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Store", func() {

	Context("#NewStore", func() {

		Context("when no credhub credentials are supplied", func() {

			It("should log a fatal error", func() {
				logger := lagertest.NewTestLogger("broker-store")
				Expect(func() {
					brokerstore.NewStore(logger, "","","","","","")
				}).Should(Panic())

				Expect(logger.Buffer()).Should(gbytes.Say("Invalid brokerstore configuration"))
			})
		})
	})

})
