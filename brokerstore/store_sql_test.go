package brokerstore_test

import (
	"errors"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"code.cloudfoundry.org/service-broker-store/brokerstore"
	"github.com/pivotal-cf/brokerapi"
	"github.com/DATA-DOG/go-sqlmock"

	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"strings"

	"code.cloudfoundry.org/goshims/sqlshim/sql_fake"
	"code.cloudfoundry.org/service-broker-store/brokerstorefakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type redactedStuff struct{}

func (a redactedStuff) Match(v driver.Value) bool {
	if s, ok := v.(string); ok {
		if strings.Contains(s, brokerstore.HashKey) {
			return true
		}
	}
	if b, ok := v.([]byte); ok {
		if strings.Contains(string(b), brokerstore.HashKey) {
			return true
		}
	}
	return false
}

var _ = Describe("SqlStore", func() {
	var (
		store                                                            brokerstore.Store
		sqlStore                                                         *brokerstore.SqlStore
		logger                                                           lager.Logger
		fakeSqlDb                                                        *sql_fake.FakeSqlDB
		fakeVariant                                                      *brokerstorefakes.FakeSqlVariant
		err                                                              error
		bindingID, serviceID, planID, orgGUID, spaceGUID, appGUID, share string
		serviceInstance                                                  brokerstore.ServiceInstance
		db                                                               *sql.DB
		mock                                                             sqlmock.Sqlmock
		bindResource                                                     brokerapi.BindResource
		parameters                                                       json.RawMessage
		bindDetails                                                      brokerapi.BindDetails
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test-broker")

		fakeSqlDb = &sql_fake.FakeSqlDB{}
		fakeVariant = &brokerstorefakes.FakeSqlVariant{}
		fakeVariant.ConnectReturns(fakeSqlDb, nil)
		fakeVariant.FlavorifyStub = func(query string) string {
			return query
		}
		store, _ = brokerstore.NewSqlStoreWithVariant(logger, fakeVariant)
		db, mock, err = sqlmock.New()
		Expect(err).ToNot(HaveOccurred())

		var ok bool
		sqlStore, ok = store.(*brokerstore.SqlStore)
		Expect(ok).To(BeTrue())
		sqlStore.Database = brokerstorefakes.FakeSQLMockConnection{SqlDB: db}
	})

	It("should open a db connection", func() {
		Expect(fakeVariant.ConnectCallCount()).To(BeNumerically(">=", 1))
	})

	It("should create tables if they don't exist", func() {
		Expect(fakeSqlDb.ExecCallCount()).To(BeNumerically(">=", 2))
		Expect(fakeSqlDb.ExecArgsForCall(0)).To(ContainSubstring("CREATE TABLE IF NOT EXISTS service_instances"))
		Expect(fakeSqlDb.ExecArgsForCall(1)).To(ContainSubstring("CREATE TABLE IF NOT EXISTS service_bindings"))
	})

	Describe("Retire", func() {
		var (
			err error
		)

		JustBeforeEach(func() {
			err = sqlStore.Retire()
		})

		Context("when the SQL insertion succeeds", func() {
			BeforeEach(func() {
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO service_instances").WithArgs("migrated-to-credhub", "true").WillReturnResult(result)
			})

			It("should succeed", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should insert the migration flag into the service_instances table", func() {
				Expect(mock.ExpectationsWereMet()).Should(Succeed())
			})
		})

		Context("when the SQL insertion fails", func() {
			BeforeEach(func() {
				mock.ExpectExec("INSERT INTO service_instances").WithArgs("migrated-to-credhub", "true").WillReturnError(errors.New("nope"))
			})

			It("should return the error", func() {
				Expect(err).To(MatchError("nope"))
			})
		})
	})

	Describe("IsRetired", func() {
		var (
			retired bool
			err     error

			columns []string
			rows    *sqlmock.Rows
		)

		BeforeEach(func() {
			columns = []string{"id", "value"}

			rows = sqlmock.NewRows(columns)
		})

		JustBeforeEach(func() {
			retired, err = sqlStore.IsRetired()
		})

		Context("given the store is still active", func() {
			BeforeEach(func() {
				mock.ExpectQuery("SELECT id, value FROM service_instances WHERE id = ?").WithArgs("migrated-to-credhub").WillReturnRows(rows)
			})

			It("should return a store and no error", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(retired).To(BeFalse())
			})
		})

		Context("given the store is retired", func() {
			BeforeEach(func() {
				rows.AddRow("migrated-to-credhub", "true")

				mock.ExpectQuery("SELECT id, value FROM service_instances WHERE id = ?").WithArgs("migrated-to-credhub").WillReturnRows(rows)
			})

			It("should return an error", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(retired).To(BeTrue())
			})
		})

		Context("given the retired check fails", func() {
			BeforeEach(func() {
				mock.ExpectQuery("SELECT id, value FROM service_instances WHERE id = ?").WithArgs("migrated-to-credhub").WillReturnError(errors.New("database-badness"))
			})

			It("should return the error from the database", func() {
				Expect(err).To(MatchError("database-badness"))
			})
		})
	})

	Describe("Restore", func() {
		BeforeEach(func() {
			err = store.Restore(logger)
		})

		It("this should be a noop", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Save", func() {
		BeforeEach(func() {
			err = store.Save(logger)
		})

		It("this should be a noop", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Cleanup", func() {
		BeforeEach(func() {
			err = store.Cleanup()
		})

		It("this should be a noop", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("RetrieveInstanceDetails", func() {
		Context("When the instance exists", func() {
			BeforeEach(func() {
				Expect(err).NotTo(HaveOccurred())
				serviceID = "instance_123"
				planID = "plan_123"
				orgGUID = "org_123"
				spaceGUID = "space_123"
				share = "share_123"

				columns := []string{"id", "value"}

				rows := sqlmock.NewRows(columns)
				jsonvalue, err := json.Marshal(brokerstore.ServiceInstance{ServiceFingerPrint: share, PlanID: planID, ServiceID: serviceID, OrganizationGUID: orgGUID, SpaceGUID: spaceGUID})
				Expect(err).NotTo(HaveOccurred())
				rows.AddRow(serviceID, jsonvalue)

				mock.ExpectQuery("SELECT id, value FROM service_instances WHERE id = ?").WithArgs(serviceID).WillReturnRows(rows)
			})

			JustBeforeEach(func() {
				serviceInstance, err = sqlStore.RetrieveInstanceDetails(serviceID)
			})

			It("should return the instance", func() {
				Expect(err).To(BeNil())
				Expect(mock.ExpectationsWereMet()).Should(Succeed())
				Expect(serviceInstance.ServiceID).To(Equal(serviceID))
				Expect(serviceInstance.PlanID).To(Equal(planID))
				Expect(serviceInstance.OrganizationGUID).To(Equal(orgGUID))
				Expect(serviceInstance.SpaceGUID).To(Equal(spaceGUID))
				Expect(serviceInstance.ServiceFingerPrint).To(Equal(share))
			})
		})

		Context("When the instance does not exist", func() {
			BeforeEach(func() {
				mock.ExpectQuery("SELECT id, value FROM service_instances WHERE id = ?").WithArgs(serviceID)
			})
			JustBeforeEach(func() {
				serviceInstance, err = sqlStore.RetrieveInstanceDetails(serviceID)
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(reflect.DeepEqual(serviceInstance, brokerstore.ServiceInstance{})).To(BeTrue())
			})
		})
	})

	Describe("RetrieveBindingDetails", func() {
		Context("When the instance exists", func() {
			BeforeEach(func() {
				Expect(err).NotTo(HaveOccurred())
				appGUID = "instance_123"
				planID = "plan_123"
				serviceID = "service_123"
				bindingID = "binding_123"
				bindResource = brokerapi.BindResource{AppGuid: appGUID, Route: "binding-route"}

				columns := []string{"id", "value"}
				rows := sqlmock.NewRows(columns)
				jsonvalue, err := json.Marshal(brokerapi.BindDetails{AppGUID: appGUID, PlanID: planID, ServiceID: serviceID, BindResource: &bindResource, RawParameters: parameters})
				Expect(err).NotTo(HaveOccurred())
				rows.AddRow(bindingID, jsonvalue)

				mock.ExpectQuery("SELECT id, value FROM service_bindings WHERE id = ?").WithArgs(bindingID).WillReturnRows(rows)
			})
			JustBeforeEach(func() {

				bindDetails, err = sqlStore.RetrieveBindingDetails(bindingID)
			})
			It("should return the binding details", func() {
				Expect(err).To(BeNil())
				Expect(mock.ExpectationsWereMet()).Should(Succeed())
				Expect(bindDetails.ServiceID).To(Equal(serviceID))
				Expect(bindDetails.PlanID).To(Equal(planID))
				Expect(bindDetails.AppGUID).To(Equal(appGUID))
				Expect(bindDetails.BindResource.AppGuid).To(Equal(appGUID))
				Expect(bindDetails.BindResource.Route).To(Equal("binding-route"))
				Expect(bindDetails.RawParameters).To(Equal(parameters))
			})
		})
		Context("When the binding does not exist", func() {
			BeforeEach(func() {
				mock.ExpectQuery("SELECT id, value FROM service_bindings WHERE id = ?").WithArgs(bindingID)
			})
			JustBeforeEach(func() {
				bindDetails, err = sqlStore.RetrieveBindingDetails(bindingID)
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(reflect.DeepEqual(bindDetails, brokerapi.BindDetails{})).To(BeTrue())
			})
		})
	})

	Describe("RetrieveAllInstanceDetails", func() {
		var serviceInstances map[string]brokerstore.ServiceInstance

		Context("when instances exist", func() {
			BeforeEach(func() {
				columns := []string{"id", "value"}
				rows := sqlmock.NewRows(columns)

				jsonvalue1, err := json.Marshal(brokerstore.ServiceInstance{
					ServiceFingerPrint: "share_1",
					PlanID:             "plan_1",
					ServiceID:          "service_1",
					OrganizationGUID:   "org_1",
					SpaceGUID:          "space_1",
				})
				Expect(err).NotTo(HaveOccurred())
				rows.AddRow("instance_1", jsonvalue1)

				jsonvalue2, err := json.Marshal(brokerstore.ServiceInstance{
					ServiceFingerPrint: "share_2",
					PlanID:             "plan_2",
					ServiceID:          "service_2",
					OrganizationGUID:   "org_2",
					SpaceGUID:          "space_2",
				})
				Expect(err).NotTo(HaveOccurred())
				rows.AddRow("instance_2", jsonvalue2)

				mock.ExpectQuery("SELECT id, value FROM service_instances").WillReturnRows(rows)
			})

			JustBeforeEach(func() {
				serviceInstances, err = sqlStore.RetrieveAllInstanceDetails()
			})

			It("should return the instance", func() {
				Expect(err).To(BeNil())
				Expect(mock.ExpectationsWereMet()).Should(Succeed())
				Expect(serviceInstances["instance_1"].ServiceID).To(Equal("service_1"))
				Expect(serviceInstances["instance_1"].PlanID).To(Equal("plan_1"))
				Expect(serviceInstances["instance_1"].OrganizationGUID).To(Equal("org_1"))
				Expect(serviceInstances["instance_1"].SpaceGUID).To(Equal("space_1"))
				Expect(serviceInstances["instance_1"].ServiceFingerPrint).To(Equal("share_1"))
				Expect(serviceInstances["instance_2"].ServiceID).To(Equal("service_2"))
				Expect(serviceInstances["instance_2"].PlanID).To(Equal("plan_2"))
				Expect(serviceInstances["instance_2"].OrganizationGUID).To(Equal("org_2"))
				Expect(serviceInstances["instance_2"].SpaceGUID).To(Equal("space_2"))
				Expect(serviceInstances["instance_2"].ServiceFingerPrint).To(Equal("share_2"))
			})
		})
	})

	Describe("RetrieveAllBindingDetails", func() {
		var (
			bindingDetails               map[string]brokerapi.BindDetails
			bindResource1, bindResource2 brokerapi.BindResource
		)

		Context("when instances exist", func() {
			BeforeEach(func() {
				columns := []string{"id", "value"}
				rows := sqlmock.NewRows(columns)

				bindResource1 = brokerapi.BindResource{AppGuid: "instance_1", Route: "binding-route-1"}
				jsonvalue1, err := json.Marshal(brokerapi.BindDetails{
					AppGUID:       "instance_1",
					PlanID:        "plan_1",
					ServiceID:     "service_1",
					BindResource:  &bindResource1,
					RawParameters: json.RawMessage([]byte(`{"a":"1"}`)),
				})
				Expect(err).NotTo(HaveOccurred())
				rows.AddRow("binding_1", jsonvalue1)

				bindResource2 = brokerapi.BindResource{AppGuid: "instance_2", Route: "binding-route-2"}
				jsonvalue2, err := json.Marshal(brokerapi.BindDetails{
					AppGUID:       "instance_2",
					PlanID:        "plan_2",
					ServiceID:     "service_2",
					BindResource:  &bindResource2,
					RawParameters: json.RawMessage([]byte(`{"a":"2"}`)),
				})
				Expect(err).NotTo(HaveOccurred())
				rows.AddRow("binding_2", jsonvalue2)

				mock.ExpectQuery("SELECT id, value FROM service_bindings").WillReturnRows(rows)
			})

			JustBeforeEach(func() {
				bindingDetails, err = sqlStore.RetrieveAllBindingDetails()
			})

			It("should return the instance", func() {
				Expect(err).To(BeNil())
				Expect(mock.ExpectationsWereMet()).Should(Succeed())
				Expect(bindingDetails["binding_1"].AppGUID).To(Equal("instance_1"))
				Expect(bindingDetails["binding_1"].PlanID).To(Equal("plan_1"))
				Expect(bindingDetails["binding_1"].ServiceID).To(Equal("service_1"))
				Expect(bindingDetails["binding_1"].BindResource).To(Equal(&bindResource1))
				Expect(bindingDetails["binding_1"].RawParameters).To(Equal(json.RawMessage([]byte(`{"a":"1"}`))))
				Expect(bindingDetails["binding_2"].AppGUID).To(Equal("instance_2"))
				Expect(bindingDetails["binding_2"].PlanID).To(Equal("plan_2"))
				Expect(bindingDetails["binding_2"].ServiceID).To(Equal("service_2"))
				Expect(bindingDetails["binding_2"].BindResource).To(Equal(&bindResource2))
				Expect(bindingDetails["binding_2"].RawParameters).To(Equal(json.RawMessage([]byte(`{"a":"2"}`))))
			})
		})
	})

	Describe("CreateInstanceDetails", func() {
		Context("when the service instance is valid", func() {
			BeforeEach(func() {
				orgGUID = "org_123"
				planID = "plan_123"
				serviceID = "service_123"
				spaceGUID = "space_123"
				share = "share_123"
				serviceInstance = brokerstore.ServiceInstance{ServiceID: serviceID, PlanID: planID, OrganizationGUID: orgGUID, SpaceGUID: spaceGUID, ServiceFingerPrint: share}
				jsonValue, err := json.Marshal(serviceInstance)
				Expect(err).NotTo(HaveOccurred())

				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO service_instances").WithArgs(serviceID, jsonValue).WillReturnResult(result)
			})
			JustBeforeEach(func() {
				err = sqlStore.CreateInstanceDetails(serviceID, serviceInstance)
			})
			It("should not error and call INSERT INTO on the db", func() {
				Expect(err).To(BeNil())
				Expect(mock.ExpectationsWereMet()).Should(Succeed())
			})
		})
		Context("when the service instance fingerprint contains a `password` key", func() {
			var fp interface{}
			BeforeEach(func() {
				orgGUID = "org_123"
				planID = "plan_123"
				serviceID = "service_123"
				spaceGUID = "space_123"
				fp = map[string]string{"password": "terribleSecrets"}
			})
			JustBeforeEach(func() {
				serviceInstance = brokerstore.ServiceInstance{ServiceID: serviceID, PlanID: planID, OrganizationGUID: orgGUID, SpaceGUID: spaceGUID, ServiceFingerPrint: fp}
				jsonValue, err2 := json.Marshal(serviceInstance)
				Expect(err2).NotTo(HaveOccurred())

				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO service_instances").WithArgs(serviceID, jsonValue).WillReturnResult(result)
				err = sqlStore.CreateInstanceDetails(serviceID, serviceInstance)
			})
			It("should fail", func() {
				Expect(err).NotTo(BeNil())
			})
			Context("when there's an array of stuff", func() {
				BeforeEach(func() {
					fp = []interface{}{
						map[string]string{"notapassword": "terribleSecrets"},
						map[string]string{"alsonotapassword": "terribleSecrets"},
						map[string]string{"password": "terribleSecrets"},
					}
				})
				It("should fail", func() {
					Expect(err).NotTo(BeNil())
				})
			})
		})

	})

	Describe("CreateBindingDetails", func() {
		BeforeEach(func() {
			Expect(err).NotTo(HaveOccurred())
			appGUID = "instance_123"
			planID = "plan_123"
			serviceID = "service_123"
			bindingID = "binding_123"
			bindResource = brokerapi.BindResource{AppGuid: appGUID, Route: "binding-route"}
			bindDetails = brokerapi.BindDetails{AppGUID: appGUID, PlanID: planID, ServiceID: serviceID, BindResource: &bindResource, RawParameters: parameters}
		})
		JustBeforeEach(func() {
			err = sqlStore.CreateBindingDetails(bindingID, bindDetails)
		})

		Context("when there are no parameters in the binding", func() {
			BeforeEach(func() {
				jsonValue, err := json.Marshal(bindDetails)
				Expect(err).NotTo(HaveOccurred())

				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO service_bindings").WithArgs(bindingID, jsonValue).WillReturnResult(result)
			})
			It("should not error and call INSERT INTO on the db", func() {
				Expect(err).To(BeNil())
				Expect(mock.ExpectationsWereMet()).Should(Succeed())
			})
		})

		Context("when there are parameters with secrets in the binding", func() {
			BeforeEach(func() {
				bindParameters := map[string]interface{}{"secret": "don't tell"}
				bindMessage, err := json.Marshal(bindParameters)
				Expect(err).NotTo(HaveOccurred())
				bindDetails = brokerapi.BindDetails{AppGUID: appGUID, PlanID: planID, ServiceID: serviceID, BindResource: &bindResource, RawParameters: bindMessage}
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO service_bindings").WithArgs(bindingID, &redactedStuff{}).WillReturnResult(result)
			})
			It("should redact parameters before saving records to the db", func() {
				Expect(err).To(BeNil())
				Expect(mock.ExpectationsWereMet()).Should(Succeed())
			})
		})
	})

	Describe("DeleteInstanceDetails", func() {
		BeforeEach(func() {
			Expect(err).NotTo(HaveOccurred())
			serviceID = "my_service"
			result := sqlmock.NewResult(1, 1)
			mock.ExpectExec("DELETE FROM service_instances WHERE id = ?").WithArgs(serviceID).WillReturnResult(result)
		})
		JustBeforeEach(func() {
			err = sqlStore.DeleteInstanceDetails(serviceID)
		})
		It("should not error and call DELETE FROM on the db", func() {
			Expect(err).To(BeNil())
			Expect(mock.ExpectationsWereMet()).Should(Succeed())
		})
	})

	Describe("DeleteBindingDetails", func() {
		BeforeEach(func() {
			Expect(err).NotTo(HaveOccurred())
			bindingID = "my_binding"
			result := sqlmock.NewResult(1, 1)
			mock.ExpectExec("DELETE FROM service_bindings WHERE id = ?").WithArgs(bindingID).WillReturnResult(result)
		})
		JustBeforeEach(func() {
			err = sqlStore.DeleteBindingDetails(bindingID)
		})
		It("should not error and call DELETE FROM on the db", func() {
			Expect(err).To(BeNil())
			Expect(mock.ExpectationsWereMet()).Should(Succeed())
		})
	})
})
