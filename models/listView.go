package models

import "time"

/*
ListView represents a listing view which could contain snippets, pouches or both.
Note 'Pouch' is only the root pouch when 'IsRoot' is true and not simply when it
is an empty string.
 */
type ListView struct {
	IsRoot bool
	Pouch    *Pouch

	LastUpgrade int64
	Expanded    bool

	Username    string
	Pouches     []*Pouch
	Personal    []*Pouch
	Snippets    []*Snippet
	UserStats

	Total    int64
	Since    time.Time
	Size     int64
}
func (rt *ListView) GetPouch(name string) *Pouch {
	for _, v := range rt.Pouches {
		if name == v.Name {
			return v
		}
	}
	for _, v := range rt.Personal {
		if name == v.Name {
			return v
		}
	}
	return nil
}

func (rt *ListView) IsPouch(name string) bool {
	for _, v := range rt.Pouches {
		if name == v.Name {
			return true
		}
	}
	for _, v := range rt.Personal {
		if name == v.Name {
			return true
		}
	}
	return false
}
