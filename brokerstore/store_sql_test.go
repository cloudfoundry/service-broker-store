package brokerstore_test

import (
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/cloudfoundry-incubator/service-broker-store/brokerstore"
	"github.com/pivotal-cf/brokerapi"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"strings"

	"code.cloudfoundry.org/goshims/sqlshim/sql_fake"
	"github.com/cloudfoundry-incubator/service-broker-store/brokerstore/brokerstorefakes"
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
		logger                                                           lager.Logger
		state                                                            brokerstore.DynamicState
		fakeSqlDb                                                        = &sql_fake.FakeSqlDB{}
		fakeVariant                                                      = &brokerstorefakes.FakeSqlVariant{}
		err                                                              error
		bindingID, serviceID, planID, orgGUID, spaceGUID, appGUID, share string
		serviceInstance                                                  brokerstore.ServiceInstance
		sqlStore                                                         brokerstore.SqlStore
		db                                                               *sql.DB
		mock                                                             sqlmock.Sqlmock
		bindResource                                                     brokerapi.BindResource
		parameters                                                       json.RawMessage
		bindDetails                                                      brokerapi.BindDetails
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test-broker")
		fakeVariant.ConnectReturns(fakeSqlDb, nil)
		fakeVariant.FlavorifyStub = func(query string) string {
			return query
		}
		store, err = brokerstore.NewSqlStoreWithVariant(logger, fakeVariant)
		Expect(err).ToNot(HaveOccurred())
		state = brokerstore.DynamicState{
			InstanceMap: map[string]brokerstore.ServiceInstance{
				"service-name": {
					ServiceFingerPrint: "server:/some-share",
				},
			},
			BindingMap: map[string]brokerapi.BindDetails{},
		}
		db, mock, err = sqlmock.New()
		sqlStore = brokerstore.SqlStore{Database: brokerstorefakes.FakeSQLMockConnection{db},
			StoreType: "mysql"}
	})

	It("should open a db connection", func() {
		Expect(fakeVariant.ConnectCallCount()).To(BeNumerically(">=", 1))
	})

	It("should create tables if they don't exist", func() {
		Expect(fakeSqlDb.ExecCallCount()).To(BeNumerically(">=", 2))
		Expect(fakeSqlDb.ExecArgsForCall(0)).To(ContainSubstring("CREATE TABLE IF NOT EXISTS service_instances"))
		Expect(fakeSqlDb.ExecArgsForCall(1)).To(ContainSubstring("CREATE TABLE IF NOT EXISTS service_bindings"))
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

	Describe("CreateInstanceDetails", func() {
		BeforeEach(func() {
			Expect(err).NotTo(HaveOccurred())
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
