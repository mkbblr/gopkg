# gopkg   
<sub>A collection of random go packages and utilities developed as hobby project. Launch `Gitpod.io` [Workspace](https://gitpod.io#https://github.com/mkbblr/gopkg) to check it out.</sub>

---

#### Working

```shell
# run the coding challenges from project euler
go run challenges/cmd/euler/euler.go

# run the leetcode challenges
go run leet/cmd/leetit.go

# install all and then launch 
go install ./... && euler && leetcode

# build check
go build ./...
```

or simplely use `make` as 

```shell

# clean, build, install, run, run to print build info
make clean && make && make install && euler -show && euler

```