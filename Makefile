GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
MAIN_SRC_DIR=./cmd/bookmarks
DIST_DIR=bin
BINARY_NAME=bookmarks
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_WIN=$(BINARY_NAME)_win.exe
GOARCH="amd64"

.DEFAULT_GOAL := test

all: clean fmt test build

clean:
	rm -rf ${DIST_DIR}

fmt:
	goimports -w .

lint: fmt
	staticcheck ./...

test: lint
	$(GOTEST) -v ./...

build: test
	@echo 'Building MacOS binary'
	GOARCH=${GOARCH} GOOS=darwin go build -o ${DIST_DIR}/${BINARY_NAME}-darwin-${GOARCH} ${MAIN_SRC_DIR}
	@echo 'Building Linux binary'
	GOARCH=${GOARCH} GOOS=linux go build -o ${DIST_DIR}/${BINARY_NAME}-linux-${GOARCH} ${MAIN_SRC_DIR}
	@echo 'Building Windows binary'
	GOARCH=${GOARCH} GOOS=windows go build -o ${DIST_DIR}/${BINARY_NAME}-windows-${GOARCH}.exe ${MAIN_SRC_DIR}

dcup:
	docker-compose up --build

docker-build:
	docker build -t sivaprasadreddy/bookmarks .
