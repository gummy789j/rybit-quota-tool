AUTH_TOKEN := [YOUR_AUTH_TOKEN]
VIP_LEVEL := [YOUR_VIP_LEVEL]

rqt-mac:
	@./exe/rqt-darwin-arm64 quota --auth="$(AUTH_TOKEN)" --vip=$(VIP_LEVEL)

rqt-linux:
	@./exe/rqt-linux-amd64 quota --auth="$(AUTH_TOKEN)" --vip=$(VIP_LEVEL)

rqt-windows:
	@./exe/rqt-windows-amd64.exe quota --auth="$(AUTH_TOKEN)" --vip=$(VIP_LEVEL)

# dev use only
batch-build:
	@GOOS=darwin GOARCH=arm64 go build -o exe/rqt-darwin-arm64 main.go
	@GOOS=linux GOARCH=amd64 go build -o exe/rqt-linux-amd64 main.go
	@GOOS=windows GOARCH=amd64 go build -o exe/rqt-windows-amd64.exe main.go