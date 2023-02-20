package in_memory_book

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/karpovicham/word-of-wisdom/service/quotes_book"

	"github.com/gojuno/minimock/v3"
	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "package in_memory_book")
}

var _ = Describe("in_memory_book tests", func() {
	var (
		mc             *minimock.Controller
		quotesFilepath string
		quotesBook     quotes_book.Book
		randSource     rand.Source
		fuzzer         = fuzz.New()
		ctx            context.Context
		service        *quotesBookService
	)

	BeforeEach(func() {
		mc = minimock.NewController(GinkgoT())
		randSource = rand.NewSource(time.Now().UnixMicro())
		ctx = context.Background()
	})

	AfterEach(func() {
		mc.Finish()
	})

	Context("NewQuotesBookService func", func() {
		When("no quotes file", func() {
			BeforeEach(func() {
				quotesFilepath = "unknown"
			})

			It("should fail with file not exists", func() {
				_, err := NewQuotesBookService(quotesFilepath, randSource)
				Ω(err).Should(MatchError(os.ErrNotExist))
			})
		})

		When("file with invalid format", func() {
			BeforeEach(func() {
				quotesFilepath = "./testdata/invalid_data_format.json"
			})

			It("should fail with invalid format data", func() {
				_, err := NewQuotesBookService(quotesFilepath, randSource)
				Ω(err).Should(HaveOccurred())
			})
		})

		When("file with no quotes", func() {
			BeforeEach(func() {
				quotesFilepath = "./testdata/no_quotes.json"
			})

			It("should fail with no quotes error", func() {
				_, err := NewQuotesBookService(quotesFilepath, randSource)
				Ω(err).Should(MatchError(ErrNoQuotes))
			})
		})

		When("file with quotes provided", func() {
			BeforeEach(func() {
				quotesFilepath = "./testdata/quotes.json"
			})

			It("should not fail", func() {
				_, err := NewQuotesBookService(quotesFilepath, randSource)
				Ω(err).ShouldNot(HaveOccurred())
			})

			When("quotes file is loaded", func() {
				BeforeEach(func() {
					quotesBook = quotes_book.Book{
						Quotes: []quotes_book.Quote{
							{
								Quote:  "Life isn’t about getting and having, it’s about giving and being.",
								Author: "Kevin Kruse",
							},
							{
								Quote:  "Whatever the mind of man can conceive and believe, it can achieve.",
								Author: "Napoleon Hill",
							},
						},
					}

					service = &quotesBookService{
						book:       quotesBook,
						randomizer: rand.New(randSource),
					}
				})

				It("should return service with inserted quotes and init randomizer", func() {
					resService, err := NewQuotesBookService(quotesFilepath, randSource)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(resService).ShouldNot(BeNil())
					Ω(resService).Should(BeEquivalentTo(service))
				})
			})
		})
	})

	Context("GetRandomQuote func", func() {
		var quote quotes_book.Quote

		When("quotes exists", func() {
			BeforeEach(func() {
				fuzzer.NilChance(0).Fuzz(&quote)
				service = &quotesBookService{
					book:       quotes_book.Book{Quotes: []quotes_book.Quote{quote}},
					randomizer: rand.New(randSource),
				}
			})

			It("should return quote", func() {
				resQuote, err := service.GetRandomQuote(ctx)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(resQuote).Should(Equal(quote))
			})
		})
	})
})
