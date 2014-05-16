package api_test

import (
	"github.com/pivotal-cf/go-service-broker/api"
)

type FakeServiceBroker struct {
	ProvisionedInstanceIDs   []string
	DeprovisionedInstanceIDs []string

	BoundInstanceIDs []string
	BoundBindingIDs  []string

	InstanceLimit int

	ProvisionError error
	BindError      error
}

func (fakeBroker FakeServiceBroker) Services() []api.Service {
	return []api.Service{
		api.Service{
			ID:          "0A789746-596F-4CEA-BFAC-A0795DA056E3",
			Name:        "p-cassandra",
			Description: "Cassandra service for application development and testing",
			Bindable:    true,
			Plans: []api.ServicePlan{
				api.ServicePlan{
					ID:          "ABE176EE-F69F-4A96-80CE-142595CC24E3",
					Name:        "default",
					Description: "The default Cassandra plan",
					Metadata: api.ServicePlanMetadata{
						Bullets:     []string{},
						DisplayName: "Cassandra",
					},
				},
			},
			Metadata: api.ServiceMetadata{
				DisplayName:      "Cassandra",
				LongDescription:  "Long description",
				DocumentationUrl: "http://thedocs.com",
				SupportUrl:       "http://helpme.no",
				Listing: api.ServiceMetadataListing{
					Blurb:    "blah blah",
					ImageUrl: "http://foo.com/thing.png",
				},
				Provider: api.ServiceMetadataProvider{
					Name: "Pivotal",
				},
			},
			Tags: []string{
				"pivotal",
				"cassandra",
			},
		},
	}
}

func (fakeBroker *FakeServiceBroker) Provision(instanceID string) error {
	if fakeBroker.ProvisionError != nil {
		return fakeBroker.ProvisionError
	}

	if len(fakeBroker.ProvisionedInstanceIDs) >= fakeBroker.InstanceLimit {
		return api.ErrInstanceLimitMet
	}

	if sliceContains(instanceID, fakeBroker.ProvisionedInstanceIDs) {
		return api.ErrInstanceAlreadyExists
	}

	fakeBroker.ProvisionedInstanceIDs = append(fakeBroker.ProvisionedInstanceIDs, instanceID)
	return nil
}

func (fakeBroker *FakeServiceBroker) Deprovision(instanceID string) error {
	fakeBroker.DeprovisionedInstanceIDs = append(fakeBroker.DeprovisionedInstanceIDs, instanceID)

	if sliceContains(instanceID, fakeBroker.ProvisionedInstanceIDs) {
		return nil
	}
	return api.ErrInstanceDoesNotExist
}

func (fakeBroker *FakeServiceBroker) Bind(instanceID, bindingID string) (interface{}, error) {
	if fakeBroker.BindError != nil {
		return nil, fakeBroker.BindError
	}

	fakeBroker.BoundInstanceIDs = append(fakeBroker.BoundInstanceIDs, instanceID)
	fakeBroker.BoundBindingIDs = append(fakeBroker.BoundBindingIDs, bindingID)

	return FakeCredentials{
		Host:     "127.0.0.1",
		Port:     3000,
		Username: "batman",
		Password: "robin",
	}, nil
}

func (fakeBroker *FakeServiceBroker) Unbind(instanceID, bindingID string) error {
	if sliceContains(instanceID, fakeBroker.ProvisionedInstanceIDs) {
		if sliceContains(bindingID, fakeBroker.BoundBindingIDs) {
			return nil
		}
		return api.ErrBindingDoesNotExist
	}

	return api.ErrInstanceDoesNotExist
}

type FakeCredentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
