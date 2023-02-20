package hashcash

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"strconv"
	"time"

	"github.com/gojuno/minimock/v3"
	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("hashcash tests", func() {
	var (
		mc              *minimock.Controller
		leadingZeroBits uint
		resource        string
		fuzzer          = fuzz.New().NilChance(0)
		ctx             context.Context
	)

	BeforeEach(func() {
		mc = minimock.NewController(GinkgoT())
		ctx = context.Background()
	})

	AfterEach(func() {
		mc.Finish()
	})

	Context("NewHashcashData func", func() {
		When("leadingZeroBits and resource are set", func() {
			BeforeEach(func() {
				fuzzer.Fuzz(&leadingZeroBits)
				fuzzer.Fuzz(&resource)
			})

			It("should return valid hashcash data", func() {
				data := NewHashcashData(leadingZeroBits, resource)
				Ω(data).ShouldNot(BeNil())
				Ω(data.Version).Should(Equal(hashcashFormatVersion))
				Ω(data.Bits).Should(Equal(leadingZeroBits))
				Ω(data.Date).ShouldNot(BeEmpty())
				Ω(data.Resource).Should(Equal(resource))
				Ω(data.Random).ShouldNot(BeEmpty())
				Ω(data.Counter).Should(BeZero())
			})
		})
	})

	Context("data.Header func", func() {
		When("data is set", func() {
			var data *data

			BeforeEach(func() {
				fuzzer.Fuzz(&data)
			})

			It("should return valid hashcash header", func() {
				expHeader := fmt.Sprintf(
					headerFormat,
					hashcashFormatVersion,
					data.Bits,
					data.Date,
					data.Resource,
					data.Random,
					base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(data.Counter, 10))),
				)

				header := data.Header()
				Ω(header).Should(Equal(expHeader))
			})
		})
	})

	Context("data.ComputeData func", func() {
		var computeTimeout time.Duration
		var data *data
		var hasher hash.Hash

		BeforeEach(func() {
			fuzzer.Fuzz(&data)
			hasher = sha256.New()
		})

		When("data leadingZeroBits is out of hasher range", func() {
			BeforeEach(func() {
				data.Bits = 257
				computeTimeout = 10 * time.Second
			})

			It("should fail", func() {
				_, err := data.ComputeData(ctx, hasher, computeTimeout)
				Ω(err).Should(MatchError("bits config 257 is out of the hash range (256 bits)"))
			})
		})

		When("data process it out of time duration is out of hasher range", func() {
			BeforeEach(func() {
				// Many bits and short timeout
				data.Bits = 50
				computeTimeout = time.Second
			})

			It("should fail with timeout error", func() {
				_, err := data.ComputeData(ctx, hasher, computeTimeout)
				Ω(err).Should(MatchError(ErrComputeTimeout))
			})
		})

		When("context is canceled", func() {
			var computeDataCtx context.Context
			var computeDataCancel context.CancelFunc
			BeforeEach(func() {
				// Many bits and long timeout
				data.Bits = 50
				computeTimeout = 10 * time.Second
				computeDataCtx, computeDataCancel = context.WithCancel(ctx)
			})

			It("should fail with canceled context error", func() {
				computeDataCancel()
				_, err := data.ComputeData(computeDataCtx, hasher, computeTimeout)
				Ω(err).Should(MatchError(ErrClosedContext))
			})
		})

		When("params are valid and the process should take short time", func() {
			BeforeEach(func() {
				// A little bits and long timeout
				data.Bits = 5
				computeTimeout = 10 * time.Second
			})

			It("should fail with canceled context error", func() {
				_, err := data.ComputeData(ctx, hasher, computeTimeout)
				Ω(err).ShouldNot(HaveOccurred())
			})

			When("Data is computed", func() {
				BeforeEach(func() {
					// Many bits and short time
					data.Bits = 5
					computeTimeout = 10 * time.Second
				})

				It("should retunr data with ony counter changed", func() {
					resData, err := data.ComputeData(ctx, hasher, computeTimeout)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(resData.Counter).ShouldNot(BeZero())

					// Only Counter should be changed
					data.Counter = resData.Counter
					Ω(resData).Should(Equal(data))
				})
			})
		})
	})
})
