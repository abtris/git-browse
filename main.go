package main

import (
	"bytes"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

// Read .git/config and get github/bitbucket URL for open
// git config --local --get-regexp remote.origin.url
func GetLink(dir string) (string, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("git", "config", "--local", "--get-regexp", "remote.origin.url")
	cmd.Stdout = &stdout
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				return "", err
			}
		}
		return "", err
	}

	line := strings.Trim(stdout.String(), "\n")
	re := regexp.MustCompile(`(remote\.origin\.url)\s((ssh|http(s?))://)?git@(?P<hostname>.+):(\d+/)?(?P<repository>.+)(\.git)`)
	groupNames := re.SubexpNames()
	var hostname, repository string
	for _, match := range re.FindAllStringSubmatch(line, -1) {
		for groupIdx, group := range match {
			name := groupNames[groupIdx]
			if name == "hostname" {
				hostname = group
			}
			if name == "repository" {
				repository = group
			}
		}
	}
	var output string
	if strings.Contains(hostname, "bitbucket") {
		bitbucket := strings.Split(repository, "/")
		output = fmt.Sprintf("https://%s/projects/%s/repos/%s/browse", hostname, bitbucket[0], bitbucket[1])
	} else {
		output = fmt.Sprintf("https://%s/%s", hostname, repository)
	}
	return output, nil
}

// OpenLink in default browser
func OpenLink(url string) {
	open.Run(url)
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	link, err := GetLink(dir)
	if err != nil {
		log.Fatal(err)
	}
	OpenLink(link)
}
