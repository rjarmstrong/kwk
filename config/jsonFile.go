package config

import (
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"encoding/json"
)

// FileSettings is an abstraction over the file system which
// marshals json to and from a specific location.
type JsonFile struct {
	DirectoryName string
	System        sys.Manager
}

func NewJsonSettings(s sys.Manager, subDirName string) Persister {
	return &JsonFile{DirectoryName: subDirName, System: s}
}

func (s *JsonFile) Upsert(key string, value interface{}) error {
	bytes, _ := json.Marshal(value)
	_, err := s.System.WriteToFile(s.DirectoryName, key, string(bytes), false)
	return err
}

func (s *JsonFile) Get(key string, value interface{}, fresherThan int64) error {
	if str, err := s.System.ReadFromFile(s.DirectoryName, key, false, fresherThan); err != nil {
		return err
	} else {
		return json.Unmarshal([]byte(str), value)
	}
}

func (s *JsonFile) Delete(key string) error {
	return s.System.Delete(s.DirectoryName, key)
}