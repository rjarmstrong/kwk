package types

import (
	"fmt"
	"github.com/rjarmstrong/kwk-types/errs"
	"github.com/rjarmstrong/kwk-types/valid"
	"net/url"
	"strconv"
	"strings"
)

func NewAlias(username string, pouch string, name string, ext string) *Alias {
	if pouch == "" {
		pouch = PouchRoot
	}
	return &Alias{
		Username: strings.ToLower(username),
		Pouch:    strings.ToLower(pouch),
		Name:     strings.ToLower(name),
		Ext:      strings.ToLower(ext),
	}
}

func NewVersionAlias(username string, pouch string, name string, ext string, version int64) *Alias {
	return &Alias{
		Username: strings.ToLower(username),
		Pouch:    strings.ToLower(pouch),
		Name:     strings.ToLower(name),
		Ext:      strings.ToLower(ext),
		Version:  version,
	}
}

func (m *Alias) IsAbs() bool {
	return m.Username != ""
}

func (m *Alias) IsSnippet() bool {
	return m.Name != ""
}

func (m *Alias) IsPouch() bool {
	return m.Pouch != "" && !m.IsSnippet()
}

func (m *Alias) IsUsername() bool {
	return !m.IsSnippet() && !m.IsPouch() && m.Username != ""
}

func (m *Alias) IsEmpty() bool {
	return !m.IsSnippet() && !m.IsPouch() && m.Username == ""
}

func (m *Alias) URI() string {
	var segs []string
	if m.Username != "" {
		segs = append(segs, "/"+m.Username)
	}
	if m.Pouch != "" {
		segs = append(segs, m.Pouch)
	}
	if m.IsSnippet() {
		segs = append(segs, m.FileName())
	}
	res := strings.Join(segs, "/")
	if m.IsPouch() || m.IsUsername() || m.IsEmpty() {
		return res + "/"
	}
	return res
}

func (m *Alias) VersionURI() string {
	if m.Version == 0 {
		return m.URI()
	}
	return fmt.Sprintf("%s?v=%d", m.URI(), m.Version)
}

func NewSnipName(name string, extension string) SnipName {
	return SnipName{Name: name, Ext: extension}
}

func (m *Alias) FileName() string {
	if m.Ext == "" {
		return m.Name
	}
	return fmt.Sprintf("%s.%s", m.Name, m.Ext)
}

func (m *SnipName) FileName() string {
	if m.Ext == "" {
		return m.Name
	}
	return fmt.Sprintf("%s.%s", m.Name, m.Ext)
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
	version := int64(0)
	if len(query) > 1 {
		v, err := url.ParseQuery(query[1])
		if err != nil {
			return nil, err
		}
		if len(v["v"]) > 0 {
			v, err := strconv.Atoi(v["v"][0])
			if err != nil {
				return nil, err
			}
			version = int64(v)
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

// IsDefNotPouchedSnippetURI determines the likely-hood that a string is an alias
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
		sn = append(sn, &SnipName{a.Name, a.Ext})
	}
	return sn, pouch, nil
}

// GetAliases Returns a slice of aliases given one username, pouch and a slice of snip names.
func GetAliases(username string, pouch string, sn []SnipName) []string {
	var as []string
	for _, v := range sn {
		a := NewAlias(username, pouch, v.Name, v.Ext)
		as = append(as, a.URI())
	}
	return as
}
