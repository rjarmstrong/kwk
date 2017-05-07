package types

import (
	"fmt"
	"github.com/kwk-super-snippets/types/errs"
	"github.com/kwk-super-snippets/types/valid"
	"net/url"
	"strconv"
	"strings"
)

func NewAlias(username string, pouch string, name string, extension string) *Alias {
	if pouch == "" {
		pouch = PouchRoot
	}
	return &Alias{
		Username: username,
		Pouch:    pouch,
		SnipName: SnipName{Name: name, Ext: extension},
	}
}

func NewVersionAlias(username string, pouch string, name string, extension string, version int) *Alias {
	return &Alias{
		Username: username,
		Pouch:    pouch,
		SnipName: SnipName{Name: name, Ext: extension},
		Version:  int64(version),
	}
}

type Alias struct {
	Pouch    string
	Username string
	SnipName
	Version int64
}

func (al *Alias) IsAbs() bool {
	return al.Username != ""
}

func (al *Alias) IsSnippet() bool {
	return al.Name != ""
}

func (al *Alias) IsPouch() bool {
	return al.Pouch != "" && !al.IsSnippet()
}

func (al *Alias) IsUsername() bool {
	return !al.IsSnippet() && !al.IsPouch() && al.Username != ""
}

func (al *Alias) IsEmpty() bool {
	return !al.IsSnippet() && !al.IsPouch() && al.Username == ""
}

func (al *Alias) String() string {
	var segs []string
	if al.Username != "" {
		segs = append(segs, "/"+al.Username)
	}
	if al.Pouch != "" {
		segs = append(segs, al.Pouch)
	}
	if al.IsSnippet() {
		segs = append(segs, al.SnipName.String())
	}
	res := strings.Join(segs, "/")
	if al.IsPouch() || al.IsUsername() || al.IsEmpty() {
		return res + "/"
	}
	return res
}

func NewSnipName(name string, extension string) SnipName {
	return SnipName{Name: name, Ext: extension}
}

type SnipName struct {
	Name string
	Ext  string
}

func (sn *SnipName) String() string {
	if sn.Ext == "" {
		return sn.Name
	}
	return fmt.Sprintf("%s.%s", sn.Name, sn.Ext)
}

// ParseAlias
// [/][<username>][/<pouch>][/][<snippet name>][.ext]
func ParseAlias(uri string) (*Alias, error) {
	if uri == "" {
		return NewAlias("", PouchRoot, "", ""), nil
	}
	uri = strings.ToLower(uri)

	// Trim fully qualified
	if strings.HasPrefix(uri, "kwk.co") {
		uri = strings.TrimPrefix(uri, "kwk.co")
	}

	// Parse query string
	query := strings.Split(uri, "?")
	version := 0
	if len(query) > 1 {
		v, err := url.ParseQuery(query[1])
		if err != nil {
			return nil, err
		}
		if len(v["v"]) > 0 {
			version, err = strconv.Atoi(v["v"][0])
			if err != nil {
				return nil, err
			}
		}
		uri = query[0]
	}

	if uri[0] == '.' {
		uri = strings.TrimPrefix(uri, ".")
	}
	var isAbs bool
	if uri[0] == '/' {
		isAbs = true
		uri = strings.TrimPrefix(uri, "/")
	}
	t := strings.Split(uri, "/")
	lastItem := t[len(t)-1]
	sn, err := ParseSnipName(lastItem)
	if err != nil {
		return nil, err
	}
	if len(t) == 1 {
		if isAbs {
			return NewAlias(uri, PouchRoot, "", ""), nil
		}
		// If its just the name
		return NewVersionAlias("", PouchRoot, sn.Name, sn.Ext, version), nil
	}
	if len(t) == 2 {
		if isAbs {
			return NewVersionAlias(t[0], PouchRoot, sn.Name, sn.Ext, version), nil
		}
		pouch := t[0]
		if !valid.Test(pouch, valid.RgxPouchName) {
			return nil, errs.PouchNameInvalid
		}
		return NewVersionAlias("", pouch, sn.Name, sn.Ext, version), nil
	}
	if len(t) == 3 {
		username := t[0]
		if !valid.Test(username, valid.RgxUsername) {
			return nil, errs.UsernameInvalid
		}
		pouch := t[1]
		if !valid.Test(pouch, valid.RgxPouchName) {
			return nil, errs.PouchNameInvalid
		}
		// If it has three parts then this will be an absolute alias
		return NewVersionAlias(username, pouch, sn.Name, sn.Ext, version), nil
	}
	return nil, errs.AliasTooManySegments
}

// IsAlias determines the likely-hood that a string is an alias
func IsDefNotPouchedSnippetURI(uri string) bool {
	if strings.HasSuffix(uri, "/") {
		return true
	}
	if strings.ContainsAny(uri, "\t\n\f\r ") {
		return true
	}
	slashes := strings.Count(uri, "/")
	if strings.HasPrefix(uri, "/") && slashes < 3 {
		return true
	}
	if slashes < 0 || slashes > 3 {
		return true
	}
	segs := strings.Split(strings.Trim(uri, "/"), "/")
	if len(segs) > 3 && !strings.HasPrefix(uri, KwkHost) {
		return true
	}
	if len(segs) == 1 {
		return true
	}
	return false
}

// ParseSnipName parses a string form of a snip name into a struct.
func ParseSnipName(snipName string) (*SnipName, error) {
	if snipName == "" {
		return &SnipName{}, nil
	}
	lIts := strings.Split(snipName, ".")
	extension := ""
	name := ""
	if len(lIts) > 1 {
		// if there is an extension
		extension = lIts[len(lIts)-1]
		name = strings.TrimSuffix(snipName, "."+extension)
	} else {
		// If there is no extension
		name = snipName
	}
	if !valid.Test(name, valid.RgxSnipName) {
		return nil, errs.SnipNameInvalid
	}
	if extension != "" && !valid.Test(extension, valid.RgxExtension) {
		return nil, errs.ExtensionInvalid
	}
	return &SnipName{Name: name, Ext: extension}, nil
}

// ParseMany Parses many URIs into alias structs.
func ParseMany(uri []string) ([]*SnipName, string, error) {
	sn := []*SnipName{}
	var pouch string
	for i, n := range uri {
		a, err := ParseAlias(n)
		if err != nil {
			return nil, "", err
		}
		if i > 0 && pouch != a.Pouch {
			return nil, "", errs.MultipleTargetPouches
		}
		pouch = a.Pouch
		sn = append(sn, &a.SnipName)
	}
	return sn, pouch, nil
}

// GetAliases Returns a slice of aliases given one username, pouch and a slice of snip names.
func GetAliases(username string, pouch string, sn []SnipName) []string {
	var as []string
	for _, v := range sn {
		a := NewAlias(username, pouch, v.Name, v.Ext)
		as = append(as, a.String())
	}
	return as
}
