package system

import "io"

type ISystem interface {
	Upgrade()
	GetVersion() string
	Exists(path string) (bool, error)
	Delete(directoryName string, fullKey string) error
	CopyToClipboard(input string)
	CopyFile(src, dst string) (err error)
	WriteToFile(directoryName string, fullKey string, uri string) (string, error)
	ReadFromFile(directoryName string, fullKey string) (string, error)
	GetDirPath(directoryName string) (string, error)
	ExecSafe(name string, arg ...string) io.ReadCloser
}