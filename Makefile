all: install

$(go env GOPATH)/bin/xbigen:
	echo "installing xbigen ..."
	go install github.com/mkbblr/gopkg/pkg/xbi/cmd/xbigen


install: $(go env GOPATH)/bin/xbigen
	echo "runnnig generate...."
	go generate ./...
	CGO_ENABLED=0 go install ./...


clean:
	-go clean ./...
	-find . -type f -name main_generated.go | xargs rm
	-rm $$(go env GOPATH)/bin/*

default: all
.PHONY: all	install clean