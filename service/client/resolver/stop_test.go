package resolver

import (
	"errors"

	"github.com/karpovicham/word-of-wisdom/internal/messenger"
	"github.com/karpovicham/word-of-wisdom/internal/proto"

	"github.com/gojuno/minimock/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stop tests", func() {
	var (
		mc *minimock.Controller

		msgrMock *messenger.MessengerMock
		r        Resolver
		reqMsg   *proto.Message
		testErr  = errors.New("test error")
	)

	BeforeEach(func() {
		mc = minimock.NewController(GinkgoT())
		msgrMock = messenger.NewMessengerMock(mc)
		r = NewClientAPIResolver(msgrMock)
	})

	AfterEach(func() {
		mc.Finish()
	})

	Context("Stop func", func() {
		It("should call msgr.Send", func() {
			reqMsg = &proto.Message{
				Type: proto.TypeStop,
				Data: nil,
			}

			msgrMock.SendMock.
				Expect(reqMsg).
				Return(testErr)

			_ = r.Stop()
		})

		When("msgr.Send return error", func() {
			BeforeEach(func() {
				msgrMock.SendMock.
					Return(testErr)
			})

			It("should return error", func() {
				err := r.Stop()
				Ω(err).Should(HaveOccurred())
			})
		})

		When("msgr.Send succeed", func() {
			BeforeEach(func() {
				msgrMock.SendMock.
					Return(nil)
			})

			It("should call msgr.Receive", func() {
				err := r.Stop()
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
