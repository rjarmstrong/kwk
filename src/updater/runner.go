package updater

import (
	"bytes"
	gu "github.com/inconshreveable/go-update"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk/src/store"
	"io"
	"os"
	"os/exec"
	"time"
)

const (
	recordFile = `updateRecord`
	Command    = `@update`
)

// SilentUpdate spawns a new process to check for updates and runs.
func SpawnUpdate() {
	cmd, err := os.Executable()
	out.Debug("RUNNER: update in new proc: %s", cmd)
	if err != nil {
		out.Debug("If you are running nacl or OpenBSD they are not supported.")
		out.LogErr(err)
	}
	exe(cmd, Command)
}

// Runner is responsible for holding the state of an update process.
type Runner interface {
	// Run starts the update process.
	Run() error
}

// New creates a new update runner.
func New(version string, repo BinRepo, a applier, rb rollbacker, doc store.Doc) Runner {
	return &runner{doc: doc, currentVersion: version, BinRepo: repo, applier: a, rollbacker: rb, UpdateHiatusSecs: 60 * 5}
}

type applier func(update io.Reader, opts gu.Options) error

type rollbacker func(err error) error

type Record struct {
	LastUpdate int64
}

type runner struct {
	doc store.Doc
	BinRepo
	applier
	rollbacker
	currentVersion string
	store.Doc
	UpdateHiatusSecs int64
}

func (r *runner) Run() error {
	li, err := r.GetLatestInfo()

	if err != nil {
		out.LogErrM("UPDATER: Couldn't get remote release info.", err)
		return err
	}
	if !r.isUpdateDue() {
		out.Debug("UPDATER: Update not due.")
		return nil
	}
	if li.Version == r.currentVersion {
		out.Debug("UPDATER: Local is same as latest version: %s", li.Version)
		r.recordUpdate()
		return nil
	}
	latest, err := r.GetLatestBinary()
	if err != nil {
		out.LogErrM("UPDATER: Couldn't get latest from remote.", err)
		return err
	}
	defer latest.Close() //TODO: Currently NOOP, should be real closer
	out.Debug("UPDATER: Applying update.")
	err = r.applier(latest, gu.Options{})
	if err != nil {
		out.LogErrM("UPDATER: Couldn't apply update.", err)
		err = r.rollbacker(err)
		r.CleanUp()
		r.recordUpdate()
		return err
	}
	r.CleanUp()
	r.recordUpdate()
	out.Debug("UPDATER: Updated to version: %s build: %s", li.Version, li.Build)
	return nil
}

func (r *runner) isUpdateDue() bool {
	ur := &Record{}
	hiatus := time.Now().Unix() - r.UpdateHiatusSecs
	out.Debug("UPDATER: Check update record greater than: %d (Unix time seconds)", hiatus)
	err := r.doc.Get(recordFile, ur, hiatus)
	if err != nil {
		out.Debug("UPDATER: Couldn't get local update record. %s", err)
		// TODO: Force update if there is any error when attempted to find this out.
		out.Debug("UPDATER: Updating...")
		return true
	}
	return false
}

func (r *runner) recordUpdate() error {
	ur := &Record{LastUpdate: time.Now().Unix()}
	out.Debug("UPDATER: saving update record.")
	return r.doc.Upsert(recordFile, ur)
}

func exe(name string, arg ...string) {
	c := exec.Command(name, arg...)
	c.Env = os.Environ()
	stdOut, err := c.StdoutPipe()
	if err != nil {
		out.LogErrM("UPDATER: If you are running nacl or OpenBSD they are not supported.", err)
	}
	defer stdOut.Close()

	var stderr bytes.Buffer
	c.Stdout = os.Stdout
	c.Stderr = &stderr
	err = c.Start()
	if err != nil {
		out.LogErrM("UPDATER: Couldn't execute command.", err)
	}
}
