# Bookmarks Go

## Tools 
* Live Reloading using [Air](https://github.com/cosmtrek/air)
* Static Analysis using [staticcheck](https://staticcheck.dev/)

```shell
$ go install github.com/cosmtrek/air@latest
$ go install honnef.co/go/tools/cmd/staticcheck@latest
```

## Run application

```shell
$ docker-compose up -d bookmarks-db
$ air
```

## Run application using docker-compose

```shell
$ docker-compose up --build -d
```

## Run tests

```shell
$ go test -v ./...

# Check coverage
$ GIN_MODE=release go test -v -coverprofile=coverage.out ./...
$ go tool cover -html=coverage.out
```

## Build the application

```shell
$ go build -o ./bin/bookmarks cmd/bookmarks/main.go
$ ./bin/bookmarks
$ open http://localhost:8080
```