GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
DIST_DIR=bin
BINARY_NAME=bookmarks
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_WIN=$(BINARY_NAME)_win.exe
GOARCH="amd64"

all: clean fmt test build

clean:
	rm -rf ${DIST_DIR}

test:
	$(GOTEST) -v ./...

build: ## show this help
	@echo 'Building MacOS binary'
	GOARCH=${GOARCH} GOOS=darwin go build -o ${DIST_DIR}/${BINARY_NAME}-darwin-${GOARCH}
	@echo 'Building Linux binary'
	GOARCH=${GOARCH} GOOS=linux go build -o ${DIST_DIR}/${BINARY_NAME}-linux-${GOARCH}
	@echo 'Building Windows binary'
	GOARCH=${GOARCH} GOOS=windows go build -o ${DIST_DIR}/${BINARY_NAME}-windows-${GOARCH}.exe

fmt:    ## format the go source files
	go fmt ./...

lint: # https://staticcheck.io/
	staticcheck ./...

dcup:
	docker-compose up --build

docker-build:
	docker build -t sivaprasadreddy/bookmarks .
