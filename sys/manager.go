package sys

const (
	OS_DARWIN = `darwin`
	OS_LINUX = `linux`
	OS_WINDOWS = `windows`
)

type Manager interface {
	GetVersion() (string, error)
	CopyToClipboard(input string) error
	Exists(path string) (bool, error)

	Delete(subDirName string, fullName string) error
	FileExists(subDirName string, fullName string, incHoldingDir bool) (bool, error)
	WriteToFile(subDirName string, fullName string, snippet string, incHoldingDir bool) (string, error)
	ReadFromFile(subDirName string, fullName string, incHoldingDir bool) (string, error)
}
