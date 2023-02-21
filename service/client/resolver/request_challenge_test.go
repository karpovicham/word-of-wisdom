package resolver

import (
	"errors"

	"github.com/karpovicham/word-of-wisdom/internal/messenger"
	"github.com/karpovicham/word-of-wisdom/internal/proto"
	"github.com/karpovicham/word-of-wisdom/pkg/pow"

	"github.com/gojuno/minimock/v3"
	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestChallenge tests", func() {
	var (
		mc *minimock.Controller

		msgrMock *messenger.MessengerMock
		r        Resolver
		reqMsg   *proto.Message
		respMsg  *proto.Message
		testErr  = errors.New("test error")

		fuzzer = fuzz.New().NilChance(0)
	)

	BeforeEach(func() {
		mc = minimock.NewController(GinkgoT())
		msgrMock = messenger.NewMessengerMock(mc)
		r = NewClientAPIResolver(msgrMock)
	})

	AfterEach(func() {
		mc.Finish()
	})

	Context("RequestChallenge func", func() {
		It("should call msgr.Send", func() {
			reqMsg = &proto.Message{
				Type: proto.TypeChallenge,
				Data: nil,
			}

			msgrMock.SendMock.
				Expect(reqMsg).
				Return(testErr)

			_, _ = r.RequestChallenge()
		})

		When("msgr.Send return error", func() {
			BeforeEach(func() {
				msgrMock.SendMock.
					Return(testErr)
			})

			It("should return error", func() {
				_, err := r.RequestChallenge()
				Ω(err).Should(HaveOccurred())
			})
		})

		When("msgr.Send succeed", func() {
			BeforeEach(func() {
				msgrMock.SendMock.
					Return(nil)
			})

			It("should call msgr.Receive", func() {
				msgrMock.ReceiveMock.
					Expect().
					Return(nil, testErr)

				_, _ = r.RequestChallenge()
			})

			When("msgr.Receive returns error", func() {
				BeforeEach(func() {
					msgrMock.ReceiveMock.
						Return(nil, testErr)
				})

				It("should return error", func() {
					_, err := r.RequestChallenge()
					Ω(err).Should(HaveOccurred())
				})
			})

			When("msgr.Receive succeeds", func() {
				BeforeEach(func() {
					respMsg = new(proto.Message)
					fuzzer.Fuzz(&respMsg.Data)
				})

				JustBeforeEach(func() {
					msgrMock.ReceiveMock.
						Return(respMsg, nil)
				})

				When("message contains returns wrong type", func() {
					It("should return error", func() {
						respMsg.Type = proto.TypeStop

						_, err := r.RequestChallenge()
						Ω(err).Should(HaveOccurred())
					})
				})

				When("message returns correct type", func() {
					BeforeEach(func() {
						respMsg.Type = proto.TypeChallenge
					})

					When("message returns valid data error", func() {
						It("should return error", func() {
							respMsg.Error = &proto.ErrorInvalidData

							_, err := r.RequestChallenge()
							Ω(err).Should(MatchError(ErrInvalidReqData))
						})
					})

					When("message returns not verified error", func() {
						It("should return error", func() {
							respMsg.Error = &proto.ErrorNotVerified

							_, err := r.RequestChallenge()
							Ω(err).Should(MatchError(ErrNotVerified))
						})
					})

					When("message returns unknown error", func() {
						It("should return error", func() {
							e := proto.Error(testErr)
							respMsg.Error = &e

							_, err := r.RequestChallenge()
							Ω(err).Should(MatchError(ErrUnknownRespError))
						})
					})

					When("msgr.Receive returns no error", func() {
						It("should return POW data", func() {
							respMsg.Error = nil

							data, err := r.RequestChallenge()
							Ω(err).ShouldNot(HaveOccurred())
							Ω(data).Should(Equal(pow.Data(respMsg.Data)))
						})
					})
				})
			})
		})
	})
})
