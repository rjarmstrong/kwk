// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"bytes"
	"flag"
)

var updateMode = flag.Bool("update", false, "Update.")

func init() {
	flag.Parse()
}

func main() {
	if *updateMode {
		fmt.Println("<silent>Updating app.<silent>")
	} else {
		SilentUpdateCheck()
		fmt.Println("Running app as normal")
	}
}

func SilentUpdateCheck() {
	fmt.Println("Checking for updates")
	c := exec.Command("./install", "--update")
	c.Stdin = os.Stdin
	out, err := c.StdoutPipe()
	if err != nil {
		// log to file
	}
	var stderr bytes.Buffer
	c.Stdout = os.Stdout
	c.Stderr = &stderr
	err = c.Run()
	if err != nil {
		// log to file
	}
	out.Close()
}
