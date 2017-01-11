package sys

import (
	"github.com/kennygrant/sanitize"
	"github.com/atotto/clipboard"
	"io/ioutil"
	"runtime"
	"strings"
	"errors"
	"path"
	"log"
	"fmt"
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
	f, err := os.OpenFile(path.Join(getCachePath(), "kwk.log"), os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger := log.New(f, "> ", log.LstdFlags)
	return f, logger
}

func (s *StdManager) WriteToFile(subDirName string, fullName string, snippet string, incHoldingDir bool) (filePath string, err error) {
	fp := s.getFilePath(subDirName, fullName, incHoldingDir)
	err = ioutil.WriteFile(fp, []byte(snippet), 0666)
	return fp, err
}

func (s *StdManager) ReadFromFile(subDirName string, fullName string, incHoldingDir bool) (string, error) {
	fp := s.getFilePath(subDirName, fullName, incHoldingDir)
	if ok, _ := s.FileExists(subDirName, fullName, incHoldingDir); !ok {
		return "", errors.New("Not found.")
	}
	bts, err := ioutil.ReadFile(fp)
	return string(bts), err
}

func (s *StdManager) FileExists(subDirName string, fullName string, incHoldingDir bool) (bool, error) {
	fp := s.getFilePath(subDirName, fullName, incHoldingDir)
	return s.Exists(fp)
}

func (s *StdManager) Exists(fullPath string) (bool, error) {
	if _, err := os.Stat(fullPath); err != nil {
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

func (s *StdManager) CopyToClipboard(input string) error {
	return clipboard.WriteAll(input)
}


func (s *StdManager) Delete(directoryName string, fullKey string) error {
	dirPath, err := s.getSubDir(directoryName)
	if err != nil {
		return err
	}
	fp := path.Join(dirPath, fullKey)
	return os.RemoveAll(fp)
}

func (s *StdManager) upsertDirectory(dir string) error {
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

func (s *StdManager) getSubDir(directoryName string) (string, error) {
	dir := path.Join(getCachePath(), directoryName)
	err := s.upsertDirectory(dir)
	return directoryName, err
}

func(s *StdManager) getHoldingDirectory(subDirName string, fullName string) string {
	hd := strings.Replace(fullName, ".", "_", -1)
	dirPath := path.Join(getCachePath(), subDirName, hd)
	if e := s.upsertDirectory(dirPath); e != nil {
		panic(e)
	}
	return hd
}

func (s *StdManager) getFilePath(subDirName string, fullName string, incHoldingDir bool) string {
	sn := sanitize.Name(fullName)
	if incHoldingDir {
		hd := s.getHoldingDirectory(subDirName, sn)
		return path.Join(getCachePath(), subDirName, hd, sn)
	} else {
		if err := s.upsertDirectory(path.Join(getCachePath(), subDirName)); err != nil {
			panic(err)
		}
		return path.Join(getCachePath(), subDirName, sn)
	}
}

func getCachePath() string {
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
