PROJECT_NAME := cyber
VERSION := v1.0.0

OS := windows darwin linux
ARCH := amd64 arm64

.PHONY: all
all: $(OS)

$(OS): %: %-amd64 %-arm64

$(OS:%=%-amd64): %-amd64: build
	@echo "Packaging $* amd64..."
	mkdir -p build/$*/amd64
	cp build/$(PROJECT_NAME)-$*-amd64 build/$*/amd64/
	tar czvf build/$(PROJECT_NAME)-$(VERSION)-$*-amd64.tar.gz -C build/$*/amd64 $(PROJECT_NAME)-$*-amd64

$(OS:%=%-arm64): %-arm64: build
	@echo "Packaging $* arm64..."
	mkdir -p build/$*/arm64
	cp build/$(PROJECT_NAME)-$*-arm64 build/$*/arm64/
	tar czvf build/$(PROJECT_NAME)-$(VERSION)-$*-arm64.tar.gz -C build/$*/arm64 $(PROJECT_NAME)-$*-arm64

.PHONY: build
build:
	@echo "Building binaries..."
	$(foreach os, $(OS), $(foreach arch, $(ARCH), \
		GOOS=$(os) GOARCH=$(arch) go build -o build/$(PROJECT_NAME)-$(os)-$(arch) -ldflags="-s -w" ; \
	))

.PHONY: clean
clean:
	rm -rf build
