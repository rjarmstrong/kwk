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

type Root struct {
	HidePrivate bool
	Username string
	Pouches []*Pouch
	Personal []*Pouch
	Snippets []*Snippet
}

func (rt *Root) IsPouch(name string) bool {
	for _, v := range rt.Pouches {
		if name == v.Name {
			return true
		}
	}
	return false
}
