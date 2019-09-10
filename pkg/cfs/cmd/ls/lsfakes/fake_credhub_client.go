// Code generated by counterfeiter. DO NOT EDIT.
package lsfakes

import (
	"sync"

	"github.com/mdelillo/credhub-fs/pkg/credhub"
)

type FakeCredhubClient struct {
	DeleteCredentialByNameStub        func(string) error
	deleteCredentialByNameMutex       sync.RWMutex
	deleteCredentialByNameArgsForCall []struct {
		arg1 string
	}
	deleteCredentialByNameReturns struct {
		result1 error
	}
	deleteCredentialByNameReturnsOnCall map[int]struct {
		result1 error
	}
	FindCredentialsByPathStub        func(string) ([]credhub.Credential, error)
	findCredentialsByPathMutex       sync.RWMutex
	findCredentialsByPathArgsForCall []struct {
		arg1 string
	}
	findCredentialsByPathReturns struct {
		result1 []credhub.Credential
		result2 error
	}
	findCredentialsByPathReturnsOnCall map[int]struct {
		result1 []credhub.Credential
		result2 error
	}
	GetCredentialByNameStub        func(string) (credhub.Credential, error)
	getCredentialByNameMutex       sync.RWMutex
	getCredentialByNameArgsForCall []struct {
		arg1 string
	}
	getCredentialByNameReturns struct {
		result1 credhub.Credential
		result2 error
	}
	getCredentialByNameReturnsOnCall map[int]struct {
		result1 credhub.Credential
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCredhubClient) DeleteCredentialByName(arg1 string) error {
	fake.deleteCredentialByNameMutex.Lock()
	ret, specificReturn := fake.deleteCredentialByNameReturnsOnCall[len(fake.deleteCredentialByNameArgsForCall)]
	fake.deleteCredentialByNameArgsForCall = append(fake.deleteCredentialByNameArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("DeleteCredentialByName", []interface{}{arg1})
	fake.deleteCredentialByNameMutex.Unlock()
	if fake.DeleteCredentialByNameStub != nil {
		return fake.DeleteCredentialByNameStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deleteCredentialByNameReturns
	return fakeReturns.result1
}

func (fake *FakeCredhubClient) DeleteCredentialByNameCallCount() int {
	fake.deleteCredentialByNameMutex.RLock()
	defer fake.deleteCredentialByNameMutex.RUnlock()
	return len(fake.deleteCredentialByNameArgsForCall)
}

func (fake *FakeCredhubClient) DeleteCredentialByNameCalls(stub func(string) error) {
	fake.deleteCredentialByNameMutex.Lock()
	defer fake.deleteCredentialByNameMutex.Unlock()
	fake.DeleteCredentialByNameStub = stub
}

func (fake *FakeCredhubClient) DeleteCredentialByNameArgsForCall(i int) string {
	fake.deleteCredentialByNameMutex.RLock()
	defer fake.deleteCredentialByNameMutex.RUnlock()
	argsForCall := fake.deleteCredentialByNameArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCredhubClient) DeleteCredentialByNameReturns(result1 error) {
	fake.deleteCredentialByNameMutex.Lock()
	defer fake.deleteCredentialByNameMutex.Unlock()
	fake.DeleteCredentialByNameStub = nil
	fake.deleteCredentialByNameReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeCredhubClient) DeleteCredentialByNameReturnsOnCall(i int, result1 error) {
	fake.deleteCredentialByNameMutex.Lock()
	defer fake.deleteCredentialByNameMutex.Unlock()
	fake.DeleteCredentialByNameStub = nil
	if fake.deleteCredentialByNameReturnsOnCall == nil {
		fake.deleteCredentialByNameReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteCredentialByNameReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeCredhubClient) FindCredentialsByPath(arg1 string) ([]credhub.Credential, error) {
	fake.findCredentialsByPathMutex.Lock()
	ret, specificReturn := fake.findCredentialsByPathReturnsOnCall[len(fake.findCredentialsByPathArgsForCall)]
	fake.findCredentialsByPathArgsForCall = append(fake.findCredentialsByPathArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("FindCredentialsByPath", []interface{}{arg1})
	fake.findCredentialsByPathMutex.Unlock()
	if fake.FindCredentialsByPathStub != nil {
		return fake.FindCredentialsByPathStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.findCredentialsByPathReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCredhubClient) FindCredentialsByPathCallCount() int {
	fake.findCredentialsByPathMutex.RLock()
	defer fake.findCredentialsByPathMutex.RUnlock()
	return len(fake.findCredentialsByPathArgsForCall)
}

func (fake *FakeCredhubClient) FindCredentialsByPathCalls(stub func(string) ([]credhub.Credential, error)) {
	fake.findCredentialsByPathMutex.Lock()
	defer fake.findCredentialsByPathMutex.Unlock()
	fake.FindCredentialsByPathStub = stub
}

func (fake *FakeCredhubClient) FindCredentialsByPathArgsForCall(i int) string {
	fake.findCredentialsByPathMutex.RLock()
	defer fake.findCredentialsByPathMutex.RUnlock()
	argsForCall := fake.findCredentialsByPathArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCredhubClient) FindCredentialsByPathReturns(result1 []credhub.Credential, result2 error) {
	fake.findCredentialsByPathMutex.Lock()
	defer fake.findCredentialsByPathMutex.Unlock()
	fake.FindCredentialsByPathStub = nil
	fake.findCredentialsByPathReturns = struct {
		result1 []credhub.Credential
		result2 error
	}{result1, result2}
}

func (fake *FakeCredhubClient) FindCredentialsByPathReturnsOnCall(i int, result1 []credhub.Credential, result2 error) {
	fake.findCredentialsByPathMutex.Lock()
	defer fake.findCredentialsByPathMutex.Unlock()
	fake.FindCredentialsByPathStub = nil
	if fake.findCredentialsByPathReturnsOnCall == nil {
		fake.findCredentialsByPathReturnsOnCall = make(map[int]struct {
			result1 []credhub.Credential
			result2 error
		})
	}
	fake.findCredentialsByPathReturnsOnCall[i] = struct {
		result1 []credhub.Credential
		result2 error
	}{result1, result2}
}

func (fake *FakeCredhubClient) GetCredentialByName(arg1 string) (credhub.Credential, error) {
	fake.getCredentialByNameMutex.Lock()
	ret, specificReturn := fake.getCredentialByNameReturnsOnCall[len(fake.getCredentialByNameArgsForCall)]
	fake.getCredentialByNameArgsForCall = append(fake.getCredentialByNameArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetCredentialByName", []interface{}{arg1})
	fake.getCredentialByNameMutex.Unlock()
	if fake.GetCredentialByNameStub != nil {
		return fake.GetCredentialByNameStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getCredentialByNameReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCredhubClient) GetCredentialByNameCallCount() int {
	fake.getCredentialByNameMutex.RLock()
	defer fake.getCredentialByNameMutex.RUnlock()
	return len(fake.getCredentialByNameArgsForCall)
}

func (fake *FakeCredhubClient) GetCredentialByNameCalls(stub func(string) (credhub.Credential, error)) {
	fake.getCredentialByNameMutex.Lock()
	defer fake.getCredentialByNameMutex.Unlock()
	fake.GetCredentialByNameStub = stub
}

func (fake *FakeCredhubClient) GetCredentialByNameArgsForCall(i int) string {
	fake.getCredentialByNameMutex.RLock()
	defer fake.getCredentialByNameMutex.RUnlock()
	argsForCall := fake.getCredentialByNameArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCredhubClient) GetCredentialByNameReturns(result1 credhub.Credential, result2 error) {
	fake.getCredentialByNameMutex.Lock()
	defer fake.getCredentialByNameMutex.Unlock()
	fake.GetCredentialByNameStub = nil
	fake.getCredentialByNameReturns = struct {
		result1 credhub.Credential
		result2 error
	}{result1, result2}
}

func (fake *FakeCredhubClient) GetCredentialByNameReturnsOnCall(i int, result1 credhub.Credential, result2 error) {
	fake.getCredentialByNameMutex.Lock()
	defer fake.getCredentialByNameMutex.Unlock()
	fake.GetCredentialByNameStub = nil
	if fake.getCredentialByNameReturnsOnCall == nil {
		fake.getCredentialByNameReturnsOnCall = make(map[int]struct {
			result1 credhub.Credential
			result2 error
		})
	}
	fake.getCredentialByNameReturnsOnCall[i] = struct {
		result1 credhub.Credential
		result2 error
	}{result1, result2}
}

func (fake *FakeCredhubClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deleteCredentialByNameMutex.RLock()
	defer fake.deleteCredentialByNameMutex.RUnlock()
	fake.findCredentialsByPathMutex.RLock()
	defer fake.findCredentialsByPathMutex.RUnlock()
	fake.getCredentialByNameMutex.RLock()
	defer fake.getCredentialByNameMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCredhubClient) recordInvocation(key string, args []interface{}) {
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
