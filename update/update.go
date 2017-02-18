package update

import (
	gu "github.com/inconshreveable/go-update"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"os/exec"
	"bytes"
	"fmt"
	"os"
	"io"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"time"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"google.golang.org/grpc/codes"
)

const RecordFile = "update-record.json"

type Runner struct {
	UpdatePeriod time.Duration
	Remoter
	Applier
	Rollbacker
	config.Persister
}

func NewRunner(p config.Persister) *Runner {
	return &Runner{Remoter:&S3Remoter{}, Applier:gu.Apply, Rollbacker:gu.RollbackError, Persister:p, UpdatePeriod:time.Hour}
}

type Applier func(update io.Reader, opts gu.Options) error

type Rollbacker func (err error) error

type Record struct {
	LastUpdate int64
}


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
	due, err := r.isUpdateDue()
	if !due {
		fmt.Println("Update not due.")
		return nil
	}
	if err != nil {
		return err
	}

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
		r.recordUpdate()
		return err
	} else {
		r.CleanUp()
		r.recordUpdate()
		fmt.Printf("Updated kwk to %s\n", ri.Current)
		return nil
	}
}

func (r *Runner) recordUpdate() error {
	ur := &Record{LastUpdate:time.Now().Unix()}
	return r.Persister.Upsert(RecordFile, ur)
}

func (r *Runner) isUpdateDue() (bool, error) {
	ur := &Record{}
	hiatus := time.Now().Unix() - int64(r.UpdatePeriod/time.Second)
	if err := r.Persister.Get(RecordFile, ur, hiatus); err != nil {
		err2, ok := err.(*models.ClientErr)
		if !ok {
			return false, err
		}
		if err2.TransportCode == codes.NotFound {
			// If no record is found then lets update.
			return true, nil
		}
		return false, err2
	}
	return false, nil
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
