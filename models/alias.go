package models

import (
	"fmt"
	"os"
	"strings"
)

func NewAlias(username string, pouch string, name string, extension string) *Alias {
	return &Alias{
		Username:username,
		Pouch:pouch,
		SnipName:SnipName{Name:name, Ext:extension},
	}
}

type Alias struct {
	Pouch string
	Username string `json:"username" schema:"username"`
	SnipName
	FullKey  string `json:"fullKey" schema:"fullKey"`
}

func (a *Alias) String() string {
	return fmt.Sprintf("%s/%s/%s", a.Username, a.Pouch, a.SnipName.String())
}

func (a *Alias) Path() string {
	return fmt.Sprintf("%s/%s", a.Pouch, a.SnipName.String())
}

func NewSetupAlias(name string, ext string) *Alias {
	if h, err := os.Hostname(); err != nil {
		panic(err)
	} else {
		return &Alias{
			Pouch:   SETTINGS_POUCH,
			SnipName:SnipName{Name:fmt.Sprintf("%s_%s", h, name), Ext:ext},
		}
	}
}

type SnipName struct {
	Name string
	Ext  string
}

func (s *SnipName) String() string {
	return fmt.Sprintf("%s.%s", s.Name, s.Ext)
}

// ParseAlias
// <username>/<pouch>/[snippet name].[ext]
// e.g.:
// bill/monkey
// bill/monkey.go
// bill/examples/monkey
// bill/examples/monkey.go
//
// FOR other users:
//
// /mandy/monkey
// /mandy/monkey.go
// /mandy/examples/monkey
// etc.
func ParseAlias(distinctName string) (*Alias, error) {
	if distinctName == "" {
		return NewAlias("", ROOT_POUCH, "", ""), nil
	}
	// When prefixed with a forward slash this refers to another users alias
	var isOtherUserAlias bool
	if distinctName[0] == '/' {
		isOtherUserAlias = true
		distinctName = strings.TrimPrefix(distinctName, "/")
	}
	t := strings.Split(distinctName, "/")
	lastItem := t[len(t)-1]
	sn, err := ParseSnipName(lastItem)
	if err != nil {
		return nil, err
	}
	if len(t) == 1 {
		if isOtherUserAlias {
			return nil, ClientErr{Title:"Incomplete alias for another user must comprise at least /username/snippet"}
		}
		// If its just the name
		return NewAlias("", ROOT_POUCH, sn.Name, sn.Ext), nil
	}
	if len(t) == 2 {
		if isOtherUserAlias {
			return NewAlias(t[0], ROOT_POUCH, sn.Name, sn.Ext), nil
		}
		return NewAlias("", t[0],  sn.Name, sn.Ext), nil
	}
	if len(t) == 3 {
		// If it has three parts then this will be an absolute alias
		return NewAlias(t[0], t[1],  sn.Name, sn.Ext), nil
	}
	return nil, ClientErr{Title:"Invalid snippet reference."}
}

func ParseSnipName (snipName string) (*SnipName, error){
	if snipName == "" {
		return &SnipName{}, nil
	}
	lIts := strings.Split(snipName, ".")
	extension := ""
	name := ""
	if len(lIts) > 1 {
		// if there is an extension
		extension = lIts[len(lIts)-1]
		name = strings.TrimSuffix(snipName, "." + extension)
	} else {
		// If there is no extension
		name = snipName
	}
	return &SnipName{Name:name, Ext:extension}, nil
}

func ParseMany(distinctNames []string) ([]*SnipName, string, error) {
	sn := []*SnipName{}
	var pouch string
	for i, n := range distinctNames {
		a, err := ParseAlias(n)
		if err != nil {
			return nil, "", err
		}
		if i > 0 && pouch != a.Pouch {
			return nil, "", ClientErr{Title:"Cannot move or rename snippets from multiple source pouches."}
		}
		pouch = a.Pouch
		sn = append(sn, &a.SnipName)
	}
	return sn, pouch, nil
}
