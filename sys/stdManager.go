package sys

import (
	"github.com/kennygrant/sanitize"
	"io/ioutil"
	"runtime"
	"strings"
	"path"
	"log"
	"fmt"
	"os"
	"os/user"
	"bitbucket.com/sharingmachine/kwkcli/models"
	lg "bitbucket.com/sharingmachine/kwkcli/log"
)

const (
	StandardFilePermission = 0700
)

var ErrFileNotFound = models.ErrOneLine(models.Code_NotFound, "File not found.")
var ErrFileExpired = models.ErrOneLine(models.Code_NotFound, "File found but expired.")

func New() Manager {
	return &StdManager{}
}

// StdManager is the default sys.Manager
type StdManager struct {
}

func NewLogger() (*os.File, *log.Logger) {
	f, err := os.OpenFile(path.Join(getCachePath(), "kwk.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger := log.New(f, "> ", log.LstdFlags)
	return f, logger
}

func (s *StdManager) WriteToFile(subDirName string, suffixPath string, snippet string, incHoldingDir bool) (filePath string, err error) {
	fp := s.getFilePath(subDirName, suffixPath, incHoldingDir)
	lg.Debug("Writing file: %s", fp)
	err = ioutil.WriteFile(fp, []byte(snippet), StandardFilePermission)
	return fp, err
}

func (s *StdManager) ReadFromFile(subDirName string, suffixPath string, incHoldingDir bool, after int64) (string, error) {
	fp := s.getFilePath(subDirName, suffixPath, incHoldingDir)
	lg.Debug("Reading file: %s", fp)
	if fi, err := os.Stat(fp); err != nil {
		if os.IsNotExist(err) {
			// TODO: PUT IN STANDARD ERROR
			return "", ErrFileNotFound
		} else {
			return "", err
		}
	} else {
		if after == 0 || after < int64(fi.ModTime().Unix()) {
			bts, err := ioutil.ReadFile(fp)
			return string(bts), err
		} else {
			return "", ErrFileExpired
		}

	}
}

func (s *StdManager) FileExists(subDirName string, suffixPath string, incHoldingDir bool) (bool, error) {
	fp := s.getFilePath(subDirName, suffixPath, incHoldingDir)
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

func (s *StdManager) Delete(directoryName string, fullKey string) error {
	dirPath, err := s.getSubDir(directoryName)
	if err != nil {
		return err
	}
	fp := path.Join(dirPath, fullKey)
	return os.RemoveAll(fp)
}

func (s *StdManager) upsertDirectory(dir string) error {
	if err := os.MkdirAll(dir, StandardFilePermission); err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func (s *StdManager) getSubDir(directoryName string) (string, error) {
	dir := path.Join(getCachePath(), directoryName)
	err := s.upsertDirectory(dir)
	return directoryName, err
}

func (s *StdManager) getHoldingDirectory(subDirName string, fullName string) string {
	hd := strings.Replace(fullName, ".", "_", -1)
	dirPath := path.Join(getCachePath(), subDirName, hd)
	if e := s.upsertDirectory(dirPath); e != nil {
		panic(e)
	}
	return hd
}

func (s *StdManager) getFilePath(subDirName string, suffixPath string, incHoldingDir bool) string {
	sn := sanitize.Name(suffixPath)
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
	var p string
	u, err := user.Current()
	if err != nil {
		// TODO: Write friendly
		panic(err)
	}
	if runtime.GOOS == OS_WINDOWS {
		// use AppDir instead
		p = "%LocalAppData%\\kwk"
	} else if runtime.GOOS == OS_LINUX {
		p = fmt.Sprintf("/%s/.kwk", u.Username)
	} else if runtime.GOOS == OS_DARWIN {
		if u.Username == "root" {
			p = "/var/root/.kwk"
		} else {
			p = fmt.Sprintf("/Users/%s/.kwk", u.Username)
		}
	} else {
		// TODO: Write friendly
		panic("OS not supported.")
	}
	if err := os.MkdirAll(p, StandardFilePermission); err != nil {
		if os.IsExist(err) {
			return p
		}
		panic(err)
	}
	return p
}
