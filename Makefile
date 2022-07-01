GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

build:
	@# for current arch. system
	@CGO_ENABLED=1 go build -ldflags="-s -w" -o '$(GOBIN)/8mar5d-hub' ./cmd/8mar5d-hub/main.go || exit
	@# for MIPS arch. system on Onion Omega2/Omega2+
	@GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags="-s -w" -o '$(GOBIN)/mips/8mar5d-hub' ./cmd/8mar5d-hub/main.go || exit

run:
	@go build -o '$(GOBIN)/8mar5d-hub' ./cmd/8mar5d-hub/main.go
	@'$(GOBIN)/8mar5d-hub' --config='$(GOBASE)/configs/8mar5d-hub.yml'

test:
	@go test -v -count=1 -race -timeout=60s ./...

install-deps: deps
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint && go mod vendor && go mod verify

lint: install-deps
	@golangci-lint run ./...

deps:
	@go mod tidy && go mod vendor && go mod verify

install:
	@go mod download

generate:
	@go generate ./...

.PHONY: build
