env GOOS=linux GOARCH=amd64 go build -o httptester-linux-amd64
env GOOS=darwin GOARCH=amd64 go build -o httptester-darwin-amd64
env GOOS=darwin GOARCH=arm64 go build -o httptester-darwin-arm64
env GOOS=windows GOARCH=amd64 go build -o httptester-windows-amd64