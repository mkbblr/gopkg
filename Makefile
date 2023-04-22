all: check install

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

PRE_COMMIT_HOOK := $(shell git rev-parse --git-dir)/hooks/pre-commit
$(PRE_COMMIT_HOOK):
	cp tools/git-hooks/pre-commit $(PRE_COMMIT_HOOK)

setup: $(PRE_COMMIT_HOOK)

.ONESHELL:
check:
	-@diff $(PRE_COMMIT_HOOK) tools/git-hooks/pre-commit > /dev/null 2>&1
	-@test $$? -eq 0 && echo  "hooks already latest" || cp tools/git-hooks/pre-commit $(PRE_COMMIT_HOOK) 


default: all
.PHONY: all	install clean setup check