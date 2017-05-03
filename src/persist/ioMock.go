package persist

type IoMock struct {
	VersionCalled             bool
	CopyToClipboardCalledWith string
}

func (*IoMock) Delete(subDirName string, suffixPath string) error {
	panic("implement me")
}

func (*IoMock) Write(subDirName string, suffixPath string, snippet string, incHoldingDir bool) (string, error) {
	panic("implement me")
}

func (*IoMock) Read(subDirName string, suffixPath string, incHoldingDir bool, fresherThan int64) (string, error) {
	panic("implement me")
}

func (*IoMock) DeleteAll() error {
	panic("implement me")
}
