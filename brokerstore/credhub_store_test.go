package brokerstore_test

import (
	"encoding/json"
	"errors"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
	"code.cloudfoundry.org/lager/v3"
	"code.cloudfoundry.org/lager/v3/lagertest"
	. "code.cloudfoundry.org/service-broker-store/brokerstore"
	"code.cloudfoundry.org/service-broker-store/brokerstore/credhub_shims/credhub_fakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi/v10"
	"github.com/pivotal-cf/brokerapi/v10/domain"
	"golang.org/x/crypto/bcrypt"
)

var _ = Describe("CredhubStore", func() {
	var (
		logger      lager.Logger
		fakeCredhub *credhub_fakes.FakeCredhub
		store       *CredhubStore
		err         error
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("CredHubStoreTest")
		fakeCredhub = &credhub_fakes.FakeCredhub{}
		store = NewCredhubStore(logger, fakeCredhub, "some-store-id")
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
			id, value := fakeCredhub.SetJSONArgsForCall(0)
			Expect(id).To(Equal("/some-store-id/12345"))
			actualJSON, err := json.Marshal(value)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualJSON).To(MatchJSON(expectedJSON))
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
			Expect(id).To(Equal("/some-store-id/12345"))
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
			Expect(id).To(Equal("/some-store-id/12345"))
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
			bindDetails  domain.BindDetails
			expectedJSON string
		)

		BeforeEach(func() {
			id = "12345"
			bindDetails = domain.BindDetails{
				AppGUID:       "app-guid",
				PlanID:        "plan-id",
				ServiceID:     "service-id",
				BindResource:  &domain.BindResource{AppGuid: "app-guid", Route: "my-app.cf.com"},
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
			id, value := fakeCredhub.SetJSONArgsForCall(0)
			Expect(id).To(Equal("/some-store-id/12345"))
			actualJSON, err := json.Marshal(value)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualJSON).To(MatchJSON(expectedJSON))
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
			Expect(id).To(Equal("/some-store-id/12345"))
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
			Expect(id).To(Equal("/some-store-id/12345"))
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
			id           string
			hashedParams string
		)

		BeforeEach(func() {
			id = "12345"
			hashedParamsBytes, err := bcrypt.GenerateFromPassword(
				[]byte(`{"username": "a-username", "password": "a-password"}`),
				bcrypt.DefaultCost,
			)
			Expect(err).NotTo(HaveOccurred())
			hashedParams = string(hashedParamsBytes)

			json := credentials.JSON{
				Value: values.JSON{
					"app_guid":   "app-guid",
					"plan_id":    "plan-id",
					"service_id": "service-id",
					"bind_resource": map[string]interface{}{
						"app_guid": "app-guid",
						"route":    "my-app.cf.com",
					},
					"parameters": map[string]interface{}{HashKey: hashedParams},
				},
			}
			fakeCredhub.GetLatestJSONReturns(json, nil)
		})

		It("returns false when instance details are the same", func() {
			isConflict := store.IsBindingConflict(id, brokerapi.BindDetails{
				AppGUID:       "app-guid",
				PlanID:        "plan-id",
				ServiceID:     "service-id",
				BindResource:  &domain.BindResource{AppGuid: "app-guid", Route: "my-app.cf.com"},
				RawParameters: json.RawMessage(`{"username": "a-username", "password": "a-password"}`),
			})
			Expect(isConflict).To(BeFalse())
		})

		It("returns true when instance details are the different", func() {
			isConflict := store.IsBindingConflict(id, brokerapi.BindDetails{
				AppGUID:       "app-guid",
				PlanID:        "other-plan-id",
				ServiceID:     "other-service-id",
				BindResource:  &brokerapi.BindResource{AppGuid: "app-guid", Route: "my-app.cf.com"},
				RawParameters: json.RawMessage(`{"paramsHash": "some-other-hash"}`),
			})
			Expect(isConflict).To(BeTrue())
		})
	})

	Context("#Activate", func() {
		JustBeforeEach(func() {
			err = store.Activate()
		})

		It("should write the record into the store", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeCredhub.SetValueCallCount()).To(Equal(1))
			id, value := fakeCredhub.SetValueArgsForCall(0)
			Expect(id).To(Equal("/some-store-id/migrated-from-sql"))
			Expect(value).To(Equal(values.Value("true")))
		})

		Context("when the record insertion fails", func() {
			BeforeEach(func() {
				fakeCredhub.SetValueReturns(credentials.Value{}, errors.New("bad-set-value"))
			})

			It("should return the error from credhub", func() {
				Expect(err).To(MatchError("bad-set-value"))
			})
		})
	})

	Context("#IsActivated", func() {
		var activated bool

		JustBeforeEach(func() {
			activated, err = store.IsActivated()
		})

		Context("when the migration from SQL has run", func() {
			BeforeEach(func() {
				fakeCredhub.FindByPathReturns(credentials.FindResults{
					Credentials: []struct {
						Name             string `json:"name" yaml:"name"`
						VersionCreatedAt string `json:"version_created_at" yaml:"version_created_at"`
					}{
						{Name: "migrated-from-sql"},
					},
				}, nil)
			})

			It("should return true", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(activated).To(BeTrue())
			})
		})

		Context("when the migration from SQL has not run", func() {
			It("should return false", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(activated).To(BeFalse())
			})
		})

		Context("when the record insertion fails", func() {
			BeforeEach(func() {
				fakeCredhub.FindByPathReturns(credentials.FindResults{}, errors.New("bad-find-by-path"))
			})

			It("should return the error from credhub", func() {
				Expect(err).To(MatchError("bad-find-by-path"))
			})
		})
	})
})
