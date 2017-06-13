package store

import (
	"github.com/kennygrant/sanitize"
	"github.com/rjarmstrong/kwk-types/errs"
	"github.com/rjarmstrong/kwk/src/out"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type File interface {
	Delete(subDirName string, suffixPath string) error
	Write(subDirName string, suffixPath string, snippet string, incHoldingDir bool) (string, error)
	// ReadFromFile fresherThan = get record as long as it was last modified after this unix time value in seconds.
	Read(subDirName string, suffixPath string, incHoldingDir bool, fresherThan int64) (string, error)
	RmDir(subDirName string, suffixPath string) error
	DeleteAll() error
}

type SnippetReadWriter interface {
	Write(uri string, content string) (string, error)
	Read(uri string) (string, error)
	RmDir(uri string) error
}

func NewDiskFile() File {
	return &diskFile{}
}

type diskFile struct {
}

func (s *diskFile) Write(subDirName string, suffixPath string, content string, incHoldingDir bool) (filePath string, err error) {
	if suffixPath == "" {
		return "", errs.New(errs.CodeInternalError, "file.*Write: No suffixPath provided.")
	}
	fp := s.getFilePath(subDirName, suffixPath, incHoldingDir)
	out.Debug("WRITE: %s", fp)
	err = ioutil.WriteFile(fp, []byte(content), out.StandardFilePermission)
	return fp, err
}

func (s *diskFile) Read(subDirName string, suffixPath string, incHoldingDir bool, after int64) (string, error) {
	fp := s.getFilePath(subDirName, suffixPath, incHoldingDir)
	out.Debug("READ: %s", fp)
	fi, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errs.FileNotFound
		}
		return "", err
	}
	if after == 0 || after < int64(fi.ModTime().Unix()) {
		bts, err := ioutil.ReadFile(fp)
		return string(bts), err
	}
	return "", errs.FileExpired
}

func (s *diskFile) RmDir(subDirName string, suffixPath string) error {
	fp := s.getFilePath(subDirName, suffixPath, true)
	parts := strings.Split(fp, "/")
	// Path will (probably) never be less than 3 segments so not checking length
	dp := strings.Join(parts[0:len(parts)-1], "/")
	out.Debug("EDIT: deleting%s", dp)
	return os.RemoveAll(dp)
}

func (s *diskFile) DeleteAll() error {
	return os.RemoveAll(out.KwkPath())
}

func (s *diskFile) Delete(directoryName string, fileName string) error {
	dirPath, err := s.getSubDir(directoryName)
	if err != nil {
		return err
	}
	fp := path.Join(dirPath, fileName)
	out.Debug("DELETING:%s", fp)
	return os.RemoveAll(fp)
}

func (s *diskFile) upsertDirectory(dir string) error {
	if err := os.MkdirAll(dir, out.StandardFilePermission); err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

// getSubDir gets the directory immediately below the root (~/.kwk/<sub dir>)
func (s *diskFile) getSubDir(directoryName string) (string, error) {
	dir := path.Join(out.KwkPath(), directoryName)
	out.Debug("DIR: %s", dir)
	err := s.upsertDirectory(dir)
	return dir, err
}

// getHoldingDirectory gets the directory which holds the file, creates it if it doesn't exist.
func (s *diskFile) getHoldingDirectory(subDirName string, fileName string) string {
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
func (s *diskFile) getFilePath(subDirName string, suffixPath string, incHoldingDir bool) string {
	sn := sanitize.Name(suffixPath)
	if incHoldingDir {
		hd := s.getHoldingDirectory(subDirName, sn)
		return path.Join(out.KwkPath(), subDirName, hd, sn)
	}
	err := s.upsertDirectory(path.Join(out.KwkPath(), subDirName))
	if err != nil {
		log.Fatal(err)
	}
	return path.Join(out.KwkPath(), subDirName, sn)
}
