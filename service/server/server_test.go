package server

import (
	"context"
	"errors"
	"net"
	"os"
	"testing"

	"github.com/karpovicham/word-of-wisdom/internal/logger"
	"github.com/karpovicham/word-of-wisdom/internal/messenger"
	"github.com/karpovicham/word-of-wisdom/internal/proto"
	"github.com/karpovicham/word-of-wisdom/pkg/pow"
	"github.com/karpovicham/word-of-wisdom/service/quotes_book"

	"github.com/gojuno/minimock/v3"
	fuzz "github.com/google/gofuzz"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "package server")
}

var _ = Describe("server tests", func() {
	var (
		mc        *minimock.Controller
		ctx       context.Context
		ctxCancel context.CancelFunc

		cfg            Config
		log            logger.Logger
		quotesBookMock *quotes_book.QuotesBookMock
		powWorkerMock  *pow.ServerWorkerMock
		msgrMock       *messenger.MessengerMock
		services       Services
		msgrFn         messenger.MsgrFn

		server *Server

		fuzzer  = fuzz.New().NilChance(0)
		testErr = errors.New("test error")
	)

	BeforeEach(func() {
		mc = minimock.NewController(GinkgoT())
		ctx, ctxCancel = context.WithCancel(context.Background())

		log = logger.NewLogger(os.Stdout)

		quotesBookMock = quotes_book.NewQuotesBookMock(mc)
		powWorkerMock = pow.NewServerWorkerMock(mc)
		msgrMock = messenger.NewMessengerMock(mc)

		msgrFn = func(conn net.Conn) messenger.Messenger {
			return msgrMock
		}

		services = Services{
			QuotesBook: quotesBookMock,
		}

		// Test conn on real port
		cfg = Config{
			Host: "",
			Port: "1234",
		}
	})

	AfterEach(func() {
		mc.Finish()
	})

	Context("NewTCPServer func", func() {
		When("parameters are valid", func() {
			BeforeEach(func() {
				server = &Server{
					Cfg:       cfg,
					Log:       log,
					Services:  services,
					POWWorker: powWorkerMock,
					MsgrFn:    msgrFn,
				}
			})

			It("should return server struct", func() {
				s := NewTCPServer(cfg, log, services, powWorkerMock, msgrFn)
				Ω(s).ShouldNot(BeIdenticalTo(server))
			})
		})
	})

	Context("NewTCPServer.Run func", func() {
		When("context is canceled", func() {
			BeforeEach(func() {
				server = NewTCPServer(cfg, log, services, powWorkerMock, msgrFn)
			})

			It("should stop listener and return no error", func() {
				ctxCancel()
				err := server.Run(ctx)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Context("NewTCPServer.handleConnection func", func() {
		When("client conn is set up", func() {
			var (
				clientConnMock  net.Conn
				receiveProtoMsg *proto.Message
				respProtoMsg    *proto.Message
				powData         pow.Data
			)

			BeforeEach(func() {
				server = NewTCPServer(cfg, log, services, powWorkerMock, msgrFn)
				clientConnMock = new(testClientTCPConn)
			})

			It("should call msgr.Receive func", func() {
				msgrMock.ReceiveMock.
					Expect().
					Return(nil, testErr)

				server.handleConnection(ctx, clientConnMock)
			})

			When("msgr.Receive return Challenge message", func() {
				BeforeEach(func() {
					receiveProtoMsg = &proto.Message{
						Type: proto.Challenge,
						Data: nil,
					}

					msgrMock.ReceiveMock.
						Return(receiveProtoMsg, nil)
				})

				It("should call powWorker.GenerateNew func", func() {
					powWorkerMock.GenerateNewMock.
						Expect(ctx, clientConnMock.RemoteAddr().String()).
						Return(nil, testErr)

					server.handleConnection(ctx, clientConnMock)
				})

				When("powWorker.GenerateNew returns valid response with powData", func() {
					BeforeEach(func() {
						fuzzer.NumElements(2, 4).Fuzz(&powData)

						powWorkerMock.GenerateNewMock.
							Return(powData, nil)

						respProtoMsg = &proto.Message{
							Type: proto.Challenge,
							Data: powData,
						}
					})

					It("should call msgr.Send func", func() {
						msgrMock.SendMock.
							Expect(respProtoMsg).
							Return(testErr)

						server.handleConnection(ctx, clientConnMock)
					})
				})
			})

			When("msgr.Receive return Quote message", func() {
				BeforeEach(func() {
					fuzzer.NumElements(2, 4).Fuzz(&powData)

					receiveProtoMsg = &proto.Message{
						Type: proto.Quote,
						Data: powData,
					}

					msgrMock.ReceiveMock.
						Expect().
						Return(receiveProtoMsg, nil)
				})

				It("should call powWorker.ValidateWorkDone func", func() {
					powWorkerMock.ValidateWorkDoneMock.
						Expect(ctx, clientConnMock.RemoteAddr().String(), powData).
						Return(testErr)

					server.handleConnection(ctx, clientConnMock)
				})

				When("powWorker.ValidateWorkDone returns no error", func() {
					var quote quotes_book.Quote

					BeforeEach(func() {
						powWorkerMock.ValidateWorkDoneMock.
							Return(nil)
					})

					It("should call quotesBook.GetRandomQuote func", func() {
						quotesBookMock.GetRandomQuoteMock.
							Expect(ctx).
							Return(quote, testErr)

						server.handleConnection(ctx, clientConnMock)
					})

					When("powWorker.ValidateWorkDone returns no error", func() {
						fuzzer.Fuzz(&quote)

						BeforeEach(func() {
							quotesBookMock.GetRandomQuoteMock.
								Return(quote, nil)

							respProtoMsg = &proto.Message{
								Type: proto.Quote,
								Data: quote.ToJson(),
							}
						})

						It("should call msgr.Send func", func() {
							msgrMock.SendMock.
								Expect(respProtoMsg).
								Return(testErr)

							server.handleConnection(ctx, clientConnMock)
						})
					})
				})
			})

			When("msgr.Receive return Stop message", func() {
				BeforeEach(func() {
					receiveProtoMsg = &proto.Message{
						Type: proto.Stop,
						Data: nil,
					}

					msgrMock.ReceiveMock.
						Expect().
						Return(receiveProtoMsg, nil)
				})

				It("should close connection", func() {
					server.handleConnection(ctx, clientConnMock)
				})
			})

			When("msgr.Receive return unsupported message", func() {
				BeforeEach(func() {
					receiveProtoMsg = &proto.Message{
						Type: 100,
						Data: nil,
					}

					msgrMock.ReceiveMock.
						Expect().
						Return(receiveProtoMsg, nil)
				})

				It("should close connection", func() {
					server.handleConnection(ctx, clientConnMock)
				})
			})
		})
	})
})

// testClientTCPConn is used to create test net.Conn structure (no real conn, just struct)
type testClientTCPConn struct {
	net.TCPConn
}

func (c testClientTCPConn) RemoteAddr() net.Addr {
	return new(testClientTCPConnAddr)
}

func (c testClientTCPConn) Close() error {
	return nil
}

// testClientTCPConnAddr is used to mock net.TCPConn called functions
type testClientTCPConnAddr struct {
	net.TCPAddr
}

func (c testClientTCPConnAddr) String() string {
	return "FakeRemoteAddr"
}
