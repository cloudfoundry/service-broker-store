package brokerstore_test

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"code.cloudfoundry.org/goshims/mysqlshim/mysql_fake"
	"code.cloudfoundry.org/goshims/sqlshim/sql_fake"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"code.cloudfoundry.org/service-broker-store/brokerstore"
	"github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MysqlVariant", func() {
	var (
		logger                 lager.Logger
		fakeSql                *sql_fake.FakeSql
		fakeMySQL              *mysql_fake.FakeMySQL
		err                    error
		database               brokerstore.SqlVariant
		skipHostnameValidation bool

		cert       string
		caCertPool *x509.CertPool
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("mysql-variant-test")
		skipHostnameValidation = false

		fakeSql = &sql_fake.FakeSql{}
		fakeMySQL = &mysql_fake.FakeMySQL{}
		fakeMySQL.ParseDSNStub = mysql.ParseDSN
	})

	JustBeforeEach(func() {
		database = brokerstore.NewMySqlVariantWithSqlObject("username", "password", "host", "port", "dbName", cert, skipHostnameValidation, fakeSql, fakeMySQL)
	})

	Describe(".Connect", func() {
		JustBeforeEach(func() {
			_, err = database.Connect(logger)
		})

		Context("when ca cert specified", func() {
			BeforeEach(func() {
				cert = exampleCaCert

				caCertPool = x509.NewCertPool()
				ok := caCertPool.AppendCertsFromPEM([]byte(cert))
				Expect(ok).To(BeTrue())
			})

			It("open call has correctly formed connection string", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeSql.OpenCallCount()).To(Equal(1))
				dbType, connectionString := fakeSql.OpenArgsForCall(0)
				Expect(dbType).To(Equal("mysql"))
				Expect(connectionString).To(Equal("username:password@tcp(host:port)/dbName?readTimeout=10m0s\u0026timeout=10m0s\u0026tls=nfs-tls\u0026writeTimeout=10m0s"))
			})

			It("registers a TLSConfig", func() {
				Expect(fakeMySQL.RegisterTLSConfigCallCount()).To(Equal(1))

				passedTLSConfigName, passedTLSConfig := fakeMySQL.RegisterTLSConfigArgsForCall(0)
				Expect(passedTLSConfigName).To(Equal("nfs-tls"))
				Expect(passedTLSConfig).To(Equal(&tls.Config{
					InsecureSkipVerify: false,
					RootCAs:            caCertPool,
				}))
			})
		})

		Context("when ca cert has leading whitespace", func() {
			BeforeEach(func() {
				cert = spacyCaCert
			})

			It("open call has correctly formed connection string", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeSql.OpenCallCount()).To(Equal(1))
				dbType, connectionString := fakeSql.OpenArgsForCall(0)
				Expect(dbType).To(Equal("mysql"))
				Expect(connectionString).To(Equal("username:password@tcp(host:port)/dbName?readTimeout=10m0s\u0026timeout=10m0s\u0026tls=nfs-tls\u0026writeTimeout=10m0s"))
			})
		})

		Context("when ca cert specified is invalid", func() {
			BeforeEach(func() {
				cert = "invalid"
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the skipHostnameValidation flag is true", func() {
			BeforeEach(func() {
				cert = exampleCaCert

				caCertPool = x509.NewCertPool()
				ok := caCertPool.AppendCertsFromPEM([]byte(cert))
				Expect(ok).To(BeTrue())

				skipHostnameValidation = true
			})

			It("open call has correctly formed connection string", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeSql.OpenCallCount()).To(Equal(1))
				dbType, connectionString := fakeSql.OpenArgsForCall(0)
				Expect(dbType).To(Equal("mysql"))
				Expect(connectionString).To(Equal("username:password@tcp(host:port)/dbName?readTimeout=10m0s\u0026timeout=10m0s\u0026tls=nfs-tls\u0026writeTimeout=10m0s"))
			})

			It("registers a TLSConfig with a custom cert verification function", func() {
				Expect(fakeMySQL.RegisterTLSConfigCallCount()).To(Equal(1))

				passedTLSConfigName, passedTLSConfig := fakeMySQL.RegisterTLSConfigArgsForCall(0)
				Expect(passedTLSConfigName).To(Equal("nfs-tls"))
				Expect(passedTLSConfig.InsecureSkipVerify).To(BeTrue())
				Expect(passedTLSConfig.RootCAs).To(Equal(caCertPool))
				// impossible to assert VerifyPeerCertificate is set to a specfic function
				Expect(passedTLSConfig.VerifyPeerCertificate).NotTo(BeNil())
			})
		})

		Context("when no ca cert specified", func() {
			BeforeEach(func() {
				cert = ""
			})

			Context("when it can connect to a valid database", func() {
				BeforeEach(func() {
					fakeSql.OpenReturns(nil, nil)
				})

				It("reports no error", func() {
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeSql.OpenCallCount()).To(Equal(1))
					dbType, connectionString := fakeSql.OpenArgsForCall(fakeSql.OpenCallCount() - 1)
					Expect(dbType).To(Equal("mysql"))
					Expect(connectionString).To(Equal("username:password@tcp(host:port)/dbName"))
				})
			})

			Context("when it cannot connect to a valid database", func() {
				BeforeEach(func() {
					fakeSql = &sql_fake.FakeSql{}
					fakeSql.OpenReturns(nil, errors.New("something wrong"))
				})

				It("reports error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})

	Describe(".Flavorify", func() {
		It("should return unaltered query", func() {
			query := `INSERT INTO service_instances (id, value) VALUES (?, ?)`
			Expect(database.Flavorify(query)).To(Equal(query))
		})
	})

	Describe(".Close", func() {
		It("doesn't fail", func() {
			err := database.Close()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("VerifyCertificatesIgnoreHostname", func() {
		BeforeEach(func() {
			caCertPool = x509.NewCertPool()
			ok := caCertPool.AppendCertsFromPEM([]byte(exampleCaCert))
			Expect(ok).To(BeTrue())
		})

		It("verifies that provided certificates are valid", func() {
			block, _ := pem.Decode([]byte(exampleCaCert))
			err := brokerstore.VerifyCertificatesIgnoreHostname([][]byte{
				block.Bytes,
			}, caCertPool)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when raw certs are not parsable", func() {
			It("returns an error", func() {
				err := brokerstore.VerifyCertificatesIgnoreHostname([][]byte{
					[]byte("foo"),
					[]byte("bar"),
				}, nil)
				Expect(err.Error()).To(ContainSubstring("tls: failed to parse certificate from server: asn1: structure error: tags don't match"))
			})
		})

		Context("when verifying an expired cert", func() {
			It("returns an error", func() {
				block, _ := pem.Decode([]byte(expiredCert))
				err := brokerstore.VerifyCertificatesIgnoreHostname([][]byte{
					block.Bytes,
				}, caCertPool)
				Expect(err.Error()).To(ContainSubstring("x509: certificate has expired or is not yet valid"))
			})
		})
	})
})
