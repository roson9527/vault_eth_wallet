# GO parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod
BUILD_TARGET=cf_eth_wallet
BUILD_LINUX_TARGET=$(BUILD_TARGET)_linux

fmt:
	$(GOCMD) fmt ./...

build: fmt
	$(GOBUILD) -o ./$(BUILD_TARGET)

clean:
	@echo "Cleaning..."
	rm -rf ./$(BUILD_TARGET)
	rm -rf ./$(BUILD_LINUX)

sum:
	shasum -a 256 ./$(BUILD_TARGET) | awk '{print $$1}'

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./$(BUILD_LINUX_TARGET) -v

download:
	$(GOMOD) download