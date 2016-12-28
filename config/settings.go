package config

type Settings interface {
	Upsert(fullKey string, data interface{}) error
	Get(fullKey string, value interface{}) error
	Delete(fullKey string) error
}