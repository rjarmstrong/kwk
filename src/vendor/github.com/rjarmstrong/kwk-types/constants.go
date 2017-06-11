package types

const (
	MaxProcessLevel      = 3 // The maximum level of embedding within an app
	PrincipalKey         = "principal"
	SnipMaxLength        = 3000
	PreviewMaxRuneLength = 50
	KwkHost              = "kwk.co"

	PouchSettings = "settings"
	PouchRoot     = ""

	TokenHeaderName  = "x-kwk-access-token"
	Prefs            = "prefs_"
	PrefsPrivateView = Prefs + "private_view"

	OsDarwin  = `darwin`
	OsLinux   = `linux`
	OsWindows = `windows`
)
