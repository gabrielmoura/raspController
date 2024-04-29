# RaspController

All: echo "RaspController Compiler\nUse\: make m8x64 or make m8x32"

m8x64: GOOS=linux GOARCH=arm64  go build
m8x32: GOOS=linux GOARCH=arm GOARM=7 go build
doc: godoc -http=:6060