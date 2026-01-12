# --------- CONFIG ---------
APP_NAME := netmap
DAEMON_BIN := $(APP_NAME)-d
DAEMON_MAIN := ./cmd/daemon/main.go

SCRIPT_BIN := $(APP_NAME)
SCRIPT_MAIN := ./cmd/script/main.go

DOCKER_IMAGE := $(APP_NAME)-backend

# --------- PHONY ---------

.PHONY: install build-script build-daemon image

# --------- DEFAULT ---------

all: build-daemon build-script

# --------- STEPS ---------

deps:
	echo "Downloading dependencies..."
	go mod download

tidy:
	echo "Tidying up modules..."
	go mod tidy

build-script: deps tidy
	echo "Building script binary..."
	go build -o $(SCRIPT_BIN) $(SCRIPT_MAIN)
	chmod +x ./$(SCRIPT_BIN)
	sudo setcap cap_net_raw+ep ./$(SCRIPT_BIN)


build-daemon: deps tidy
	echo "Building daemon binary..."
	go build -o $(DAEMON_BIN) $(DAEMON_MAIN)

# --------- DOCKER ---------
image:
	echo "Building Docker image..."
	docker build . -t $(DOCKER_IMAGE)