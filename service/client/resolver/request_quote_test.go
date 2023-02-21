package resolver

import (
	"encoding/json"
	"errors"

	"github.com/karpovicham/word-of-wisdom/internal/messenger"
	"github.com/karpovicham/word-of-wisdom/internal/proto"
	"github.com/karpovicham/word-of-wisdom/pkg/pow"
	"github.com/karpovicham/word-of-wisdom/service/quotes_book"

	"github.com/gojuno/minimock/v3"
	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestQuote tests", func() {
	var (
		mc *minimock.Controller

		msgrMock   *messenger.MessengerMock
		r          Resolver
		reqPOWData pow.Data
		reqMsg     *proto.Message
		respMsg    *proto.Message
		testErr    = errors.New("test error")

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

	Context("RequestQuote func", func() {
		It("should call msgr.Send", func() {
			reqMsg = &proto.Message{
				Type: proto.TypeQuote,
				Data: reqPOWData,
			}

			msgrMock.SendMock.
				Expect(reqMsg).
				Return(testErr)

			_, _ = r.RequestQuote(reqPOWData)
		})

		When("msgr.Send return error", func() {
			BeforeEach(func() {
				msgrMock.SendMock.
					Return(testErr)
			})

			It("should return error", func() {
				_, err := r.RequestQuote(reqPOWData)
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

				_, _ = r.RequestQuote(reqPOWData)
			})

			When("msgr.Receive returns error", func() {
				BeforeEach(func() {
					msgrMock.ReceiveMock.
						Return(nil, testErr)
				})

				It("should return error", func() {
					_, err := r.RequestQuote(reqPOWData)
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

						_, err := r.RequestQuote(reqPOWData)
						Ω(err).Should(HaveOccurred())
					})
				})

				When("message returns correct type", func() {
					BeforeEach(func() {
						respMsg.Type = proto.TypeQuote
					})

					When("message returns valid data error", func() {
						It("should return error", func() {
							respMsg.Error = proto.ErrorPtr(proto.ErrorInvalidData)

							_, err := r.RequestQuote(reqPOWData)
							Ω(err).Should(MatchError(ErrInvalidReqData))
						})
					})

					When("message returns not verified error", func() {
						It("should return error", func() {
							respMsg.Error = proto.ErrorPtr(proto.ErrorNotVerified)

							_, err := r.RequestQuote(reqPOWData)
							Ω(err).Should(MatchError(ErrNotVerified))
						})
					})

					When("message returns unknown error", func() {
						It("should return error", func() {
							respMsg.Error = proto.ErrorPtr(proto.Error(100))

							_, err := r.RequestQuote(reqPOWData)
							Ω(err).Should(MatchError(ErrUnknownRespError))
						})
					})

					When("msgr.Receive returns no error", func() {
						BeforeEach(func() {
							respMsg.Error = nil
						})

						When("message data is not quote JSON", func() {
							It("should return error", func() {
								_, err := r.RequestQuote(reqPOWData)
								Ω(err).Should(MatchError(ErrInvalidRespData))
							})
						})

						When("message data is not quote JSON", func() {
							var quote *quotes_book.Quote

							BeforeEach(func() {
								fuzzer.Fuzz(&quote)
								quoteData, _ := json.Marshal(quote)
								respMsg.Data = quoteData
							})

							It("should return quote", func() {
								data, err := r.RequestQuote(reqPOWData)
								Ω(err).ShouldNot(HaveOccurred())
								Ω(data).Should(Equal(quote))
							})
						})
					})
				})
			})
		})
	})
})
