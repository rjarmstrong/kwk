package update

import (
	gu "github.com/inconshreveable/go-update"
	"os/exec"
	"bytes"
	"time"
	"os"
	"io"
	"bitbucket.com/sharingmachine/types/errs"
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"bitbucket.com/sharingmachine/kwkcli/src/persist"
)

const RecordFile = "update-record.json"

// SilentCheckAndRun spawns a new process to check for updates and run.
func SilentCheckAndRun() {
	cmd, err := os.Executable()
	models.Debug("Initiating silent update check for: %s", cmd)
	if err != nil {
		models.Debug("If you are running nacl or OpenBSD they are not supported.")
		models.LogErr(err)
	}
	exe(false, cmd,"update", "silent")
}

type Runner struct {
	UpdatePeriod   time.Duration
	Remoter
	Applier
	Rollbacker
	persist.Persister
	currentVersion string
}

func NewRunner(p persist.Persister, version string) *Runner {
	return &Runner{currentVersion: version, Remoter:&S3Remoter{}, Applier:gu.Apply, Rollbacker:gu.RollbackError, Persister:p, UpdatePeriod:time.Hour}
}

type Applier func(update io.Reader, opts gu.Options) error

type Rollbacker func (err error) error

type Record struct {
	LastUpdate int64
}

func (r *Runner) Run() error {
	due, err := r.isUpdateDue()
	if !due {
		models.Debug("Update not due.")
		return nil
	}
	if err != nil {
		models.Debug("%+v", err)
		return err
	}

	li, err := r.LatestInfo()
	if err != nil {
		models.LogErrM("Couldn't get remote release info.", err)
		return err
	}
	if li.Version == r.currentVersion {
		models.Debug("Local is same as latest version: %s", li.Version)
		r.recordUpdate()
		return nil
	}
	latest, err := r.LatestBinary()
	if err != nil {
		models.LogErrM("Couldn't get latest from remote.", err)
		return err
	}
	defer latest.Close() //TODO: Currently NOOP, should be real closer
	models.Debug("Applying update.")
	err = r.Applier(latest, gu.Options{})
	if err != nil {
		models.LogErrM("Couldn't apply update.", err)
		err = r.Rollbacker(err)
		r.CleanUp()
		r.recordUpdate()
		return err
	} else {
		r.CleanUp()
		r.recordUpdate()
		models.Debug("Updated to version: %s+%s", li.Version, li.Build)
		return nil
	}
}

func (r *Runner) recordUpdate() error {
	ur := &Record{LastUpdate:time.Now().Unix()}
	models.Debug("Updating update record.")
	return r.Persister.Upsert(RecordFile, ur)
}

func (r *Runner) isUpdateDue() (bool, error) {
	if !models.Prefs().RegulateUpdates {
		models.Debug("Updates not regulated.")
		return true, nil
	}
	ur := &Record{}
	hiatus := time.Now().Unix() - int64(r.UpdatePeriod/time.Second)
	models.Debug("Checking update is newer than: %d (Unix time seconds)", hiatus)
	if err := r.Persister.Get(RecordFile, ur, hiatus); err != nil {
		models.LogErrM("Couldn't get local update record.", err)
		err2, ok := err.(*errs.Error)
		if !ok {
			return false, err
		}
		if err2.Code == errs.CodeNotFound {
			// If no record is found then lets update.
			return true, nil
		}
		return false, err2
	}
	return false, nil
}

func exe(wait bool, name string, arg ...string) {
	c := exec.Command(name, arg...)
	c.Stdin = os.Stdin
	out, err := c.StdoutPipe()
	if err != nil {
		models.LogErrM("If you are running nacl or OpenBSD they are not supported.", err)
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
		models.LogErrM("Couldn't execute command.", err)
	}
	out.Close()
}
