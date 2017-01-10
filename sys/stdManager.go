package sys

import (
	"github.com/kennygrant/sanitize"
	"github.com/atotto/clipboard"
	"encoding/json"
	"io/ioutil"
	"runtime"
	"os/exec"
	"bytes"
	"errors"
	"path"
	"log"
	"fmt"
	"io"
	"os"
)

const (
	APP_VERSION = "APP_VERSION"
)

func New() Manager {
	return &StdManager{}
}

// StdManager is the default sys.Manager
type StdManager struct {
}

func NewLogger() (*os.File, *log.Logger) {
	f, err := os.OpenFile(path.Join(GetCachePath(), "kwk.log"), os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger := log.New(f, "> ", log.LstdFlags)
	return f, logger
}

func (s *StdManager) ExecSafe(name string, arg ...string) io.ReadCloser {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	out, _ := cmd.StdoutPipe()
	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}
	return out
}

func PrettyPrint(obj interface{}) {
	fmt.Println("")
	p, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Print(string(p))
	fmt.Print("\n\n")
}

func (s *StdManager) UpsertDirectory(dir string) error {
	ok, err := s.Exists(dir)
	if ok {
		return nil
	}
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}
	return nil
}

func (s *StdManager) WriteToFile(directoryName string, fullKey string, uri string) (string, error) {
	dirPath, err := s.GetDirPath(directoryName)
	filePath := path.Join(dirPath, fullKey, sanitize.Name(fullKey))
	err = ioutil.WriteFile(filePath, []byte(uri), 0666)
	return filePath, err
}

func (s *StdManager) ReadFromFile(directoryName string, fullKey string) (string, error) {
	dirPath, err := s.GetDirPath(directoryName)
	if err != nil {
		return "", err
	}
	fp := path.Join(dirPath, fullKey)
	if ok, _ := s.Exists(fp); !ok {
		return "", errors.New("Not found.")
	}
	bts, err := ioutil.ReadFile(fp)
	return string(bts), err
}

func (s *StdManager) Delete(directoryName string, fullKey string) error {
	dirPath, err := s.GetDirPath(directoryName)
	if err != nil {
		return err
	}
	fp := path.Join(dirPath, fullKey)
	return os.RemoveAll(fp)
}

func (s *StdManager) GetDirPath(directoryName string) (string, error) {
	dir := path.Join(GetCachePath(), directoryName)
	err := s.UpsertDirectory(dir)
	return dir, err
}

func GetCachePath() string {
	p := ""
	u := os.Getenv("USER")
	if runtime.GOOS == "windows" {
		// check that other users can't access this
		p = "%temp%"
	} else if runtime.GOOS == "linux" {
		p = fmt.Sprintf("/%s/.kwk", u)
	} else {
		p = fmt.Sprintf("/Users/%s/Library/Caches/kwk", u)
	}
	if err := os.Mkdir(p, os.ModeDir); err != nil {
		if os.IsExist(err) {
			return p
		}
		panic(err)
	}
	return p
}

func (s *StdManager) Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (s *StdManager) GetVersion() (string, error) {
	if v := os.Getenv(APP_VERSION); v == "" {
		return "", errors.New(APP_VERSION + " has not been set.")
	} else {
		return v, nil
	}
}

func (s *StdManager) Upgrade() error {
	distributionUri := "/Volumes/development/go/src/bitbucket.com/sharingmachine/kwkcli/kwkcli"
	installPath := "/usr/local/bin/kwk"
	//download in future
	return s.CopyFile(distributionUri, installPath)
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func (s *StdManager) CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func (s *StdManager) CopyToClipboard(input string) error {
	return clipboard.WriteAll(input)
}
