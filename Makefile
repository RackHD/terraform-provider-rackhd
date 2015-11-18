

default: test

deps:
	go get -t -v ./...

test: deps
	go test -v ./...

build: deps
	go build -o bin/terraform-provider-rackhd

install: deps
	cp -f ./bin/terraform-provider-rackhd $(shell dirname `which terraform`)
