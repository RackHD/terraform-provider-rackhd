

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

cross:
	env GOOS=linux GOARCH=amd64 go build -o bin/terraform-provider-rackhd
	scp bin/terraform-provider-rackhd onrack@10.240.16.168:~/terraform/terraform-provider-rackhd
