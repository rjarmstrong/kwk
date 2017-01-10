package sys

import "io"

type ManagerMock struct {
	UpgradeCalled             bool
	VersionCalled             bool
	CopyToClipboardCalledWith string
}

func (m *ManagerMock) Upgrade() error {
	m.UpgradeCalled = true
	return nil
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
func (m *ManagerMock) WriteToFile(directoryName string, fullKey string, uri string) (string, error) {
	return "", nil
}
func (m *ManagerMock) ReadFromFile(directoryName string, fullKey string) (string, error) {
	return "", nil
}
func (m *ManagerMock) GetDirPath(directoryName string) (string, error) {
	return "", nil
}
func (m *ManagerMock) ExecSafe(name string, arg ...string) io.ReadCloser {
	return nil
}
