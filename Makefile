# GO parameters
GOCMD=go
GOBUILD=$(GOCMD) build
BUILD_TARGET_PLUGIN=cf_eth_wallet

fmt:
	$(GOCMD) fmt ./...

build: fmt
	$(GOBUILD) -o ./$(BUILD_TARGET_PLUGIN)

clean:
	@echo "Cleaning..."
	rm -rf ./$(BUILD_TARGET_PLUGIN)