package config

import (
	"reflect"
)

type PersisterMock struct {
	GetCalledWith             []interface{}
	ChangeDirectoryCalledWith string
	UpsertCalledWith          []interface{}
	DeleteCalledWith          string
	GetHydrates               interface{}
	GetReturns               error
}

func (s *PersisterMock) Delete(fullKey string) error {
	s.DeleteCalledWith = fullKey
	return nil
}

func (s *PersisterMock) Get(fullKey string, input interface{}, fresherThan int64) error {
	s.GetCalledWith = []interface{}{fullKey, input}
	if s.GetReturns != nil {
		return s.GetReturns
	}
	if s.GetHydrates != nil {
		v1 := reflect.ValueOf(input).Elem()
		v2 := reflect.ValueOf(s.GetHydrates).Elem()
		v1.Set(v2)
	}
	return nil
}

func (s *PersisterMock) Upsert(dir string, data interface{}) error {
	s.UpsertCalledWith = []interface{}{dir, data}
	return nil
}
