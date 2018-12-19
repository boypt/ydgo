bin=ydcv
.PHONY: clean
all: clean $(bin)

$(bin):
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $(bin)_linux -ldflags "-s -w -X main.VERSION=git-$$(git rev-parse --short HEAD)"; \
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o $(bin).exe -ldflags "-s -w -X main.VERSION=git-$$(git rev-parse --short HEAD)"; \
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o $(bin)_darwin -ldflags "-s -w -X main.VERSION=git-$$(git rev-parse --short HEAD)"; \

clean:
	rm -fv $(bin)_*
