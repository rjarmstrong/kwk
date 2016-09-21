package settings

type ISettings interface {
	ChangeDirectory(username string)
	Upsert(fullKey string, data interface{})
	Get(fullKey string, value interface{}) error
	Delete(fullKey string) error
}
