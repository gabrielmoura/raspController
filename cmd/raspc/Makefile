# RaspController

ALL:
	@echo "RaspController Compiler"
	@echo "Use: make m8x64 for >=Raspberry Pi 3+"
	@echo "Use: make m8x32 for =<Raspberry Pi 3+"

m8x64:
	GOOS=linux GOARCH=arm64 go build -o raspc_aarch64

m8x32:
	GOOS=linux GOARCH=arm GOARM=7 go build -o raspc_arm

doc:
	godoc -http=:6060
