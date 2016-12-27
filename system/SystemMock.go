package system

import "io"

type MockSystem struct {
	UpgradeCalled             bool
	VersionCalled             bool
	CopyToClipboardCalledWith string
}

func (s *MockSystem) Upgrade() error {
	s.UpgradeCalled = true
	return nil
}

func (s *MockSystem) GetVersion() (string, error) {
	s.VersionCalled = true
	return "0.0.1", nil
}

func (s *MockSystem) CopyToClipboard(input string) error {
	s.CopyToClipboardCalledWith = input
	return nil
}

func (s *MockSystem) Exists(path string) (bool, error) {
	//return s.ExistsResponse
	return false, nil
}
func (s *MockSystem) CopyFile(src, dst string) error {
	return nil
}
func (s *MockSystem) Delete(directoryName string, fullKey string) error {
	return nil
}
func (s *MockSystem) WriteToFile(directoryName string, fullKey string, uri string) (string, error) {
	return "", nil
}
func (s *MockSystem) ReadFromFile(directoryName string, fullKey string) (string, error) {
	return "", nil
}
func (s *MockSystem) GetDirPath(directoryName string) (string, error) {
	return "", nil
}
func (s *MockSystem) ExecSafe(name string, arg ...string) io.ReadCloser {
	return nil
}
