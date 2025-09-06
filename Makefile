build:
	go build -o go-digest-proxy main.go

run:
	go run main.go

test:
	go test -v ./...

clean:
	rm go-digest-proxy

docker-build:
	docker build -t go-digest-proxy .
