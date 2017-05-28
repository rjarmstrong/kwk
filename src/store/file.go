package app

import (
	"github.com/kennygrant/sanitize"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/types/errs"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type IO interface {
	Delete(subDirName string, suffixPath string) error
	Write(subDirName string, suffixPath string, snippet string, incHoldingDir bool) (string, error)
	// ReadFromFile fresherThan = get record as long as it was last modified after this unix time value in seconds.
	Read(subDirName string, suffixPath string, incHoldingDir bool, fresherThan int64) (string, error)
	DeleteAll() error
}

func NewIO() IO {
	return &StdIO{}
}

type StdIO struct {
}

func (s *StdIO) Write(subDirName string, suffixPath string, snippet string, incHoldingDir bool) (filePath string, err error) {
	fp := s.getFilePath(subDirName, suffixPath, incHoldingDir)
	out.Debug("WRITE: %s", fp)
	err = ioutil.WriteFile(fp, []byte(snippet), out.StandardFilePermission)
	return fp, err
}

func (s *StdIO) Read(subDirName string, suffixPath string, incHoldingDir bool, after int64) (string, error) {
	fp := s.getFilePath(subDirName, suffixPath, incHoldingDir)
	out.Debug("READ: %s", fp)
	if fi, err := os.Stat(fp); err != nil {
		if os.IsNotExist(err) {
			return "", errs.FileNotFound
		} else {
			return "", err
		}
	} else {
		if after == 0 || after < int64(fi.ModTime().Unix()) {
			bts, err := ioutil.ReadFile(fp)
			return string(bts), err
		} else {
			return "", errs.FileExpired
		}

	}
}

func (s *StdIO) DeleteAll() error {
	return os.RemoveAll(out.KwkPath())
}

func (s *StdIO) Delete(directoryName string, fileName string) error {
	dirPath, err := s.getSubDir(directoryName)
	if err != nil {
		return err
	}
	fp := path.Join(dirPath, fileName)
	out.Debug("DELETING:%s", fp)
	return os.RemoveAll(fp)
}

func (s *StdIO) upsertDirectory(dir string) error {
	if err := os.MkdirAll(dir, out.StandardFilePermission); err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

// getSubDir gets the directory immediately below the root (~/.kwk/<sub dir>)
func (s *StdIO) getSubDir(directoryName string) (string, error) {
	dir := path.Join(out.KwkPath(), directoryName)
	out.Debug("DIR: %s", dir)
	err := s.upsertDirectory(dir)
	return dir, err
}

// getHoldingDirectory gets the directory which holds the file, creates it if it doesn't exist.
func (s *StdIO) getHoldingDirectory(subDirName string, fileName string) string {
	hd := strings.Replace(fileName, ".", "_", -1)
	dirPath := path.Join(out.KwkPath(), subDirName, hd)
	if e := s.upsertDirectory(dirPath); e != nil {
		out.Debug("Could not create directory:")
		out.LogErr(e)
		return ""
	}
	return hd
}

// getFilePath gets the file path of the actual document.
func (s *StdIO) getFilePath(subDirName string, suffixPath string, incHoldingDir bool) string {
	sn := sanitize.Name(suffixPath)
	if incHoldingDir {
		hd := s.getHoldingDirectory(subDirName, sn)
		return path.Join(out.KwkPath(), subDirName, hd, sn)
	} else {
		if err := s.upsertDirectory(path.Join(out.KwkPath(), subDirName)); err != nil {
			panic(err)
		}
		return path.Join(out.KwkPath(), subDirName, sn)
	}
}
