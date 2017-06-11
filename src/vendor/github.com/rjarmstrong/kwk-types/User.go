package types

import "strings"

func (m *User) Is(username string) bool {
	return strings.ToLower(m.Username) == strings.ToLower(username)
}
