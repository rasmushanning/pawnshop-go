unittest: 
	go test ./server/pkg/... -cover -count 1

endtoendtests:
	go test ./endtoend_tests/... -count 1 -v

test: unittest endtoendtests

lint:
	golangci-lint run

run-server:
	go run server/cmd/pawnshop/main.go

run-client:
	go run client/cmd/client/main.go

build-server:
	go build -o bin/server server/cmd/pawnshop/main.go

build-client:
	go build -o bin/client client/cmd/client/main.go

build: build-server build-client