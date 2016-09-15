package system

type ISystem interface {
	Upgrade()
	GetVersion() string
	ChangeDirectory(username string)
}
