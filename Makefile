# RaspController Build and Utility Commands

# Default target: displays usage instructions
ALL:
	@echo "RaspController Compiler"
	@echo "Usage:"
	@echo "  make build-arm64  # For Raspberry Pi 3 and newer (64-bit)"
	@echo "  make build-arm32  # For Raspberry Pi 3 and older (32-bit)"
	@echo "  make compress     # Compress the binary using UPX"
	@echo "  make doc          # Run Go documentation server on port 6060"

# Build for ARM64 (Raspberry Pi 3 and newer)
build-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o raspc cmd/raspc/main.go
	$(MAKE) show


# Build for ARM32 (Raspberry Pi 3 and older)
build-arm32:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -ldflags="-s -w" -o raspc cmd/raspc/main.go
	$(MAKE) show


# Compress the binary using UPX
compress:
	upx --best raspc

show:
	@file raspc
	@sha256sum raspc

# Run Go documentation server on port 6060
doc:
	godoc -http=:6060

.PHONY: ALL build-arm64 build-arm32 compress doc
