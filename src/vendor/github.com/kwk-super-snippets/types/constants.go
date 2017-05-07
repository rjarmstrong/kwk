package types

const (
	MaxProcessLevel      = 3 // The maximum level of embedding within an app
	PrincipalKey         = "principal"
	SnipMaxLength        = 3000
	PreviewMaxRuneLength = 50
	KwkHost              = "kwk.co"

	UseStatusUnknown UseStatus = 0
	UseStatusSuccess UseStatus = 1
	UseStatusFail    UseStatus = 2

	UseTypeUnknown UseType = 0
	UseTypeView    UseType = 1
	UseTypeRun     UseType = 2
	UseTypeClone   UseType = 3

	RoleStandard    SnipRole = 0
	RolePreferences SnipRole = 1
	RoleEnvironment SnipRole = 2
	RoleMessage     SnipRole = 3

	PouchTypePhysical  PouchType = 0
	PouchTypeVirtual   PouchType = 1
	PouchTypePersonal  PouchType = 2
	PouchTypeCommunity PouchType = 3

	PouchSettings = "settings"
	PouchRoot     = ""

	TokenHeaderName = "token"

	OsDarwin  = `darwin`
	OsLinux   = `linux`
	OsWindows = `windows`
)
