test:
	go test -v ./...

build:
	go build -o bin/kujo .

docker:
	docker build -t kujo .
