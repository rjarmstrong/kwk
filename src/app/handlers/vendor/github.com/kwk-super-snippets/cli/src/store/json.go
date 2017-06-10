package store

import "encoding/json"

/*
DocStore is a simplified file system interface, the idea is that you can read and write to the native file system with json or whichever format you like.
*/
type Doc interface {
	Upsert(fullKey string, data interface{}) error
	Get(fullKey string, value interface{}, fresherThan int64) error
	Delete(fullKey string) error
	DeleteAll() error
}

// FileSettings is an abstraction over the file system which
// marshals json to and from a specific location.
type Json struct {
	DirectoryName string
	File
}

func NewJson(f File, subDirName string) Doc {
	return &Json{DirectoryName: subDirName, File: f}
}

func (s *Json) Upsert(key string, value interface{}) error {
	bytes, _ := json.Marshal(value)
	_, err := s.File.Write(s.DirectoryName, key, string(bytes), false)
	return err
}

func (s *Json) Get(key string, value interface{}, fresherThan int64) error {
	if str, err := s.File.Read(s.DirectoryName, key, false, fresherThan); err != nil {
		return err
	} else {
		return json.Unmarshal([]byte(str), value)
	}
}

func (s *Json) Delete(key string) error {
	return s.File.Delete(s.DirectoryName, key)
}

func (s *Json) DeleteAll() error {
	return s.File.DeleteAll()
}
