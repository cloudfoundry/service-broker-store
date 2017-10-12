package brokerstore_test

import (
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/cloudfoundry-incubator/service-broker-store/brokerstore"

	"errors"

	"code.cloudfoundry.org/goshims/ioutilshim/ioutil_fake"
	"code.cloudfoundry.org/goshims/osshim/os_fake"
	"code.cloudfoundry.org/goshims/sqlshim/sql_fake"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PostgresVariant", func() {
	var (
		logger     *lagertest.TestLogger
		fakeSql    *sql_fake.FakeSql
		fakeIoUtil *ioutil_fake.FakeIoutil
		fakeOs     *os_fake.FakeOs
		fakeFile   *os_fake.FakeFile
		database   brokerstore.SqlVariant

		cert string
	)

	Describe(".Connect", func() {
		var (
			err error
		)

		JustBeforeEach(func() {
			database = brokerstore.NewPostgresVariantWithShims("username", "password", "host", "port", "dbName", cert, fakeSql, fakeIoUtil, fakeOs)
			_, err = database.Connect(logger)
		})

		Context("given no ca cert", func() {
			BeforeEach(func() {
				logger = lagertest.NewTestLogger("portgress-variant-test")

				fakeSql = &sql_fake.FakeSql{}
				fakeIoUtil = &ioutil_fake.FakeIoutil{}

				cert = ""
			})

			Context("when connecting to valid database", func() {
				BeforeEach(func() {
					fakeSql.OpenReturns(nil, nil)
				})

				It("connects successfully", func() {
					Expect(err).NotTo(HaveOccurred())

					dbType, connectionString := fakeSql.OpenArgsForCall(0)
					Expect(dbType).To(Equal("postgres"))
					Expect(connectionString).To(Equal("postgres://username:password@host:port/dbName?sslmode=disable"))
				})
			})

			Context("when connecting to invalid database", func() {
				BeforeEach(func() {
					fakeSql.OpenReturns(nil, errors.New("something wrong"))
				})

				JustBeforeEach(func() {
					_, err = database.Connect(logger)
				})

				It("reports error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("given a valid ca cert", func() {

			BeforeEach(func() {
				logger = lagertest.NewTestLogger("portgress-variant-test")

				fakeSql = &sql_fake.FakeSql{}
				fakeIoUtil = &ioutil_fake.FakeIoutil{}
				fakeFile = &os_fake.FakeFile{}

				cert = exampleCaCert
			})

			Context("and can create a temp file", func() {
				BeforeEach(func() {
					fakeIoUtil.TempFileReturns(fakeFile, nil)
				})

				Context("and can write to temp file", func() {
					BeforeEach(func() {
						fakeFile.WriteAtReturns(0, nil)
					})

					Context("and can close temp file", func() {
						BeforeEach(func() {
							fakeFile.CloseReturns(nil)
						})

						Context("when a connection is attempted", func() {
							BeforeEach(func() {
								fakeFile.NameReturns("/a/temp.file")
							})

							It("should have a correctly formed connection string", func() {
								_, connectionString := fakeSql.OpenArgsForCall(0)
								Expect(connectionString).To(Equal("postgres://username:password@host:port/dbName?sslmode=verify-ca&sslrootcert=/a/temp.file"))
							})
						})
					})

					Context("when the close fails", func() {
						BeforeEach(func() {
							fakeFile.CloseReturns(errors.New("badness"))
						})
						It("should return an error", func() {
							Expect(err).To(HaveOccurred())
						})
					})
				})

				Context("when the write fails", func() {
					BeforeEach(func() {
						fakeFile.WriteReturns(1, errors.New("badness"))
					})
					It("should return an error", func() {
						Expect(err).To(HaveOccurred())
					})
				})
			})

			Context("when the create file fails", func() {
				BeforeEach(func() {
					fakeIoUtil.TempFileReturns(nil, errors.New("badness"))
				})
				It("return an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("given an invalid cert", func() {

			BeforeEach(func() {
				logger = lagertest.NewTestLogger("portgress-variant-test")

				fakeSql = &sql_fake.FakeSql{}
				fakeIoUtil = &ioutil_fake.FakeIoutil{}
				fakeFile = &os_fake.FakeFile{}

				cert = "invalid-cert"
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe(".Close", func() {

		BeforeEach(func() {
			fakeOs = &os_fake.FakeOs{}
			fakeFile = &os_fake.FakeFile{}
			fakeFile.NameReturns("/a/temp.file")

			database = brokerstore.NewPostgresVariantWithShims("username", "password", "host", "port", "dbName", "somefile", fakeSql, fakeIoUtil, fakeOs)
		})

		JustBeforeEach(func() {
			database.Close()
		})

		It("should delete temp file", func() {
			Expect(fakeOs.RemoveCallCount()).To(Equal(1))
			deletedFile := fakeOs.RemoveArgsForCall(0)
			Expect(deletedFile).To(Equal("somefile"))
		})
	})

	Describe(".Flavorify", func() {
		var (
			result string
		)

		BeforeEach(func() {
			fakeOs = &os_fake.FakeOs{}
			fakeFile = &os_fake.FakeFile{}
			fakeFile.NameReturns("/a/temp.file")

			database = brokerstore.NewPostgresVariantWithShims("username", "password", "host", "port", "dbName", "somefile", fakeSql, fakeIoUtil, fakeOs)
		})

		JustBeforeEach(func() {
			query := `INSERT INTO service_instances (id, value) VALUES (?, ?)`
			result = database.Flavorify(query)
		})

		It("should postgresify the query", func() {
			Expect(result).To(Equal(`INSERT INTO service_instances (id, value) VALUES ($1, $2)`))
		})
	})
})
