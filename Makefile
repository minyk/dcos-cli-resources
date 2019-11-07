PWD := `pwd`

default: build

build: cleanall linux darwin windows

cleanall: clean-linux clean-darwin clean-windows

clean-linux:
	rm -f ./build/dcos-resources-linux

clean-darwin:
	rm -f ./build/dcos-resources-darwin

clean-windows:
	rm -f ./build/dcos-resources-windows.exe

linux: clean-linux
	docker run --rm -e "GOOS=linux" -e "GOARCH=amd64" -v $(PWD):/usr/src/github.com/minyk/dcos-resources -w /usr/src/github.com/minyk/dcos-resources golang:1.13 go build -ldflags="-s -w ${GO_LDFLAGS}" -v -o build/dcos-resources-linux

darwin: clean-darwin
	docker run --rm -e "GOOS=darwin" -e "GOARCH=amd64" -v $(PWD):/usr/src/github.com/minyk/dcos-resources -w /usr/src/github.com/minyk/dcos-resources golang:1.13 go build -ldflags="-s -w ${GO_LDFLAGS}" -v -o build/dcos-resources-darwin

windows: clean-windows
	docker run --rm -e "GOOS=linux" -e "GOARCH=amd64" -v $(PWD):/usr/src/github.com/minyk/dcos-resources -w /usr/src/github.com/minyk/dcos-resources golang:1.13 go build -ldflags="-s -w ${GO_LDFLAGS}" -v -o build/dcos-resources-windows.exe
