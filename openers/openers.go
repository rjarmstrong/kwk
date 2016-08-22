package openers

import (
	"fmt"
	"os/exec"
	"bytes"
)

func Open(uri string) {
	execSafe("open", uri)
}

func OpenCovert(uri string){
	execSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", uri)
	execSafe("osascript", "-e", "activate application \"Google Chrome\"")
}

func execSafe(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
}