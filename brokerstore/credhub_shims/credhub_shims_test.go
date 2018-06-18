package credhub_shims_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"

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
		caCert              string
	)

	BeforeEach(func() {
		fakeCredhubAuthShim = &credhub_fakes.FakeCredhubAuth{}
		fakeBuilder = auth.Noop
		fakeCredhubAuthShim.UaaClientCredentialsReturns(fakeBuilder)
		caCert = generateCaCert()
	})

	Describe("NewCredhubShim", func() {
		It("instantiates credhub with UAA auth", func() {
			_, err := credhub_shims.NewCredhubShim("http://some-url", caCert, "some-client-id", "some-client-secret", fakeCredhubAuthShim)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeCredhubAuthShim.UaaClientCredentialsCallCount()).To(Equal(1))
			clientID, clientSecret := fakeCredhubAuthShim.UaaClientCredentialsArgsForCall(0)
			Expect(clientID).To(Equal("some-client-id"))
			Expect(clientSecret).To(Equal("some-client-secret"))
		})

		Context("when CA cert is not provided", func() {
			BeforeEach(func() {
				caCert = generateCaCert()
			})

			It("instantiates credhub without CA cert", func() {
				_, err := credhub_shims.NewCredhubShim("http://some-url", "", "some-client-id", "some-client-secret", fakeCredhubAuthShim)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeCredhubAuthShim.UaaClientCredentialsCallCount()).To(Equal(1))
				clientID, clientSecret := fakeCredhubAuthShim.UaaClientCredentialsArgsForCall(0)
				Expect(clientID).To(Equal("some-client-id"))
				Expect(clientSecret).To(Equal("some-client-secret"))
			})
		})
	})
})

func generateCaCert() string {
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	Expect(err).NotTo(HaveOccurred())

	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	Expect(err).NotTo(HaveOccurred())
	certOut := new(bytes.Buffer)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	return certOut.String()
}
