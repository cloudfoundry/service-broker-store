// Code generated by counterfeiter. DO NOT EDIT.
package credhub_fakes

import (
	"sync"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
	"code.cloudfoundry.org/service-broker-store/brokerstore/credhub_shims"
)

type FakeCredhub struct {
	SetJSONStub        func(name string, value values.JSON) (credentials.JSON, error)
	setJSONMutex       sync.RWMutex
	setJSONArgsForCall []struct {
		name  string
		value values.JSON
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
	SetValueStub        func(name string, value values.Value) (credentials.Value, error)
	setValueMutex       sync.RWMutex
	setValueArgsForCall []struct {
		name  string
		value values.Value
	}
	setValueReturns struct {
		result1 credentials.Value
		result2 error
	}
	setValueReturnsOnCall map[int]struct {
		result1 credentials.Value
		result2 error
	}
	GetLatestValueStub        func(name string) (credentials.Value, error)
	getLatestValueMutex       sync.RWMutex
	getLatestValueArgsForCall []struct {
		name string
	}
	getLatestValueReturns struct {
		result1 credentials.Value
		result2 error
	}
	getLatestValueReturnsOnCall map[int]struct {
		result1 credentials.Value
		result2 error
	}
	FindByPathStub        func(path string) (credentials.FindResults, error)
	findByPathMutex       sync.RWMutex
	findByPathArgsForCall []struct {
		path string
	}
	findByPathReturns struct {
		result1 credentials.FindResults
		result2 error
	}
	findByPathReturnsOnCall map[int]struct {
		result1 credentials.FindResults
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

func (fake *FakeCredhub) SetJSON(name string, value values.JSON) (credentials.JSON, error) {
	fake.setJSONMutex.Lock()
	ret, specificReturn := fake.setJSONReturnsOnCall[len(fake.setJSONArgsForCall)]
	fake.setJSONArgsForCall = append(fake.setJSONArgsForCall, struct {
		name  string
		value values.JSON
	}{name, value})
	fake.recordInvocation("SetJSON", []interface{}{name, value})
	fake.setJSONMutex.Unlock()
	if fake.SetJSONStub != nil {
		return fake.SetJSONStub(name, value)
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

func (fake *FakeCredhub) SetJSONArgsForCall(i int) (string, values.JSON) {
	fake.setJSONMutex.RLock()
	defer fake.setJSONMutex.RUnlock()
	return fake.setJSONArgsForCall[i].name, fake.setJSONArgsForCall[i].value
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

func (fake *FakeCredhub) SetValue(name string, value values.Value) (credentials.Value, error) {
	fake.setValueMutex.Lock()
	ret, specificReturn := fake.setValueReturnsOnCall[len(fake.setValueArgsForCall)]
	fake.setValueArgsForCall = append(fake.setValueArgsForCall, struct {
		name  string
		value values.Value
	}{name, value})
	fake.recordInvocation("SetValue", []interface{}{name, value})
	fake.setValueMutex.Unlock()
	if fake.SetValueStub != nil {
		return fake.SetValueStub(name, value)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.setValueReturns.result1, fake.setValueReturns.result2
}

func (fake *FakeCredhub) SetValueCallCount() int {
	fake.setValueMutex.RLock()
	defer fake.setValueMutex.RUnlock()
	return len(fake.setValueArgsForCall)
}

func (fake *FakeCredhub) SetValueArgsForCall(i int) (string, values.Value) {
	fake.setValueMutex.RLock()
	defer fake.setValueMutex.RUnlock()
	return fake.setValueArgsForCall[i].name, fake.setValueArgsForCall[i].value
}

func (fake *FakeCredhub) SetValueReturns(result1 credentials.Value, result2 error) {
	fake.SetValueStub = nil
	fake.setValueReturns = struct {
		result1 credentials.Value
		result2 error
	}{result1, result2}
}

func (fake *FakeCredhub) SetValueReturnsOnCall(i int, result1 credentials.Value, result2 error) {
	fake.SetValueStub = nil
	if fake.setValueReturnsOnCall == nil {
		fake.setValueReturnsOnCall = make(map[int]struct {
			result1 credentials.Value
			result2 error
		})
	}
	fake.setValueReturnsOnCall[i] = struct {
		result1 credentials.Value
		result2 error
	}{result1, result2}
}

func (fake *FakeCredhub) GetLatestValue(name string) (credentials.Value, error) {
	fake.getLatestValueMutex.Lock()
	ret, specificReturn := fake.getLatestValueReturnsOnCall[len(fake.getLatestValueArgsForCall)]
	fake.getLatestValueArgsForCall = append(fake.getLatestValueArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("GetLatestValue", []interface{}{name})
	fake.getLatestValueMutex.Unlock()
	if fake.GetLatestValueStub != nil {
		return fake.GetLatestValueStub(name)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getLatestValueReturns.result1, fake.getLatestValueReturns.result2
}

func (fake *FakeCredhub) GetLatestValueCallCount() int {
	fake.getLatestValueMutex.RLock()
	defer fake.getLatestValueMutex.RUnlock()
	return len(fake.getLatestValueArgsForCall)
}

func (fake *FakeCredhub) GetLatestValueArgsForCall(i int) string {
	fake.getLatestValueMutex.RLock()
	defer fake.getLatestValueMutex.RUnlock()
	return fake.getLatestValueArgsForCall[i].name
}

func (fake *FakeCredhub) GetLatestValueReturns(result1 credentials.Value, result2 error) {
	fake.GetLatestValueStub = nil
	fake.getLatestValueReturns = struct {
		result1 credentials.Value
		result2 error
	}{result1, result2}
}

func (fake *FakeCredhub) GetLatestValueReturnsOnCall(i int, result1 credentials.Value, result2 error) {
	fake.GetLatestValueStub = nil
	if fake.getLatestValueReturnsOnCall == nil {
		fake.getLatestValueReturnsOnCall = make(map[int]struct {
			result1 credentials.Value
			result2 error
		})
	}
	fake.getLatestValueReturnsOnCall[i] = struct {
		result1 credentials.Value
		result2 error
	}{result1, result2}
}

func (fake *FakeCredhub) FindByPath(path string) (credentials.FindResults, error) {
	fake.findByPathMutex.Lock()
	ret, specificReturn := fake.findByPathReturnsOnCall[len(fake.findByPathArgsForCall)]
	fake.findByPathArgsForCall = append(fake.findByPathArgsForCall, struct {
		path string
	}{path})
	fake.recordInvocation("FindByPath", []interface{}{path})
	fake.findByPathMutex.Unlock()
	if fake.FindByPathStub != nil {
		return fake.FindByPathStub(path)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.findByPathReturns.result1, fake.findByPathReturns.result2
}

func (fake *FakeCredhub) FindByPathCallCount() int {
	fake.findByPathMutex.RLock()
	defer fake.findByPathMutex.RUnlock()
	return len(fake.findByPathArgsForCall)
}

func (fake *FakeCredhub) FindByPathArgsForCall(i int) string {
	fake.findByPathMutex.RLock()
	defer fake.findByPathMutex.RUnlock()
	return fake.findByPathArgsForCall[i].path
}

func (fake *FakeCredhub) FindByPathReturns(result1 credentials.FindResults, result2 error) {
	fake.FindByPathStub = nil
	fake.findByPathReturns = struct {
		result1 credentials.FindResults
		result2 error
	}{result1, result2}
}

func (fake *FakeCredhub) FindByPathReturnsOnCall(i int, result1 credentials.FindResults, result2 error) {
	fake.FindByPathStub = nil
	if fake.findByPathReturnsOnCall == nil {
		fake.findByPathReturnsOnCall = make(map[int]struct {
			result1 credentials.FindResults
			result2 error
		})
	}
	fake.findByPathReturnsOnCall[i] = struct {
		result1 credentials.FindResults
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
	fake.setValueMutex.RLock()
	defer fake.setValueMutex.RUnlock()
	fake.getLatestValueMutex.RLock()
	defer fake.getLatestValueMutex.RUnlock()
	fake.findByPathMutex.RLock()
	defer fake.findByPathMutex.RUnlock()
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
