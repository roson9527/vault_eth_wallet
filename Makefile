# GO parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod
BUILD_TARGET=cf_eth_wallet
BUILD_LINUX_TARGET=$(BUILD_TARGET_PLUGIN)_linux

fmt:
	$(GOCMD) fmt ./...

build: fmt
	$(GOBUILD) -o ./$(BUILD_TARGET)

clean:
	@echo "Cleaning..."
	rm -rf ./$(BUILD_TARGET)
	rm -rf ./$(BUILD_LINUX)

sum:
	shasum -a 256 ./$(BUILD_TARGET_PLUGIN)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_LINUX) -v

download:
	$(GOMOD) download