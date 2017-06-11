package types

import (
	"github.com/satori/go.uuid"
	"time"
)

func NewBlankPouch() *Pouch {
	return &Pouch{
		SharedWith: &SharedWith{Usernames:map[string]bool{}},
		Stats: &PouchStats{},
	}
}

func NewPouch(username string, name string, mkPrivate bool, personal bool) *Pouch {
	return &Pouch{
		Id:          uuid.NewV4().String(),
		Username:    username,
		Name:        name,
		MakePrivate: mkPrivate,
		Created:    time.Now().Unix(),
		Updated:    time.Now().Unix(),
	}
}

func (m *Pouch) Use() int64 {
	if m.Stats != nil {
		return m.Stats.Use()
	}
	return -1
}

func (m *PouchStats) Use() int64 {
	return (m.Clones * 10) + (m.Runs * 6) + (m.Views * 4)
}
