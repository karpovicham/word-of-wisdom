FROM golang:1.19.6-alpine AS builder

WORKDIR /app

# Copies everything from your root directory into /app
COPY . .

RUN go mod download
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

# multistage build to copy only binary and assets
FROM scratch

COPY --from=builder /app/main /
COPY --from=builder /app/assets/quotes.json /assets/quotes.json

ENTRYPOINT ["/main"]
