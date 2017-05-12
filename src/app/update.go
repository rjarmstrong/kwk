package app

import (
	"bytes"
	gu "github.com/inconshreveable/go-update"
	"github.com/kwk-super-snippets/types/errs"
	"io"
	"os"
	"os/exec"
	"time"
)

const RecordFile = "update-record.json"

// SilentCheckAndRun spawns a new process to check for updates and run.
func SilentCheckAndRun() {
	cmd, err := os.Executable()
	Debug("Initiating silent update check for: %s", cmd)
	if err != nil {
		Debug("If you are running nacl or OpenBSD they are not supported.")
		LogErr(err)
	}
	exe(false, cmd, "update", "silent")
}

type Updater interface {
	Run() error
}

type UpdateRunner struct {
	UpdatePeriod time.Duration
	BinRepo
	Applier
	Rollbacker
	Persister
	currentVersion string
}

func NewUpdateRunner(p Persister, version string) Updater {
	return &UpdateRunner{currentVersion: version, BinRepo: &S3Repo{}, Applier: gu.Apply, Rollbacker: gu.RollbackError, Persister: p, UpdatePeriod: time.Hour}
}

type Applier func(update io.Reader, opts gu.Options) error

type Rollbacker func(err error) error

type Record struct {
	LastUpdate int64
}

func (r *UpdateRunner) Run() error {
	due, err := r.isUpdateDue()
	if !due {
		Debug("Update not due.")
		return nil
	}
	if err != nil {
		Debug("%+v", err)
		return err
	}

	li, err := r.LatestInfo()
	if err != nil {
		LogErrM("Couldn't get remote release info.", err)
		return err
	}
	if li.Version == r.currentVersion {
		Debug("Local is same as latest version: %s", li.Version)
		r.recordUpdate()
		return nil
	}
	latest, err := r.LatestBinary()
	if err != nil {
		LogErrM("Couldn't get latest from remote.", err)
		return err
	}
	defer latest.Close() //TODO: Currently NOOP, should be real closer
	Debug("Applying update.")
	err = r.Applier(latest, gu.Options{})
	if err != nil {
		LogErrM("Couldn't apply update.", err)
		err = r.Rollbacker(err)
		r.CleanUp()
		r.recordUpdate()
		return err
	} else {
		r.CleanUp()
		r.recordUpdate()
		Debug("Updated to version: %s+%s", li.Version, li.Build)
		return nil
	}
}

func (r *UpdateRunner) recordUpdate() error {
	ur := &Record{LastUpdate: time.Now().Unix()}
	Debug("Updating update record.")
	return r.Persister.Upsert(RecordFile, ur)
}

func (r *UpdateRunner) isUpdateDue() (bool, error) {
	if !Prefs().RegulateUpdates {
		Debug("Updates not regulated.")
		return true, nil
	}
	ur := &Record{}
	hiatus := time.Now().Unix() - int64(r.UpdatePeriod/time.Second)
	Debug("Checking update is newer than: %d (Unix time seconds)", hiatus)
	if err := r.Persister.Get(RecordFile, ur, hiatus); err != nil {
		LogErrM("Couldn't get local update record.", err)
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
		LogErrM("If you are running nacl or OpenBSD they are not supported.", err)
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
		LogErrM("Couldn't execute command.", err)
	}
	out.Close()
}
