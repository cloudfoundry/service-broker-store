package credhub_shims_test

import (
	"code.cloudfoundry.org/service-broker-store/brokerstore/credhub_shims/credhub_fakes"

	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"code.cloudfoundry.org/service-broker-store/brokerstore/credhub_shims"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CredhubShims", func() {
	var (
		fakeCredhubAuthShim *credhub_fakes.FakeCredhubAuth
		fakeBuilder         auth.Builder
	)

	BeforeEach(func() {
		fakeCredhubAuthShim = &credhub_fakes.FakeCredhubAuth{}
		fakeBuilder = auth.Noop
		fakeCredhubAuthShim.UaaClientCredentialsReturns(fakeBuilder)
	})

	Describe("NewCredhubShim", func() {
		It("instantiates credhub with UAA auth", func() {
			_, err := credhub_shims.NewCredhubShim("http://some-url", "some-client-id", "some-client-secret", true, fakeCredhubAuthShim)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeCredhubAuthShim.UaaClientCredentialsCallCount()).To(Equal(1))
			clientID, clientSecret := fakeCredhubAuthShim.UaaClientCredentialsArgsForCall(0)
			Expect(clientID).To(Equal("some-client-id"))
			Expect(clientSecret).To(Equal("some-client-secret"))
		})
	})
})
