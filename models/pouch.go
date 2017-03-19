package models

import "time"

type Pouch struct {
	Username    string
	Name        string
	MakePrivate bool  // Shape
	Encrypt     bool  // Shape
	PouchStats
	LastUsed      int64  // Last time any snippet was used. = Brightness
	SharedWith []string
	Modified   time.Time
	PouchId    string
	UnOpened   int64
}

type PouchStats struct {
	Views int64
	Runs int64
	Broke int64
	Clones int64
	Snips int64
}
