package sys

const (
	OS_DARWIN = `darwin`
	OS_LINUX = `linux`
	OS_WINDOWS = `windows`
	KWK_TESTMODE = `KWK_TESTMODE`
)

type Manager interface {
	CopyToClipboard(input string) error
	Exists(path string) (bool, error)

	Delete(subDirName string, suffixPath string) error
	FileExists(subDirName string, suffixPath string, incHoldingDir bool) (bool, error)
	WriteToFile(subDirName string, suffixPath string, snippet string, incHoldingDir bool) (string, error)
	ReadFromFile(subDirName string, suffixPath string, incHoldingDir bool, fresherThan int64) (string, error)
}
