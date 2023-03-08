export CC=arm-linux-gnueabihf-gcc
export GOOS=windows
export GOARCH=amd64

go build -o server_http.exe server.go