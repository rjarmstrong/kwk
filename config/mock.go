package config

import (
	"reflect"
)

type Mock struct {
	GetCalledWith             []interface{}
	ChangeDirectoryCalledWith string
	UpsertCalledWith          []interface{}
	DeleteCalledWith          string
	GetHydrateWith            interface{}
}

// get and check username exists
// Save to settings
// Print confirmation
//fmt.Println(gui.Colour(gui.LightBlue, "Switched to kwk.co/" + args[0] + "/"))
func (s *Mock) ChangeDirectory(username string) error {
	s.ChangeDirectoryCalledWith = username
	return nil
}

func (s *Mock) Delete(fullKey string) error {
	s.DeleteCalledWith = fullKey
	return nil
}

func (s *Mock) Get(fullKey string, input interface{}) error {
	s.GetCalledWith = []interface{}{fullKey, input}
	if s.GetHydrateWith != nil {
		v1 := reflect.ValueOf(input).Elem()
		v2 := reflect.ValueOf(s.GetHydrateWith).Elem()
		v1.Set(v2)
	}
	return nil
}

func (s *Mock) Upsert(dir string, data interface{}) error {
	s.UpsertCalledWith = []interface{}{dir, data}
	return nil
}
