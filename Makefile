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


$(git rev-parse --git-dir)/hooks/pre-commit:
	cp tools/git-hooks/pre-commit $(git rev-parse --git-dir)/hooks/

setup: $(git rev-parse --git-dir)/hooks/pre-commit

# .oneshell:
# check:
# 	@echo off
# 	@diff $(git rev-parse --git-dir)/hooks/pre-commit git-hooks/pre-commit > /dev/null 2>&1
# 	@test "$?" -ne "0" && echo  "hooks not installed" || echo "all ok !!!"


default: all
.PHONY: all	install clean setup check