package sys

type ManagerMock struct {
	VersionCalled             bool
	CopyToClipboardCalledWith string
}


func (*ManagerMock) FileExists(subDirName string, fullName string, incHoldingDir bool) (bool, error) {
	panic("implement me")
}

func (m *ManagerMock) GetVersion() (string, error) {
	m.VersionCalled = true
	return "0.0.1", nil
}

func (m *ManagerMock) CopyToClipboard(input string) error {
	m.CopyToClipboardCalledWith = input
	return nil
}

func (m *ManagerMock) Exists(path string) (bool, error) {
	//return s.ExistsResponse
	return false, nil
}
func (m *ManagerMock) CopyFile(src, dst string) error {
	return nil
}
func (m *ManagerMock) Delete(directoryName string, fullKey string) error {
	return nil
}
func (m *ManagerMock) WriteToFile(subDirName string, fullName string, snippet string, incHoldingDir bool) (string, error) {
	return "", nil
}
func (m *ManagerMock) ReadFromFile(subDirName string, fullName string, incHoldingDir bool, fresherThan int64) (string, error) {
	return "", nil
}
