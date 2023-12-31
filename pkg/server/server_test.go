package server_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/dgruber/qsub/pkg/server"
)

var _ = Describe("Server", func() {

	Context("Secrets", func() {

		It("should create a secret", func() {
			secret, err := server.GetOrCreateSecret()
			Expect(err).To(BeNil())
			Expect(secret).ToNot(BeNil())

			// delete secret in home
			err = server.DeleteSecret()
			Expect(err).To(BeNil())

			// create secret again
			secret2, err := server.GetOrCreateSecret()
			Expect(err).To(BeNil())
			Expect(secret2).ToNot(BeNil())
			// secret should be different
			Expect(secret2).NotTo(Equal(secret))

			secret3, err := server.GetOrCreateSecret()
			Expect(err).To(BeNil())
			Expect(secret3).ToNot(BeNil())
			// secret should be the same
			Expect(secret3).To(Equal(secret2))

			server.DeleteSecret()
		})

	})

})
