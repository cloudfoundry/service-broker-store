package brokerstore_test

import (
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	. "code.cloudfoundry.org/service-broker-store/brokerstore"
	"code.cloudfoundry.org/service-broker-store/brokerstore/credhub_shims/credhub_fakes"

	"encoding/json"
	"errors"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

//type Store interface {
//	RetrieveInstanceDetails(id string) (ServiceInstance, error)
//	RetrieveBindingDetails(id string) (brokerapi.BindDetails, error)
//
//	CreateInstanceDetails(id string, details ServiceInstance) error
//	CreateBindingDetails(id string, details brokerapi.BindDetails) error
//
//	DeleteInstanceDetails(id string) error
//	DeleteBindingDetails(id string) error
//
//	IsInstanceConflict(id string, details ServiceInstance) bool
//	IsBindingConflict(id string, details brokerapi.BindDetails) bool
//
//	Restore(logger lager.Logger) error
//	Save(logger lager.Logger) error
//	Cleanup() error
//}

var _ = Describe("CredhubStore", func() {

	var (
		logger      lager.Logger
		fakeCredhub *credhub_fakes.FakeCredhub
		store       Store
		err         error
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("CredHubStoreTest")
		fakeCredhub = &credhub_fakes.FakeCredhub{}
		store = NewCredhubStore(logger, fakeCredhub)
	})

	Context("#CreateInstanceDetails", func() {
		var (
			id              string
			serviceInstance ServiceInstance
			expectedJSON    string
		)

		BeforeEach(func() {
			id = "12345"
			fp := map[string]interface{}{
				"username": "a-username",
				"password": "a-password",
			}
			serviceInstance = ServiceInstance{
				ServiceID:          "service-id",
				PlanID:             "plan-id",
				OrganizationGUID:   "org-guid",
				SpaceGUID:          "space-guid",
				ServiceFingerPrint: fp,
			}
			expectedJSON = `{
					"service_id":         "service-id",
					"plan_id":            "plan-id",
					"organization_guid":  "org-guid",
					"space_guid":         "space-guid",
					"ServiceFingerPrint": {"username": "a-username", "password": "a-password"}
				}`
		})

		JustBeforeEach(func() {
			err = store.CreateInstanceDetails(id, serviceInstance)
		})

		It("should store it in credhub", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeCredhub.SetJSONCallCount()).To(Equal(1))
			id, value, mode := fakeCredhub.SetJSONArgsForCall(0)
			Expect(id).To(Equal("12345"))
			actualJSON, err := json.Marshal(value)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualJSON).To(MatchJSON(expectedJSON))
			Expect(mode).To(Equal(credhub.NoOverwrite))
		})

		Context("when SetJSON returns an error", func() {
			BeforeEach(func() {
				fakeCredhub.SetJSONReturns(credentials.JSON{}, errors.New("bad-set-json"))
			})
			It("should return the error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("bad-set-json"))
			})
		})
	})

	Context("#RetrieveInstanceDetails", func() {
		var (
			id              string
			serviceInstance ServiceInstance
		)

		BeforeEach(func() {
			id = "12345"
			json := credentials.JSON{
				Value: values.JSON{
					"service_id":         "service-id",
					"plan_id":            "plan-id",
					"organization_guid":  "org-guid",
					"space_guid":         "space-guid",
					"ServiceFingerPrint": map[string]interface{}{"username": "a-username", "password": "a-password"},
				},
			}

			fakeCredhub.GetLatestJSONReturns(json, nil)
		})

		JustBeforeEach(func() {
			serviceInstance, err = store.RetrieveInstanceDetails(id)
		})

		It("should retrieve them in credhub", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeCredhub.GetLatestJSONCallCount()).To(Equal(1))
			id := fakeCredhub.GetLatestJSONArgsForCall(0)
			Expect(id).To(Equal("12345"))
			Expect(serviceInstance).To(Equal(ServiceInstance{
				ServiceID:          "service-id",
				PlanID:             "plan-id",
				OrganizationGUID:   "org-guid",
				SpaceGUID:          "space-guid",
				ServiceFingerPrint: map[string]interface{}{"username": "a-username", "password": "a-password"},
			}))
		})

		Context("when credhub returns an error", func() {
			BeforeEach(func() {
				fakeCredhub.GetLatestJSONReturns(credentials.JSON{}, errors.New("bad-get-latest-json"))
			})
			It("should return the error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("bad-get-latest-json"))
			})
		})
	})

	Context("#RetrieveBindingDetails", func() {
		var (
			id             string
			bindingDetails brokerapi.BindDetails
		)

		BeforeEach(func() {
			id = "12345"
			json := credentials.JSON{
				Value: values.JSON{
					"app_guid":   "app-guid",
					"plan_id":    "plan-id",
					"service_id": "service-id",
					"bind_resource": map[string]interface{}{
						"app_guid": "app-guid",
						"route":    "my-app.cf.com",
					},
					"parameters": map[string]interface{}{"username": "a-username", "password": "a-password"},
				},
			}
			fakeCredhub.GetLatestJSONReturns(json, nil)
		})

		JustBeforeEach(func() {
			bindingDetails, err = store.RetrieveBindingDetails(id)
		})

		It("should retrieve them from credhub", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeCredhub.GetLatestJSONCallCount()).To(Equal(1))
			id := fakeCredhub.GetLatestJSONArgsForCall(0)
			Expect(id).To(Equal("12345"))
			Expect(bindingDetails).To(Equal(brokerapi.BindDetails{
				AppGUID:       "app-guid",
				PlanID:        "plan-id",
				ServiceID:     "service-id",
				BindResource:  &brokerapi.BindResource{AppGuid: "app-guid", Route: "my-app.cf.com"},
				RawParameters: json.RawMessage([]byte(`{"password":"a-password","username":"a-username"}`)),
			}))
		})

		Context("when credhub returns an error", func() {
			BeforeEach(func() {
				fakeCredhub.GetLatestJSONReturns(credentials.JSON{}, errors.New("bad-get-latest-json"))
			})
			It("should return the error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("bad-get-latest-json"))
			})
		})
	})

	Context("#CreateBindingDetails", func() {
		var (
			id           string
			bindDetails  brokerapi.BindDetails
			expectedJSON string
		)

		BeforeEach(func() {
			id = "12345"
			bindDetails = brokerapi.BindDetails{
				AppGUID:       "app-guid",
				PlanID:        "plan-id",
				ServiceID:     "service-id",
				BindResource:  &brokerapi.BindResource{AppGuid: "app-guid", Route: "my-app.cf.com"},
				RawParameters: json.RawMessage([]byte(`{"password":"a-password","username":"a-username"}`)),
			}
			expectedJSON = `{
					"app_guid":      "app-guid",
					"plan_id":       "plan-id",
					"service_id":    "service-id",
					"bind_resource": {"app_guid": "app-guid", "route": "my-app.cf.com"},
					"parameters":    {"username": "a-username", "password": "a-password"}
				}`
		})

		JustBeforeEach(func() {
			err = store.CreateBindingDetails(id, bindDetails)
		})

		It("should store it in credhub", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeCredhub.SetJSONCallCount()).To(Equal(1))
			id, value, mode := fakeCredhub.SetJSONArgsForCall(0)
			Expect(id).To(Equal("12345"))
			actualJSON, err := json.Marshal(value)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualJSON).To(MatchJSON(expectedJSON))
			Expect(mode).To(Equal(credhub.NoOverwrite))
		})

		Context("when SetJSON returns an error", func() {
			BeforeEach(func() {
				fakeCredhub.SetJSONReturns(credentials.JSON{}, errors.New("bad-set-json"))
			})

			It("should return the error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("bad-set-json"))
			})
		})
	})

	Context("#CreateBindingDetails", func() {
		var (
			id           string
			bindDetails  brokerapi.BindDetails
			expectedJSON string
		)

		BeforeEach(func() {
			id = "12345"
			bindDetails = brokerapi.BindDetails{
				AppGUID:       "app-guid",
				PlanID:        "plan-id",
				ServiceID:     "service-id",
				BindResource:  &brokerapi.BindResource{AppGuid: "app-guid", Route: "my-app.cf.com"},
				RawParameters: json.RawMessage([]byte(`{"password":"a-password","username":"a-username"}`)),
			}
			expectedJSON = `{
					"app_guid":      "app-guid",
					"plan_id":       "plan-id",
					"service_id":    "service-id",
					"bind_resource": {"app_guid": "app-guid", "route": "my-app.cf.com"},
					"parameters":    {"username": "a-username", "password": "a-password"}
				}`
		})

		JustBeforeEach(func() {
			err = store.CreateBindingDetails(id, bindDetails)
		})

		It("should store it in credhub", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeCredhub.SetJSONCallCount()).To(Equal(1))
			id, value, mode := fakeCredhub.SetJSONArgsForCall(0)
			Expect(id).To(Equal("12345"))
			actualJSON, err := json.Marshal(value)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualJSON).To(MatchJSON(expectedJSON))
			Expect(mode).To(Equal(credhub.NoOverwrite))
		})

		Context("when SetJSON returns an error", func() {
			BeforeEach(func() {
				fakeCredhub.SetJSONReturns(credentials.JSON{}, errors.New("bad-set-json"))
			})

			It("should return the error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("bad-set-json"))
			})
		})
	})

	Context("#DeleteInstanceDetails", func() {
		var (
			id string
		)

		BeforeEach(func() {
			id = "12345"
		})

		JustBeforeEach(func() {
			err = store.DeleteInstanceDetails(id)
		})

		It("should remove key from credhub", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeCredhub.DeleteCallCount()).To(Equal(1))
			id := fakeCredhub.DeleteArgsForCall(0)
			Expect(id).To(Equal("12345"))
		})

		Context("when Delete returns an error", func() {
			BeforeEach(func() {
				fakeCredhub.DeleteReturns(errors.New("bad-delete"))
			})
			It("should return the error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("bad-delete"))
			})
		})
	})

	Context("#DeleteBindingDetails", func() {
		var (
			id string
		)

		BeforeEach(func() {
			id = "12345"
		})

		JustBeforeEach(func() {
			err = store.DeleteBindingDetails(id)
		})

		It("should store the binding from credhub", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeCredhub.DeleteCallCount()).To(Equal(1))
			id := fakeCredhub.DeleteArgsForCall(0)
			Expect(id).To(Equal("12345"))
		})

		Context("when SetJSON returns an error", func() {
			BeforeEach(func() {
				fakeCredhub.DeleteReturns(errors.New("bad-delete"))
			})

			It("should return the error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("bad-delete"))
			})
		})
	})

	Context("#IsInstanceConflict", func() {
		var id string

		BeforeEach(func() {
			id = "12345"
			json := credentials.JSON{
				Value: values.JSON{
					"service_id":         "service-id",
					"plan_id":            "plan-id",
					"organization_guid":  "org-guid",
					"space_guid":         "space-guid",
					"ServiceFingerPrint": map[string]interface{}{"username": "a-username", "password": "a-password"},
				},
			}

			fakeCredhub.GetLatestJSONReturns(json, nil)
		})

		It("returns false when instance details are the same", func() {
			isConflict := store.IsInstanceConflict(id, ServiceInstance{
				ServiceID:          "service-id",
				PlanID:             "plan-id",
				OrganizationGUID:   "org-guid",
				SpaceGUID:          "space-guid",
				ServiceFingerPrint: map[string]interface{}{"username": "a-username", "password": "a-password"},
			})
			Expect(isConflict).To(BeFalse())
		})

		It("returns true when instance details are the different", func() {
			isConflict := store.IsInstanceConflict(id, ServiceInstance{
				ServiceID:          "service-id",
				PlanID:             "other-plan-id",
				OrganizationGUID:   "other-org-guid",
				SpaceGUID:          "other-space-guid",
				ServiceFingerPrint: map[string]interface{}{"username": "b-username", "password": "b-password"},
			})
			Expect(isConflict).To(BeTrue())
		})
	})

	Context("#IsBindingConflict", func() {
		var (
			id             string
			hashedPassword string
		)

		BeforeEach(func() {
			id = "12345"
			hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(`"username": "a-username", "password": "a-password"}`), bcrypt.DefaultCost)
			Expect(err).NotTo(HaveOccurred())

			json := credentials.JSON{
				Value: values.JSON{
					"app_guid":   "app-guid",
					"plan_id":    "plan-id",
					"service_id": "service-id",
					"bind_resource": map[string]interface{}{
						"app_guid": "app-guid",
						"route":    "my-app.cf.com",
					},
					"parameters": map[string]interface{}{HashKey: hashedPassword},
				},
			}
			fakeCredhub.GetLatestJSONReturns(json, nil)
		})

		FIt("returns false when instance details are the same", func() {
			isConflict := store.IsBindingConflict(id, brokerapi.BindDetails{
				AppGUID:       "app-guid",
				PlanID:        "plan-id",
				ServiceID:     "service-id",
				BindResource:  &brokerapi.BindResource{AppGuid: "app-guid", Route: "my-app.cf.com"},
				RawParameters: json.RawMessage([]byte(`{"paramsHash":"` + hashedPassword + `"}`)),
			})
			Expect(isConflict).To(BeFalse())
		})

		It("returns true when instance details are the different", func() {
			isConflict := store.IsBindingConflict(id, brokerapi.BindDetails{
				AppGUID:       "app-guid",
				PlanID:        "other-plan-id",
				ServiceID:     "other-service-id",
				BindResource:  &brokerapi.BindResource{AppGuid: "app-guid", Route: "my-app.cf.com"},
				RawParameters: json.RawMessage([]byte(`{"paramsHash": "some-other-hash"}`)),
			})
			Expect(isConflict).To(BeTrue())
		})
	})
})
