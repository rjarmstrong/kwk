package system

import (
	"encoding/json"
)

type Settings struct {
	DirectoryName string
}

func NewSettings(directoryName string) *Settings {
	return &Settings{ DirectoryName: directoryName}
}


func (s *Settings) Upsert(key string, value interface{}) {
	bytes, _ := json.Marshal(value)
	WriteToFile(s.DirectoryName, key, string(bytes))
}

func (s *Settings) Get(key string, value interface{}) error {
	str, err := ReadFromFile(s.DirectoryName, key)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(str), value)
	return err
}

func (s *Settings) Delete(key string) error {
	return Delete(s.DirectoryName, key)
}