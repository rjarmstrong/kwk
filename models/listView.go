package models

import "time"

/*
ListView represents a listing view which could contain snippets, pouches or both.
Note 'Pouch' is only the root pouch when 'IsRoot' is true and not simply when it
is an empty string.
 */
type ListView struct {
	IsRoot bool
	Pouch    string

	LastUpdate  int64
	HidePrivate bool
	Expanded bool

	Username    string
	Pouches     []*Pouch
	Personal    []*Pouch
	Snippets    []*Snippet

	Total    int64
	Since    time.Time
	Size     int64
}

func (rt *ListView) IsPouch(name string) bool {
	for _, v := range rt.Pouches {
		if name == v.Name {
			return true
		}
	}
	return false
}
