

install:
	go mod download
	go mod tidy

build-script:
	go build -o ./netmap ./cmd/scanner/main.go

build-daemon:
	go build -o ./netmap-d ./cmd/server/main.go

image:
	docker build . -t netmap-backend