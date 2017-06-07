package types

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

func NewBlankSnippet() *Snippet {
	return &Snippet{
		Alias:        &Alias{},
		Content:      "",
		SupportedOn:  &Runtimes{Oss: map[string]bool{}},
		Apps:         &Apps{Aliases: []*Alias{}},
		Dependencies: &Dependencies{Aliases: []*Alias{}},
		Tags:         &Tags{Names: map[string]bool{}},
		Stats:        &SnipStats{},
		Public:       true,
	}
}

func NewSnippet(a *Alias, content string) *Snippet {
	return &Snippet{
		Alias:        a,
		Content:      content,
		SupportedOn:  &Runtimes{Oss: map[string]bool{}},
		Apps:         &Apps{Aliases: []*Alias{}},
		Dependencies: &Dependencies{Aliases: []*Alias{}},
		Tags:         &Tags{Names: map[string]bool{}},
	}
}

func KwkTime(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func (m *Snippet) IsApp() bool {
	if m.Dependencies == nil {
		return false
	}
	return len(m.Dependencies.Aliases) > 0
}

func (m *Snippet) VerifyChecksum() bool {
	s := sha256.Sum256([]byte(m.Content))
	actual := fmt.Sprintf("%x", s)
	return actual == m.Checksum
}

func (m *Snippet) Username() string {
	if m.Alias == nil {
		return ""
	}
	return strings.ToLower(m.Alias.Username)
}

func (m *Snippet) Pouch() string {
	if m.Alias == nil {
		return ""
	}
	return strings.ToLower(m.Alias.Pouch)
}

func (m *Snippet) Name() string {
	if m.Alias == nil {
		return ""
	}
	return strings.ToLower(m.Alias.Name)
}

func (m *Snippet) Ext() string {
	if m.Alias == nil {
		return ""
	}
	return strings.ToLower(m.Alias.Ext)
}

func (m *Snippet) Version() int64 {
	if m.Alias == nil {
		return -1
	}
	return m.Alias.Version
}

func (m *RootResponse) GetPouch(name string) *Pouch {
	for _, v := range m.Pouches {
		if name == v.Name {
			return v
		}
	}
	for _, v := range m.Personal {
		if name == v.Name {
			return v
		}
	}
	return nil
}

func (m *RootResponse) IsPouch(name string) bool {
	p := m.GetPouch(name)
	return p != nil
}

func Limit(in string, length int) string {
	in = strings.Replace(in, "\n", "  ", -1)
	in = strings.TrimSpace(in)
	var numRunes = 0
	b := bytes.Buffer{}
	for _, r := range in {
		if numRunes == length {
			return strings.TrimSpace(b.String())
		}
		numRunes++
		if r == '\n' {
			b.WriteRune(' ')
			b.WriteRune(' ')
			continue
		}
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}
