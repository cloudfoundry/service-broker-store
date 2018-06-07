// Code generated by counterfeiter. DO NOT EDIT.
package credhub_fakes

import (
	"sync"

	"code.cloudfoundry.org/service-broker-store/brokerstore/credhub_shims"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/credentials/values"
)

type FakeCredhub struct {
	SetJSONStub        func(name string, value values.JSON, overwrite credhub.Mode) (credentials.JSON, error)
	setJSONMutex       sync.RWMutex
	setJSONArgsForCall []struct {
		name      string
		value     values.JSON
		overwrite credhub.Mode
	}
	setJSONReturns struct {
		result1 credentials.JSON
		result2 error
	}
	setJSONReturnsOnCall map[int]struct {
		result1 credentials.JSON
		result2 error
	}
	GetLatestJSONStub        func(name string) (credentials.JSON, error)
	getLatestJSONMutex       sync.RWMutex
	getLatestJSONArgsForCall []struct {
		name string
	}
	getLatestJSONReturns struct {
		result1 credentials.JSON
		result2 error
	}
	getLatestJSONReturnsOnCall map[int]struct {
		result1 credentials.JSON
		result2 error
	}
	DeleteStub        func(name string) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		name string
	}
	deleteReturns struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCredhub) SetJSON(name string, value values.JSON, overwrite credhub.Mode) (credentials.JSON, error) {
	fake.setJSONMutex.Lock()
	ret, specificReturn := fake.setJSONReturnsOnCall[len(fake.setJSONArgsForCall)]
	fake.setJSONArgsForCall = append(fake.setJSONArgsForCall, struct {
		name      string
		value     values.JSON
		overwrite credhub.Mode
	}{name, value, overwrite})
	fake.recordInvocation("SetJSON", []interface{}{name, value, overwrite})
	fake.setJSONMutex.Unlock()
	if fake.SetJSONStub != nil {
		return fake.SetJSONStub(name, value, overwrite)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.setJSONReturns.result1, fake.setJSONReturns.result2
}

func (fake *FakeCredhub) SetJSONCallCount() int {
	fake.setJSONMutex.RLock()
	defer fake.setJSONMutex.RUnlock()
	return len(fake.setJSONArgsForCall)
}

func (fake *FakeCredhub) SetJSONArgsForCall(i int) (string, values.JSON, credhub.Mode) {
	fake.setJSONMutex.RLock()
	defer fake.setJSONMutex.RUnlock()
	return fake.setJSONArgsForCall[i].name, fake.setJSONArgsForCall[i].value, fake.setJSONArgsForCall[i].overwrite
}

func (fake *FakeCredhub) SetJSONReturns(result1 credentials.JSON, result2 error) {
	fake.SetJSONStub = nil
	fake.setJSONReturns = struct {
		result1 credentials.JSON
		result2 error
	}{result1, result2}
}

func (fake *FakeCredhub) SetJSONReturnsOnCall(i int, result1 credentials.JSON, result2 error) {
	fake.SetJSONStub = nil
	if fake.setJSONReturnsOnCall == nil {
		fake.setJSONReturnsOnCall = make(map[int]struct {
			result1 credentials.JSON
			result2 error
		})
	}
	fake.setJSONReturnsOnCall[i] = struct {
		result1 credentials.JSON
		result2 error
	}{result1, result2}
}

func (fake *FakeCredhub) GetLatestJSON(name string) (credentials.JSON, error) {
	fake.getLatestJSONMutex.Lock()
	ret, specificReturn := fake.getLatestJSONReturnsOnCall[len(fake.getLatestJSONArgsForCall)]
	fake.getLatestJSONArgsForCall = append(fake.getLatestJSONArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("GetLatestJSON", []interface{}{name})
	fake.getLatestJSONMutex.Unlock()
	if fake.GetLatestJSONStub != nil {
		return fake.GetLatestJSONStub(name)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getLatestJSONReturns.result1, fake.getLatestJSONReturns.result2
}

func (fake *FakeCredhub) GetLatestJSONCallCount() int {
	fake.getLatestJSONMutex.RLock()
	defer fake.getLatestJSONMutex.RUnlock()
	return len(fake.getLatestJSONArgsForCall)
}

func (fake *FakeCredhub) GetLatestJSONArgsForCall(i int) string {
	fake.getLatestJSONMutex.RLock()
	defer fake.getLatestJSONMutex.RUnlock()
	return fake.getLatestJSONArgsForCall[i].name
}

func (fake *FakeCredhub) GetLatestJSONReturns(result1 credentials.JSON, result2 error) {
	fake.GetLatestJSONStub = nil
	fake.getLatestJSONReturns = struct {
		result1 credentials.JSON
		result2 error
	}{result1, result2}
}

func (fake *FakeCredhub) GetLatestJSONReturnsOnCall(i int, result1 credentials.JSON, result2 error) {
	fake.GetLatestJSONStub = nil
	if fake.getLatestJSONReturnsOnCall == nil {
		fake.getLatestJSONReturnsOnCall = make(map[int]struct {
			result1 credentials.JSON
			result2 error
		})
	}
	fake.getLatestJSONReturnsOnCall[i] = struct {
		result1 credentials.JSON
		result2 error
	}{result1, result2}
}

func (fake *FakeCredhub) Delete(name string) error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("Delete", []interface{}{name})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deleteReturns.result1
}

func (fake *FakeCredhub) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeCredhub) DeleteArgsForCall(i int) string {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return fake.deleteArgsForCall[i].name
}

func (fake *FakeCredhub) DeleteReturns(result1 error) {
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCredhub) DeleteReturnsOnCall(i int, result1 error) {
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCredhub) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.setJSONMutex.RLock()
	defer fake.setJSONMutex.RUnlock()
	fake.getLatestJSONMutex.RLock()
	defer fake.getLatestJSONMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCredhub) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ credhub_shims.Credhub = new(FakeCredhub)
