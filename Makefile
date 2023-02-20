check:
	go vet -v ./...

gen_code:
	minimock -g -i ./pkg/pow.* -o ./pkg/pow -s _mock.go
	minimock -g -i ./internal/messenger.* -o ./internal/messenger -s _mock.go
	minimock -g -i ./service/quotes_book.* -o ./service/quotes_book -s _mock.go
	minimock -g -i ./service/client/resolver.* -o ./service/client/resolver -s _mock.go
	go mod tidy

test: check
	go clean --testcache
	go test -v -race -coverprofile=coverage.txt ./...

run_server:
	go run cmd/server/main.go

run_client:
	go run cmd/client/main.go

start:
	docker compose up --abort-on-container-exit --force-recreate --build server --build client
