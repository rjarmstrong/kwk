package sys

import (
	"github.com/kennygrant/sanitize"
	"io/ioutil"
	"strings"
	"path"
	"os"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/log"
	"bitbucket.com/sharingmachine/kwkcli/cache"
)

var ErrFileNotFound = models.ErrOneLine(models.Code_NotFound, "File not found.")
var ErrFileExpired = models.ErrOneLine(models.Code_NotFound, "File found but expired.")

func New() Manager {
	return &StdManager{}
}

type StdManager struct {
}

func (s *StdManager) WriteToFile(subDirName string, suffixPath string, snippet string, incHoldingDir bool) (filePath string, err error) {
	fp := s.getFilePath(subDirName, suffixPath, incHoldingDir)
	log.Debug("WRITE: %s", fp)
	err = ioutil.WriteFile(fp, []byte(snippet), cache.StandardFilePermission)
	return fp, err
}

func (s *StdManager) ReadFromFile(subDirName string, suffixPath string, incHoldingDir bool, after int64) (string, error) {
	fp := s.getFilePath(subDirName, suffixPath, incHoldingDir)
	log.Debug("READ: %s", fp)
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
	if err := os.MkdirAll(dir, cache.StandardFilePermission); err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func (s *StdManager) getSubDir(directoryName string) (string, error) {
	dir := path.Join(cache.Path(), directoryName)
	err := s.upsertDirectory(dir)
	return directoryName, err
}

func (s *StdManager) getHoldingDirectory(subDirName string, fullName string) string {
	hd := strings.Replace(fullName, ".", "_", -1)
	dirPath := path.Join(cache.Path(), subDirName, hd)
	if e := s.upsertDirectory(dirPath); e != nil {
		panic(e)
	}
	return hd
}

func (s *StdManager) getFilePath(subDirName string, suffixPath string, incHoldingDir bool) string {
	sn := sanitize.Name(suffixPath)
	if incHoldingDir {
		hd := s.getHoldingDirectory(subDirName, sn)
		return path.Join(cache.Path(), subDirName, hd, sn)
	} else {
		if err := s.upsertDirectory(path.Join(cache.Path(), subDirName)); err != nil {
			panic(err)
		}
		return path.Join(cache.Path(), subDirName, sn)
	}
}
