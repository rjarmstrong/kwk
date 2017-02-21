package sys

var KWK_TEST_MODE = false

var Version string

var Build string

type Manager interface {
	Exists(path string) (bool, error)
	Delete(subDirName string, suffixPath string) error
	FileExists(subDirName string, suffixPath string, incHoldingDir bool) (bool, error)
	WriteToFile(subDirName string, suffixPath string, snippet string, incHoldingDir bool) (string, error)

	/*
	ReadFromFile fresherThan = get record as long as it was last modified after this unix time value in seconds.
	 */
	ReadFromFile(subDirName string, suffixPath string, incHoldingDir bool, fresherThan int64) (string, error)
}
