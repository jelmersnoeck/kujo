test:
	go test -v ./...

build:
	go build -o bin/kujo .

install:
	go install .

docker:
	docker build -t kujo .
