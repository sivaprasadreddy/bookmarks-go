# Bookmarks Go

## Live Reloading using [Air](https://github.com/cosmtrek/air)

```shell
$ go install github.com/cosmtrek/air@latest
```

## Run application

```shell
$ docker-compose up -d bookmarks-db
$ air
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
$ go build -o ./dist/bookmarks cmd/bookmarks/main.go
$ ./dist/bookmarks
$ open http://localhost:8080
```