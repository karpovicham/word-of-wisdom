package client

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/karpovicham/word-of-wisdom/internal/logger"
	"github.com/karpovicham/word-of-wisdom/pkg/pow"
	"github.com/karpovicham/word-of-wisdom/service/client/resolver"
	"github.com/karpovicham/word-of-wisdom/service/quotes_book"

	"github.com/gojuno/minimock/v3"
	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "package client")
}

var _ = Describe("client tests", func() {
	var (
		mc  *minimock.Controller
		ctx context.Context

		cfg           Config
		log           logger.Logger
		powWorkerMock *pow.ClientWorkerMock

		client *Client

		fuzzer  = fuzz.New().NilChance(0)
		testErr = errors.New("test error")
	)

	BeforeEach(func() {
		mc = minimock.NewController(GinkgoT())
		ctx = context.Background()
		log = logger.NewLogger(os.Stdout)

		powWorkerMock = pow.NewClientWorkerMock(mc)

		// Test conn on real port
		cfg = Config{
			Host: "",
			Port: "1234",
		}
	})

	AfterEach(func() {
		mc.Finish()
	})

	Context("NewClient func", func() {
		When("parameters are valid", func() {
			BeforeEach(func() {
				client = &Client{
					Cfg:       cfg,
					Log:       log,
					POWWorker: powWorkerMock,
				}
			})

			It("should return client struct", func() {
				c := NewClient(cfg, log, powWorkerMock)
				Ω(c).ShouldNot(BeIdenticalTo(client))
			})
		})
	})

	Context("Client.ProcessQuote func", func() {
		var (
			challengeRespPowData pow.Data
			doWorkRespPowData    pow.Data
			resolverMock         *resolver.ResolverMock
		)

		BeforeEach(func() {
			client = NewClient(cfg, log, powWorkerMock)
			resolverMock = resolver.NewResolverMock(mc)
		})

		It("should call resolver.RequestChallenge func", func() {
			resolverMock.RequestChallengeMock.
				Expect().
				Return(nil, testErr)

			_ = client.ProcessQuote(ctx, resolverMock)
		})

		When("resolver.RequestChallenge return error", func() {
			BeforeEach(func() {
				resolverMock.RequestChallengeMock.
					Return(nil, testErr)
			})

			It("should return error", func() {
				err := client.ProcessQuote(ctx, resolverMock)
				Ω(err).Should(MatchError(testErr))
			})

			When("resolver.RequestChallenge returns success response", func() {
				BeforeEach(func() {
					fuzzer.NumElements(2, 4).Fuzz(&challengeRespPowData)

					resolverMock.RequestChallengeMock.
						Return(challengeRespPowData, nil)
				})

				It("should call powWorker.DoWork func", func() {
					powWorkerMock.DoWorkMock.
						Expect(ctx, challengeRespPowData).
						Return(nil, testErr)

					_ = client.ProcessQuote(ctx, resolverMock)
				})

				When("powWorker.DoWork return error response", func() {
					BeforeEach(func() {
						powWorkerMock.DoWorkMock.
							Return(nil, testErr)
					})

					It("should return error", func() {
						err := client.ProcessQuote(ctx, resolverMock)
						Ω(err).Should(MatchError(testErr))
					})
				})

				When("powWorker.DoWork return success response", func() {
					BeforeEach(func() {
						fuzzer.NumElements(2, 4).Fuzz(&doWorkRespPowData)

						powWorkerMock.DoWorkMock.
							Return(doWorkRespPowData, nil)
					})

					It("should call resolver.RequestQuote", func() {
						resolverMock.RequestQuoteMock.
							Expect(doWorkRespPowData).
							Return(nil, testErr)

						_ = client.ProcessQuote(ctx, resolverMock)
					})

					When("resolver.RequestQuote return error response", func() {
						BeforeEach(func() {
							resolverMock.RequestQuoteMock.
								Return(nil, testErr)
						})

						It("should return error", func() {
							err := client.ProcessQuote(ctx, resolverMock)
							Ω(err).Should(MatchError(testErr))
						})
					})

					When("resolver.RequestQuote return success response", func() {
						BeforeEach(func() {
							var quote quotes_book.Quote
							fuzzer.Fuzz(&quote)

							resolverMock.RequestQuoteMock.
								Return(&quote, nil)
						})

						It("should no error", func() {
							err := client.ProcessQuote(ctx, resolverMock)
							Ω(err).ShouldNot(HaveOccurred())
						})
					})
				})
			})
		})
	})
})
