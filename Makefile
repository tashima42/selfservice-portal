VERSION=dev
TARGET=dist/selfserviceportal

default: build

download:
	go mod download

build: download
	go build -o $(TARGET) -ldflags '-w -X main.Version=$(VERSION)' .

watch:
	air
