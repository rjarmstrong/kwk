package types

import (
	"github.com/gocql/gocql"
	"time"
)

type Pouch struct {
	Id          string
	Username    string
	Name        string
	MakePrivate bool // All items that are place in this pouch are MADE private
	Encrypt     bool // All items that are placed in this pouch are encrypted
	PouchCounts
	SharedWith []string //list of usernames which can view this if private.
	Modified   time.Time
	LastUse    time.Time
	UnOpened   int64
	Personal   bool
	Type       PouchType
}

type PouchStats struct {
	Username string
	PouchId  gocql.UUID
	PouchCounts
}

type PouchCounts struct {
	Red    int64
	Green  int64
	Clones int64
	Views  int64
	Runs   int64
	Snips  int64
}

func (pc PouchCounts) Use() int64 {
	return (pc.Clones * 10) + (pc.Runs * 6) + (pc.Views * 4)
}

type UserStats struct {
	LastPouch        string
	MaxUsePerPouch   int64
	MaxSnipsPerPouch int64
	RecentPouches    []string

	//AppCount int64
	//Snips int64
	//Pouches int64
	//PrivatePouches int64
}
