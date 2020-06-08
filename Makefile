PKGS := $(shell go list ./...)

SHA1 := $(shell git rev-parse HEAD)
NOW = $(shell date +'%Y-%m-%d_%T')

BINARY := trpc
VERSION ?= dev
PLATFORMS := darwin linux freebsd
os = $(word 1, $@)

VER_PKG := github.com/shric/trpc/cmd/trpc

$(BINARY):
	go build -ldflags "-X $(VER_PKG).sha1ver=$(SHA1) -X $(VER_PKG).version=$(VERSION) -X $(VER_PKG).buildTime=$(NOW)" 

test:
	go test $(PKGS)

lint:
	golint ./...

$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -ldflags "-X $(VER_PKG).sha1ver=$(SHA1) -X $(VER_PKG).version=$(VERSION) -X $(VER_PKG).buildTime=$(NOW)" -o release/$(BINARY)-$(VERSION)-$(os)-amd64

release: $(PLATFORMS)

clean:
	rm -rf release
	rm -f $(BINARY)
