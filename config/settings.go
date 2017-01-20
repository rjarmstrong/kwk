package config

import (
)

type Settings interface {
	Upsert(fullKey string, data interface{}) error
	Get(fullKey string, value interface{}, fresherThan int64) error
	Delete(fullKey string) error
	SetPersistedPrefs(pp *PersistedPrefs)
}