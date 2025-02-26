.PHONY: staticcheck dependency clean build release all

PKGS       := $(shell go list ./...)
REPO       := github.com/guessi/eks-node-diagnostic
BUILDTIME  := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GITVERSION := $(shell git describe --tags --abbrev=8)
GOVERSION  := $(shell go version | cut -d' ' -f3)
LDFLAGS    := -s -w -X "$(REPO)/internal/variables.GitVersion=$(GITVERSION)" -X "$(REPO)/internal/variables.GoVersion=$(GOVERSION)" -X "$(REPO)/internal/variables.BuildTime=$(BUILDTIME)"

default: build

staticcheck:
	@echo "Setup staticcheck..."
	@go install honnef.co/go/tools/cmd/staticcheck@2025.1 # https://github.com/dominikh/go-tools/releases/tag/2025.1
	@echo "Check staticcheck version..."
	staticcheck --version
	@echo "Run staticcheck..."
	@for i in $(PKGS); do echo $${i}; staticcheck $${i}; done

test:
	go version
	go fmt ./...
	go vet ./...
	# go test -v ./...

dependency:
	go mod download

build-linux-x86_64:
	@echo "Creating Build for Linux (x86_64)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./releases/$(GITVERSION)/Linux-x86_64/eks-node-diagnostic
	@cp ./LICENSE ./releases/$(GITVERSION)/Linux-x86_64/LICENSE
	@tar zcf ./releases/$(GITVERSION)/eks-node-diagnostic-Linux-x86_64.tar.gz -C releases/$(GITVERSION)/Linux-x86_64 eks-node-diagnostic LICENSE

build-linux-arm64:
	@echo "Creating Build for Linux (arm64)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o ./releases/$(GITVERSION)/Linux-arm64/eks-node-diagnostic
	@cp ./LICENSE ./releases/$(GITVERSION)/Linux-arm64/LICENSE
	@tar zcf ./releases/$(GITVERSION)/eks-node-diagnostic-Linux-arm64.tar.gz -C releases/$(GITVERSION)/Linux-arm64 eks-node-diagnostic LICENSE

build-darwin-x86_64:
	@echo "Creating Build for macOS (x86_64)..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./releases/$(GITVERSION)/Darwin-x86_64/eks-node-diagnostic
	@cp ./LICENSE ./releases/$(GITVERSION)/Darwin-x86_64/LICENSE
	@tar zcf ./releases/$(GITVERSION)/eks-node-diagnostic-Darwin-x86_64.tar.gz -C releases/$(GITVERSION)/Darwin-x86_64 eks-node-diagnostic LICENSE

build-darwin-arm64:
	@echo "Creating Build for macOS (arm64)..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o ./releases/$(GITVERSION)/Darwin-arm64/eks-node-diagnostic
	@cp ./LICENSE ./releases/$(GITVERSION)/Darwin-arm64/LICENSE
	@tar zcf ./releases/$(GITVERSION)/eks-node-diagnostic-Darwin-arm64.tar.gz -C releases/$(GITVERSION)/Darwin-arm64 eks-node-diagnostic LICENSE

build-windows-x86_64:
	@echo "Creating Build for Windows (x86_64)..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./releases/$(GITVERSION)/Windows-x86_64/eks-node-diagnostic.exe
	@cp ./LICENSE ./releases/$(GITVERSION)/Windows-x86_64/LICENSE.txt
	@tar zcf ./releases/$(GITVERSION)/eks-node-diagnostic-Windows-x86_64.tar.gz -C releases/$(GITVERSION)/Windows-x86_64 eks-node-diagnostic.exe LICENSE.txt

build-linux: build-linux-x86_64 build-linux-arm64
build-darwin: build-darwin-x86_64 build-darwin-arm64
build-windows: build-windows-x86_64

build: build-linux build-darwin build-windows

clean:
	@echo "Cleanup Releases..."
	rm -rvf ./releases/*

release:
	@echo "Creating Releases..."
	@curl -LO https://github.com/tcnksm/ghr/releases/download/v0.17.0/ghr_v0.17.0_linux_amd64.tar.gz
	@tar --strip-components=1 -xvf ghr_v0.17.0_linux_amd64.tar.gz ghr_v0.17.0_linux_amd64/ghr
	./ghr -version
	./ghr -replace -recreate -token ${GITHUB_TOKEN} $(GITVERSION) releases/$(GITVERSION)/
	sha1sum releases/$(GITVERSION)/*.tar.gz > releases/$(GITVERSION)/SHA1SUM

all: staticcheck dependency clean build
