package credentials_test

import (
	"github.com/mdelillo/credhub-fs/test/fake-credhub/credentials"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store", func() {
	It("sets, gets, and deletes credentials", func() {
		cred1 := credentials.Credential{Name: "/cred1"}
		cred2 := credentials.Credential{Name: "/cred2"}

		store := credentials.NewStore()

		store.Set(cred1)
		store.Set(cred2)

		actualCred1, found := store.GetByName(cred1.Name)
		Expect(found).To(BeTrue())
		Expect(actualCred1).To(Equal(cred1))

		actualCred2, found := store.GetByName(cred2.Name)
		Expect(found).To(BeTrue())
		Expect(actualCred2).To(Equal(cred2))

		deleted := store.Delete(cred1.Name)
		Expect(deleted).To(BeTrue())

		_, found = store.GetByName(cred1.Name)
		Expect(found).To(BeFalse())
		_, found = store.GetByName(cred2.Name)
		Expect(found).To(BeTrue())

		deleted = store.Delete(cred1.Name)
		Expect(deleted).To(BeFalse())
	})

	Context("when the credential does not exist", func() {
		It("returns false", func() {
			store := credentials.NewStore()

			_, found := store.GetByName("some-nonexistent-name")
			Expect(found).To(BeFalse())
		})
	})

	Context("when a credential with a duplicate name is set", func() {
		It("overwrites the credential", func() {
			oldCred := credentials.Credential{Name: "/cred", Value: "old"}
			newCred := credentials.Credential{Name: "/cred", Value: "new"}

			store := credentials.NewStore()

			store.Set(oldCred)
			store.Set(newCred)

			actualCred, found := store.GetByName(oldCred.Name)
			Expect(found).To(BeTrue())
			Expect(actualCred).To(Equal(newCred))
		})
	})

	Describe("GetByPath", func() {
		It("gets credentials within a given path", func() {
			nestedCred1 := credentials.Credential{Name: "/nested/cred1"}
			nestedCred2 := credentials.Credential{Name: "/nested/deeply/cred2"}
			otherCred := credentials.Credential{Name: "/nested-lie"}

			store := credentials.NewStore()

			store.Set(nestedCred1)
			store.Set(nestedCred2)
			store.Set(otherCred)

			nestedCreds := store.GetByPath("/nested")
			Expect(nestedCreds).To(ConsistOf(nestedCred1, nestedCred2))

			nestedCreds = store.GetByPath("/nested/")
			Expect(nestedCreds).To(ConsistOf(nestedCred1, nestedCred2))
		})

		It("does not return credentials exactly matching the path", func() {
			cred := credentials.Credential{Name: "some-cred"}

			store := credentials.NewStore()
			store.Set(cred)

			creds := store.GetByPath(cred.Name)
			Expect(creds).To(BeEmpty())
		})
	})
})
