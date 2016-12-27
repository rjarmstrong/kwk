package system

import "io"

type SystemMock struct {
	UpgradeCalled             bool
	VersionCalled             bool
	CopyToClipboardCalledWith string
}

func (s *SystemMock) Upgrade() error {
	s.UpgradeCalled = true
	return nil
}

func (s *SystemMock) GetVersion() (string, error) {
	s.VersionCalled = true
	return "0.0.1", nil
}

func (s *SystemMock) CopyToClipboard(input string) error {
	s.CopyToClipboardCalledWith = input
	return nil
}

func (s *SystemMock) Exists(path string) (bool, error) {
	//return s.ExistsResponse
	return false, nil
}
func (s *SystemMock) CopyFile(src, dst string) error {
	return nil
}
func (s *SystemMock) Delete(directoryName string, fullKey string) error {
	return nil
}
func (s *SystemMock) WriteToFile(directoryName string, fullKey string, uri string) (string, error) {
	return "", nil
}
func (s *SystemMock) ReadFromFile(directoryName string, fullKey string) (string, error) {
	return "", nil
}
func (s *SystemMock) GetDirPath(directoryName string) (string, error) {
	return "", nil
}
func (s *SystemMock) ExecSafe(name string, arg ...string) io.ReadCloser {
	return nil
}
