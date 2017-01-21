package config


type Persister interface {
	Upsert(fullKey string, data interface{}) error
	Get(fullKey string, value interface{}, fresherThan int64) error
	Delete(fullKey string) error
}