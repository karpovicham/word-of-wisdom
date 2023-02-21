package resolver

import (
	"testing"

	"github.com/karpovicham/word-of-wisdom/internal/messenger"

	"github.com/gojuno/minimock/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "package resolver")
}

var _ = Describe("resolver tests", func() {
	var (
		mc *minimock.Controller

		msgrMock *messenger.MessengerMock
		res      Resolver
	)

	BeforeEach(func() {
		mc = minimock.NewController(GinkgoT())
		msgrMock = messenger.NewMessengerMock(mc)
	})

	AfterEach(func() {
		mc.Finish()
	})

	Context("NewClientAPIResolver func", func() {
		When("parameters are valid", func() {
			BeforeEach(func() {
				res = &resolver{
					Msgr: msgrMock,
				}
			})

			It("should return client struct", func() {
				resolver := NewClientAPIResolver(msgrMock)
				Ω(resolver).ShouldNot(BeIdenticalTo(res))
			})
		})
	})
})
