package update

import (
	gu "github.com/inconshreveable/go-update"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"os/exec"
	"bytes"
	"fmt"
	"os"
	"io"
)

type Runner struct {
	Remoter
	Applier
	Rollbacker
}

func NewRunner() *Runner {
	return &Runner{Remoter:&S3Remoter{}, Applier:gu.Apply, Rollbacker:gu.RollbackError}
}

type Applier func(update io.Reader, opts gu.Options) error

type Rollbacker func (err error) error


// updater
/*
kwk [anything]

  update.ForkRun()

	1. Check for updates

	Get current app version
	TODO: Read file last check from file
		Check when last update happened, it it was older than 24 hours then...
		else 	Write file

	2. Download update

	Download from 'latest' tag for 'os'-'arch'

	3. Apply update

	Patch or switcher-oo

	4. Clean up

	5. Write file

 */
func (r *Runner) Run() error {
	ri, err := r.ReleaseInfo()
	if err != nil {
		return err
	}
	if ri.Current == sys.Version {
		return nil
	}
	latest, err := r.Latest()
	if err != nil {
		return err
	}
	defer latest.Close() //TODO: Currently NOOP, should be real closer
	err = r.Applier(latest, gu.Options{})
	if err != nil {
		fmt.Println("Update error")
		err = r.Rollbacker(err)
		r.CleanUp()
		return err
	} else {
		r.CleanUp()
		fmt.Printf("Updated kwk to %s\n", ri.Current)
		return nil
	}
}


type ReleaseInfo struct {
	Current string `json:"current"`
	Build string `json:"build"`
}

func SilentCheckAndRun() {
	fmt.Println("Checking for updates")
	//var cmd string
	//if sys.KWK_TEST_MODE {
	//	cmd = "./kwkcli"
	//} else {
	//	cmd = "kwk"
	//}
	cmd, err := os.Executable()
	fmt.Println("exe:", cmd)
	if err != nil {
		fmt.Println("If you are running nacl or OpenBSD they are not supported.")
		fmt.Println(err)
	}
	exe(false, cmd,"update", "silent")
}

func exe(wait bool, name string, arg ...string) {
	c := exec.Command(name, arg...)
	c.Stdin = os.Stdin
	out, err := c.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		// log to file
	}
	var stderr bytes.Buffer
	c.Stdout = os.Stdout
	c.Stderr = &stderr
	if wait {
		err = c.Run()
	} else {
		err = c.Start()
	}

	if err != nil {
		fmt.Println(err)
		// log to file
	}
	out.Close()
}
