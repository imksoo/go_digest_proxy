FROM golang:1.22-alpine
WORKDIR /app
COPY . .
RUN go build -o go-digest-proxy main.go
CMD ["./go-digest-proxy"]