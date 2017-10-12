package brokerstore_test

import (
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/cloudfoundry-incubator/service-broker-store/brokerstore"

	"errors"

	"code.cloudfoundry.org/goshims/sqlshim/sql_fake"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MysqlVariant", func() {

	var (
		logger   lager.Logger
		fakeSql  *sql_fake.FakeSql
		err      error
		database brokerstore.SqlVariant

		cert string
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("mysql-variant-test")

		fakeSql = &sql_fake.FakeSql{}
	})

	JustBeforeEach(func() {
		database = brokerstore.NewMySqlVariantWithSqlObject("username", "password", "host", "port", "dbName", cert, fakeSql)
	})

	Describe(".Connect", func() {

		JustBeforeEach(func() {
			_, err = database.Connect(logger)
		})

		Context("when ca cert specified", func() {
			BeforeEach(func() {
				cert = exampleCaCert
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
})
