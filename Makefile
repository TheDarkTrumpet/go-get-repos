BINARY_NAME=go-get-repos

all: clean .linux .windows

.linux:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o build/${BINARY_NAME} main.go

.windows:
	mkdir -p build
	GOOS=windows GOARCH=amd64 go build -o build/${BINARY_NAME}.exe main.go

clean:
	go clean
	rm -Rf build


