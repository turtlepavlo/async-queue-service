.PHONY: build test lint integration-test

build:
	go build ./...

test:
	go test -short -race -count=1 ./...

lint:
	golangci-lint run -c .golangci.yml ./...

integration-test:
	docker compose -f test.compose.yaml up -d
	sleep 10
	go test -count=1 -timeout=120s ./...
	docker compose -f test.compose.yaml down
