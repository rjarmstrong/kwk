package settings

type ISettings interface {
	ChangeDirectory(username string) error
	Upsert(fullKey string, data interface{}) error
	Get(fullKey string, value interface{}) error
	Delete(fullKey string) error
}
