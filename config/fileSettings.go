package config

import (
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"encoding/json"
)

type FileSettings struct {
	DirectoryName string
	System        sys.Manager
}

func New(s sys.Manager, directoryName string) *FileSettings {
	return &FileSettings{DirectoryName: directoryName, System: s}
}

func (s *FileSettings) Upsert(key string, value interface{}) error {
	bytes, _ := json.Marshal(value)
	_, err := s.System.WriteToFile(s.DirectoryName, key, string(bytes))
	return err
}

func (s *FileSettings) Get(key string, value interface{}) error {
	str, err := s.System.ReadFromFile(s.DirectoryName, key)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(str), value)
	return err
}

func (s *FileSettings) Delete(key string) error {
	return s.System.Delete(s.DirectoryName, key)
}
