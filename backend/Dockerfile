FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o ./main.go ./cmd/server/main.go

FROM debian:latest
WORKDIR /root/

RUN apt-get update && apt-get install -y ca-certificates --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/ .

RUN chmod +x ./main.go

EXPOSE 8080

CMD ["./main.go"]