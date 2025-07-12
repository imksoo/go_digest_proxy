build:
	go build -o proxy main.go

run:
	go run main.go

test:
	go test -v ./...

clean:
	rm proxy

docker-build:
	docker build -t go-digest-proxy .
