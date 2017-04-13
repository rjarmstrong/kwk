package file

var KWK_TEST_MODE = false

type IO interface {
	Delete(subDirName string, suffixPath string) error
	Write(subDirName string, suffixPath string, snippet string, incHoldingDir bool) (string, error)
	// ReadFromFile fresherThan = get record as long as it was last modified after this unix time value in seconds.
	Read(subDirName string, suffixPath string, incHoldingDir bool, fresherThan int64) (string, error)
}
