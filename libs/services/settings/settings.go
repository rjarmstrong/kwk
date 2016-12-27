package settings

import (
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"encoding/json"
)

type Settings struct {
	DirectoryName string
	System        system.ISystem
}

func New(system system.ISystem, directoryName string) *Settings {
	return &Settings{DirectoryName: directoryName, System: system}
}

func (s *Settings) Upsert(key string, value interface{}) error {
	bytes, _ := json.Marshal(value)
	_, err := s.System.WriteToFile(s.DirectoryName, key, string(bytes))
	return err
}

func (s *Settings) Get(key string, value interface{}) error {
	str, err := s.System.ReadFromFile(s.DirectoryName, key)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(str), value)
	return err
}

func (s *Settings) Delete(key string) error {
	return s.System.Delete(s.DirectoryName, key)
}

func (s *Settings) ChangeDirectory(username string) error {
	return nil
}
