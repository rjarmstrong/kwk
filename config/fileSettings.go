package config

import (
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"encoding/json"
//	"fmt"
)

// FileSettings is an abstraction over the file system which
// marshals json to and from a specific location.
type FileSettings struct {
	DirectoryName string
	System        sys.Manager
	Prefs         *Preferences
}

func NewFileSettings(s sys.Manager, subDirName string) Settings {
	return &FileSettings{DirectoryName: subDirName, System: s, Prefs:DefaultPrefs()}
}

func (s *FileSettings) Upsert(key string, value interface{}) error {
	bytes, _ := json.Marshal(value)
	_, err := s.System.WriteToFile(s.DirectoryName, key, string(bytes), false)
	return err
}

func (s *FileSettings) Get(key string, value interface{}, fresherThan int64) error {
	if str, err := s.System.ReadFromFile(s.DirectoryName, key, false, fresherThan); err != nil {
		return err
	} else {
		return json.Unmarshal([]byte(str), value)
	}
}

func (s *FileSettings) Delete(key string) error {
	return s.System.Delete(s.DirectoryName, key)
}

func (s *FileSettings) GetPrefs() *Preferences {
	return s.Prefs
}

func (s *FileSettings) SetPersistedPrefs(pp *PersistedPrefs) {
	s.Prefs.PersistedPrefs = *pp
}
