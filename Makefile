VERSION := v1.0.2
BINARY_NAME := cyber
TARGETS := windows-amd64 windows-arm64 darwin-amd64 darwin-arm64 linux-amd64 linux-arm64

.PHONY: all
all: $(TARGETS)

$(TARGETS):
	@echo "Building $(BINARY_NAME)-$(VERSION)-$@"
	@GOOS=$(firstword $(subst -, ,$@)) GOARCH=$(lastword $(subst -, ,$@)) go build -o build/$(BINARY_NAME)-$(VERSION)-$@$(if $(filter windows-%,$@),.exe,) main.go
	@tar czf build/$(BINARY_NAME)-$(VERSION)-$@.tar.gz -C build $(BINARY_NAME)-$(VERSION)-$@$(if $(filter windows-%,$@),.exe,)
	@rm build/$(BINARY_NAME)-$(VERSION)-$@$(if $(filter windows-%,$@),.exe,)

.PHONY: clean
clean:
	@rm -rf build

