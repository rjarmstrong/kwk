package settings

import (
	"encoding/json"
	"github.com/kwk-links/kwk-cli/libs/services/system"
)

type Settings struct {
	DirectoryName string
	System        system.ISystem
}

func NewSettings(system system.ISystem, directoryName string) *Settings {
	return &Settings{DirectoryName: directoryName, System:system}
}

func (s *Settings) Upsert(key string, value interface{}) {
	bytes, _ := json.Marshal(value)
	s.System.WriteToFile(s.DirectoryName, key, string(bytes))
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

func (s *Settings) ChangeDirectory(username string) {

}