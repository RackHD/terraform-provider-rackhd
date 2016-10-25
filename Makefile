

default: test

deps:
	go get -t -v ./...

test: deps
	go test -v ./...

clean:
	rm -f ./bin/terraform-provider-rackhd

build: deps
	go build -o bin/terraform-provider-rackhd

install: clean build
	cp -f ./bin/terraform-provider-rackhd $(shell dirname `which terraform`)
