package persist

import (
	"encoding/json"
)

// FileSettings is an abstraction over the file system which
// marshals json to and from a specific location.
type Json struct {
	DirectoryName string
	System        IO
}

func NewJson(s IO, subDirName string) Persister {
	return &Json{DirectoryName: subDirName, System: s}
}

func (s *Json) Upsert(key string, value interface{}) error {
	bytes, _ := json.Marshal(value)
	_, err := s.System.Write(s.DirectoryName, key, string(bytes), false)
	return err
}

func (s *Json) Get(key string, value interface{}, fresherThan int64) error {
	if str, err := s.System.Read(s.DirectoryName, key, false, fresherThan); err != nil {
		return err
	} else {
		return json.Unmarshal([]byte(str), value)
	}
}

func (s *Json) Delete(key string) error {
	return s.System.Delete(s.DirectoryName, key)
}

func (s *Json) DeleteAll() error {
	return s.System.DeleteAll()
}