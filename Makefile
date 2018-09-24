# Go parameters

GOPATH ?= $(HOME)/go

#This is how we want to name the binary output
BINARY=go-api-storage-mongodb

all: test build

build:
	cd $(GOPATH)/src; go install github.com/puneethreddy20/go-api-storage-mongodb

test:
	go test -v ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

deps:
	go get -t ./...