package models

import "time"

type Pouch struct {
	Username    string
	Name        string
	MakePrivate bool
	Encrypt     bool
	SnipCount   int64
	SharedWith  []string
	Modified    time.Time
	PouchId     string
	UnOpened    int64
}
