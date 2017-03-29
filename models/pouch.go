package models

import "time"

type Pouch struct {
	Username    string
	Name        string
	MakePrivate bool  // Shape
	Encrypt     bool  // Shape
	PouchStats
	LastUse     int64 // Last time any snippet was used. = Brightness
	SharedWith  []string
	Modified    time.Time
	PouchId     string
	UnOpened    int64
	Type        PouchType
}

type PouchType int64

const (
	PouchType_Physical  PouchType = 0
	PouchType_Virtual   PouchType = 1
	PouchType_Personal  PouchType = 2
	PouchType_Community PouchType = 3
)

type PouchStats struct {
	Use int64
	Views int64
	Runs int64
	Green int64
	Red int64
	Clones int64
	Snips int64
}

type UserStats struct {
	LastPouch string
	MaxUsePerPouch int64
	MaxSnipsPerPouch int64
	RecentPouches []string

	//AppCount int64
	//Snips int64
	//Pouches int64
	//PrivatePouches int64
}
