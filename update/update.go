package update

import (
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/models"
	gu "github.com/inconshreveable/go-update"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"bitbucket.com/sharingmachine/kwkcli/log"
	"os/exec"
	"bytes"
	"time"
	"os"
	"io"
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

func (r *Runner) Run() error {
	due, err := r.isUpdateDue()
	if !due {
		log.Debug("Update not due.")
		return nil
	}
	if err != nil {
		log.Debug("%+v", err)
		return err
	}

	ri, err := r.ReleaseInfo()
	if err != nil {
		log.Error("Couldn't get remote release info.", err)
		return err
	}
	if ri.Current == sys.Version {
		log.Debug("Local is same as latest version: %s", ri.Current)
		return nil
	}
	latest, err := r.Latest()
	if err != nil {
		log.Error("Couldn't get latest from remote.", err)
		return err
	}
	defer latest.Close() //TODO: Currently NOOP, should be real closer
	log.Debug("Applying update.")
	err = r.Applier(latest, gu.Options{})
	if err != nil {
		log.Error("Couldn't apply update.", err)
		err = r.Rollbacker(err)
		r.CleanUp()
		r.recordUpdate()
		return err
	} else {
		r.CleanUp()
		r.recordUpdate()
		log.Debug("Updated to version: %s+%s", ri.Current, ri.Build)
		return nil
	}
}

func (r *Runner) recordUpdate() error {
	ur := &Record{LastUpdate:time.Now().Unix()}
	log.Debug("Updating update record.")
	return r.Persister.Upsert(RecordFile, ur)
}

func (r *Runner) isUpdateDue() (bool, error) {
	ur := &Record{}
	hiatus := time.Now().Unix() - int64(r.UpdatePeriod/time.Second)
	log.Debug("Checking update is newer than: %d (Unix time seconds)", hiatus)
	if err := r.Persister.Get(RecordFile, ur, hiatus); err != nil {
		log.Error("Couldn't get local update record.", err)
		err2, ok := err.(*models.ClientErr)
		if !ok {
			return false, err
		}
		if err2.Contains(models.Code_NotFound) {
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
	cmd, err := os.Executable()
	log.Debug("Initiating silent update check for: %s", cmd)
	if err != nil {
		log.Error("If you are running nacl or OpenBSD they are not supported.", err)
	}
	exe(false, cmd,"update", "silent")
}

func exe(wait bool, name string, arg ...string) {
	c := exec.Command(name, arg...)
	c.Stdin = os.Stdin
	out, err := c.StdoutPipe()
	if err != nil {
		log.Error("If you are running nacl or OpenBSD they are not supported.", err)
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
		log.Error("Couldn't execute command.", err)
	}
	out.Close()
}
