bin=ydcv
.PHONY: clean
all: clean $(bin)

$(bin):
	for ARCH in amd64 386; do \
	    for OS in windows linux darwin; do \
		CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=$(OS) go build -o $(bin)_$(OS)_$(ARCH) -ldflags "-s -w -X main.VERSION=git-$$(git rev-parse --short HEAD)"; \
	    done; \
	done; 

clean:
	rm -fv $(bin)_*
