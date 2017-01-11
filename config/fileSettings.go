package config

import (
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"encoding/json"
)

// FileSettings is an abstraction over the file system which
// marshals json to and from a specific location.
type FileSettings struct {
	DirectoryName string
	System        sys.Manager
}

func New(s sys.Manager, subDirName string) *FileSettings {
	return &FileSettings{DirectoryName: subDirName, System: s}
}

func (s *FileSettings) Upsert(key string, value interface{}) error {
	bytes, _ := json.Marshal(value)
	_, err := s.System.WriteToFile(s.DirectoryName, key, string(bytes), false)
	return err
}

func (s *FileSettings) Get(key string, value interface{}) error {
	if str, err := s.System.ReadFromFile(s.DirectoryName, key, false); err != nil {
		return err
	} else {
		return json.Unmarshal([]byte(str), value)
	}
}

func (s *FileSettings) Delete(key string) error {
	return s.System.Delete(s.DirectoryName, key)
}
