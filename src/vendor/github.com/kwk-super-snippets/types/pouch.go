package types

import (
	"github.com/gocql/gocql"
	"github.com/satori/go.uuid"
	"time"
)

func NewPouch(username string, name string, mkPrivate bool, personal bool) *Pouch {
	return &Pouch{
		Id:          uuid.NewV4().String(),
		Username:    username,
		Name:        name,
		MakePrivate: mkPrivate,
		Modified:    time.Now().Unix(),
	}
}

func (m *Pouch) Use() int64 {
	if m.Stats != nil {
		return m.Stats.Use()
	}
	return -1
}

type PouchStats struct {
	Username string
	PouchId  gocql.UUID
	PouchCounts
}

func (m *PouchCounts) Use() int64 {
	return (m.Clones * 10) + (m.Runs * 6) + (m.Views * 4)
}
