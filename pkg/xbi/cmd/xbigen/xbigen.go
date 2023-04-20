/*
xbigen generates code to inject additional build information. It creates a file
named main_generated.go. Ideally it should be invoked using go:generate
directive at the top of main package code. like -

	package main
	import "fmt"
	//go:generate sh -c "$(go env GOPATH)/bin/xbigen"
	func main() {
		fmt.Println("hello world!!!")
	}

Then run `go generate`, then `go build`. THAT IS ALL.

Now we can invoke the generated binary with the following
command line flags to display the build information in different format.

	./main [-onelinebi | -fullbi | -jsonbi | -show]

Currently it injects the following information into the binary.
 1. Git branch name
 2. Git log (latest 5 commits),
 3. Git status (modified, untracked, deleted)
 4. Git local commits (latest 5 commits that are not pushed to default remote i.e. origin)
 5. Build time
 6. Build host
 7. Build user name
*/
package main

import (
	"bufio"
	_ "embed"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/mkbblr/gopkg/pkg/xbi"
	"github.com/pkg/errors"
)

func main() {

	ext := make(map[string]string)

	status, err := runCommand("git", "status", "--porcelain=v1", "-b", "-uall")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ext[xbi.X_BI_KEY_GIT_STATUS] = base64.StdEncoding.EncodeToString(status)

	remote, err := runCommand("git", "config", "--get", "remote.origin.url")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ext[xbi.X_BI_KEY_GIT_ORIGIN] = string(remote)

	gitLog, err := runCommand("git", "log", "-n", "5", "--pretty=format:%h: %s")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if gitBranch, err := runCommand("git branch | grep \"*\" | cut -c 3- | grep -v detached"); err == nil {
		if gitLocalCommit, err := runCommand("git log origin/" + string(gitBranch) + ".." + string(gitBranch) + " -n 5  --pretty=format:%h: %s"); err == nil {
			ext[xbi.X_BI_KEY_GIT_LOCAL_COMMITS] = base64.StdEncoding.EncodeToString(gitLocalCommit)
		}
	}

	ext[xbi.X_BI_KEY_GIT_LOG] = base64.StdEncoding.EncodeToString(gitLog)

	ext[xbi.X_BI_KEY_BUILD_TIME] = time.Now().Format("2006-01-02T15:04:05")

	wd, err := os.Getwd()
	if err == nil {
		ext[xbi.X_BI_KEY_BUILD_PATH] = wd
	}

	host, err := os.Hostname()
	if err == nil {
		ext[xbi.X_BI_KEY_BUILD_HOST] = host
	}

	user, err := user.Current()
	if err == nil {
		ext[xbi.X_BI_KEY_BUILD_USER] = user.Name
	}

	generateCode(ext)
}

func runCommand(name string, arg ...string) (stdout []byte, err error) {
	c := exec.Command(name, arg...)
	o, _ := c.StdoutPipe()
	// e, _ := c.StderrPipe()

	sch := make(chan []byte)
	ech := make(chan error)
	reader := bufio.NewReader(o)

	go func() {
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			ech <- err
		} else {
			sch <- data
		}
		close(sch)
		close(ech)
	}()

	if err := c.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start command")
	}

	select {
	case data := <-sch:
		return data, nil
	case err := <-ech:
		return nil, errors.Wrap(err, "failed to run command")
	case <-time.After(5 * time.Second):
		return nil, errors.New("timeout running command")
	}
}

//go:embed main.go.template
var fileBytes []byte

func generateCode(ext map[string]string) {

	fileContent := string(fileBytes)

	xbiKV := ""
	for k, v := range ext {
		xbiKV = xbiKV + k + ":" + v + ","
	}
	xbiKV = base64.StdEncoding.EncodeToString([]byte(xbiKV))

	k := xbi.X_BI_KEY_KV_PAIR
	fileContent = strings.Replace(fileContent, "%"+k+"%", k+":"+xbiKV, -1)

	err := ioutil.WriteFile("main_generated.go", []byte(fileContent), 0644)
	if err != nil {
		fmt.Println("failed to generate file, ", err)
		os.Exit(1)
	}
}
