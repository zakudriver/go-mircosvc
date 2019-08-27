# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=exec
BINARY_UNIX=$(BINARY_NAME)_unix

all: run

build:
	$(GOBUILD) -o ./gateway/build/$(BINARY_NAME) -v ./gateway/cmd/*.go
	$(GOBUILD) -o ./servers/user/build/$(BINARY_NAME) -v ./servers/user/cmd/*.go

test:
	$(GOTEST) -v ./

clean:
	$(GOCLEAN) -i -n
	rm -f ./gateway/build/$(BINARY_NAME)
	rm -f ./gateway/build/$(BINARY_UNIX)
	rm -f ./servers/user/build/$(BINARY_NAME)
	rm -f ./servers/user/build/$(BINARY_UNIX)

run:
	$(GOBUILD) -o ./gateway/build/$(BINARY_NAME) -v ./gateway/cmd/gateway
	./gateway/build/$(BINARY_NAME) &
	$(GOBUILD) -o ./servers/order/build/$(BINARY_NAME) -v ./servers/order/cmd/order
	./servers/order/build/$(BINARY_NAME) &

restart:
	kill -INT $$(cat pid)
	$(GOBUILD) -o ./build/$(BINARY_NAME) -v ./
	./build/$(BINARY_NAME)

deps:
	$(GOGET) github.com/kardianos/govendor
	cd ./gateway && govendor sync
	cd ./servers/user && govendor sync

cross:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./gateway/build/$(BINARY_NAME) -v ./gateway/cmd/gateway
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./servers/user/build/$(BINARY_NAME) -v ./servers/user/cmd
