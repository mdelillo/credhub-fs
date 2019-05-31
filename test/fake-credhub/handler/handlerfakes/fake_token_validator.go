// Code generated by counterfeiter. DO NOT EDIT.
package handlerfakes

import (
	"sync"
)

type FakeTokenValidator struct {
	ValidateTokenWithClaimsStub        func(string, map[string]string) error
	validateTokenWithClaimsMutex       sync.RWMutex
	validateTokenWithClaimsArgsForCall []struct {
		arg1 string
		arg2 map[string]string
	}
	validateTokenWithClaimsReturns struct {
		result1 error
	}
	validateTokenWithClaimsReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTokenValidator) ValidateTokenWithClaims(arg1 string, arg2 map[string]string) error {
	fake.validateTokenWithClaimsMutex.Lock()
	ret, specificReturn := fake.validateTokenWithClaimsReturnsOnCall[len(fake.validateTokenWithClaimsArgsForCall)]
	fake.validateTokenWithClaimsArgsForCall = append(fake.validateTokenWithClaimsArgsForCall, struct {
		arg1 string
		arg2 map[string]string
	}{arg1, arg2})
	fake.recordInvocation("ValidateTokenWithClaims", []interface{}{arg1, arg2})
	fake.validateTokenWithClaimsMutex.Unlock()
	if fake.ValidateTokenWithClaimsStub != nil {
		return fake.ValidateTokenWithClaimsStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.validateTokenWithClaimsReturns
	return fakeReturns.result1
}

func (fake *FakeTokenValidator) ValidateTokenWithClaimsCallCount() int {
	fake.validateTokenWithClaimsMutex.RLock()
	defer fake.validateTokenWithClaimsMutex.RUnlock()
	return len(fake.validateTokenWithClaimsArgsForCall)
}

func (fake *FakeTokenValidator) ValidateTokenWithClaimsCalls(stub func(string, map[string]string) error) {
	fake.validateTokenWithClaimsMutex.Lock()
	defer fake.validateTokenWithClaimsMutex.Unlock()
	fake.ValidateTokenWithClaimsStub = stub
}

func (fake *FakeTokenValidator) ValidateTokenWithClaimsArgsForCall(i int) (string, map[string]string) {
	fake.validateTokenWithClaimsMutex.RLock()
	defer fake.validateTokenWithClaimsMutex.RUnlock()
	argsForCall := fake.validateTokenWithClaimsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeTokenValidator) ValidateTokenWithClaimsReturns(result1 error) {
	fake.validateTokenWithClaimsMutex.Lock()
	defer fake.validateTokenWithClaimsMutex.Unlock()
	fake.ValidateTokenWithClaimsStub = nil
	fake.validateTokenWithClaimsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTokenValidator) ValidateTokenWithClaimsReturnsOnCall(i int, result1 error) {
	fake.validateTokenWithClaimsMutex.Lock()
	defer fake.validateTokenWithClaimsMutex.Unlock()
	fake.ValidateTokenWithClaimsStub = nil
	if fake.validateTokenWithClaimsReturnsOnCall == nil {
		fake.validateTokenWithClaimsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validateTokenWithClaimsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTokenValidator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.validateTokenWithClaimsMutex.RLock()
	defer fake.validateTokenWithClaimsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTokenValidator) recordInvocation(key string, args []interface{}) {
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
