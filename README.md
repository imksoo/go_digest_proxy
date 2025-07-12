# Go Digest Reverse Proxy

This is a reverse proxy server written in Go that authenticates to a backend server using HTTP Digest Authentication and relays the response to the client without requiring authentication.

## Features
- Transparent reverse proxy
- Supports HTTP Digest Authentication
- Configurable via environment variables or command-line arguments

## Usage

### 1. Set Environment Variables
Create a `.env` file or export variables directly:
```bash
export DIGEST_USER="your_username"
export DIGEST_PASS="your_password"
export BACKEND_URL="http://example.com"
export PORT=8080
```

### 2. Run
```bash
go run main.go
```

or with flags:
```bash
go run main.go --user your_username --pass your_password --url http://example.com --port 8080
```

### 3. Docker
```bash
docker build -t go-digest-proxy .
docker run -p 8080:8080 --env-file .env go-digest-proxy
```

## Test
Run unit tests:
```bash
go test -v
```

## License
MIT