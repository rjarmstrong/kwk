package updater

import (
	"bytes"
	gu "github.com/inconshreveable/go-update"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"io"
	"os"
	"os/exec"
	"time"
	"github.com/kwk-super-snippets/cli/src/store"
)

// SilentUpdate spawns a new process to check for updates and runs.
func SilentUpdate() {
	cmd, err := os.Executable()
	out.Debug("Initiating silent update check for: %s", cmd)
	if err != nil {
		out.Debug("If you are running nacl or OpenBSD they are not supported.")
		out.LogErr(err)
	}
	exe(false, cmd, "update", "silent")
}

type Runner interface {
	Run() error
}

func New(version string) Runner {
	return &runner{currentVersion: version, BinRepo: &S3Repo{}, Applier: gu.Apply, Rollbacker: gu.RollbackError, UpdatePeriod: time.Hour}
}

type Applier func(update io.Reader, opts gu.Options) error

type Rollbacker func(err error) error

type Record struct {
	LastUpdate int64
}


type runner struct {
	UpdatePeriod   time.Duration
	BinRepo
	Applier
	Rollbacker
	currentVersion string
	store.Doc
}

func (r *runner) Run() error {
	li, err := r.LatestInfo()
	if err != nil {
		out.LogErrM("Couldn't get remote release info.", err)
		return err
	}
	if li.Version == r.currentVersion {
		out.Debug("Local is same as latest version: %s", li.Version)
		r.recordUpdate()
		return nil
	}
	latest, err := r.LatestBinary()
	if err != nil {
		out.LogErrM("Couldn't get latest from remote.", err)
		return err
	}
	defer latest.Close() //TODO: Currently NOOP, should be real closer
	out.Debug("Applying update.")
	err = r.Applier(latest, gu.Options{})
	if err != nil {
		out.LogErrM("Couldn't apply update.", err)
		err = r.Rollbacker(err)
		r.CleanUp()
		r.recordUpdate()
		return err
	} else {
		r.CleanUp()
		r.recordUpdate()
		out.Debug("Updated to version: %s+%s", li.Version, li.Build)
		return nil
	}
}

func (r *runner) recordUpdate() error {
	out.Debug("Updated: %+v", Record{LastUpdate: time.Now().Unix()})
	return nil
}

func exe(wait bool, name string, arg ...string) {
	c := exec.Command(name, arg...)
	c.Stdin = os.Stdin
	stdOut, err := c.StdoutPipe()
	if err != nil {
		out.LogErrM("If you are running nacl or OpenBSD they are not supported.", err)
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
		out.LogErrM("Couldn't execute command.", err)
	}
	stdOut.Close()
}
