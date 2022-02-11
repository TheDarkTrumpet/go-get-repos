BINARY_NAME=go-get-repos

build_linux:
	mkdir -p build/linux
	GOOS=linux GOARCH=amd64 go build --ldflags '-linkmode external -extldflags "-static"' -o build/${BINARY_NAME} main.go

build_windows:
	mkdir -p build/windows
	GOOS=windows GOARCH=amd64 go build --ldflags '-linkmode external -extldflags "-static"' -o build/${BINARY_NAME}.exe main.go

clean:
	go clean
	rm -Rf build

all: build_linux build_windows
