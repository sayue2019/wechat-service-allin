set -ex

CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o wxDriver.exe main.go