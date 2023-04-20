# mygocode
A monorepo container to hold any golang coding experiment.

# Start with gitpod.io - click this link
https://gitpod.io#https://github.com/mkbblr/gopkg


# Working 

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