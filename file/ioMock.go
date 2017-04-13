package file

type IoMock struct {
	VersionCalled             bool
	CopyToClipboardCalledWith string
}


func (*IoMock) FileExists(subDirName string, fullName string, incHoldingDir bool) (bool, error) {
	panic("implement me")
}

func (m *IoMock) GetVersion() (string, error) {
	m.VersionCalled = true
	return "0.0.1", nil
}

func (m *IoMock) Exists(path string) (bool, error) {
	//return s.ExistsResponse
	return false, nil
}
func (m *IoMock) CopyFile(src, dst string) error {
	return nil
}
func (m *IoMock) Delete(directoryName string, fullKey string) error {
	return nil
}
func (m *IoMock) Write(subDirName string, fullName string, snippet string, incHoldingDir bool) (string, error) {
	return "", nil
}
func (m *IoMock) Read(subDirName string, fullName string, incHoldingDir bool, fresherThan int64) (string, error) {
	return "", nil
}