package proto

import (
	"encoding/json"
	"testing"

	"github.com/gojuno/minimock/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "package proto")
}

var _ = Describe("proto tests", func() {
	var (
		mc  *minimock.Controller
		msg Message
	)

	BeforeEach(func() {
		mc = minimock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mc.Finish()
	})

	Context("Parse func", func() {
		var data []byte

		When("data is invalid json", func() {
			BeforeEach(func() {
				data = []byte("invalid json")
			})

			It("should return error", func() {
				_, err := Parse(data)
				Ω(err).Should(HaveOccurred())
			})
		})

		When("data is valid json", func() {
			BeforeEach(func() {

			})

			It("should return error", func() {
				msg = Message{
					Type:  TypeChallenge,
					Data:  []byte("some"),
					Error: ErrorPtr(ErrorNotVerified),
				}

				msgBytes, err := json.Marshal(msg)
				Ω(err).ShouldNot(HaveOccurred())

				_, err = Parse(msgBytes)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Context("ToJSON func", func() {
		It("should return json", func() {
			msg = Message{
				Type:  TypeChallenge,
				Data:  []byte("some"),
				Error: ErrorPtr(ErrorNotVerified),
			}

			msgBytes, err := json.Marshal(msg)
			Ω(err).ShouldNot(HaveOccurred())

			res := msg.ToJSON()
			Ω(res).Should(Equal(msgBytes))
		})
	})
})
