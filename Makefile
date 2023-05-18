rqt-mac:
	@./exe/rqt-darwin-arm64 quota

rqt-linux:
	@./exe/rqt-linux-amd64 quota

rqt-windows:
	@./exe/rqt-windows-amd64.exe quota

# dev use only
batch-build:
	@GOOS=darwin GOARCH=arm64 go build -o exe/rqt-darwin-arm64 main.go
	@GOOS=linux GOARCH=amd64 go build -o exe/rqt-linux-amd64 main.go
	@GOOS=windows GOARCH=amd64 go build -o exe/rqt-windows-amd64.exe main.go