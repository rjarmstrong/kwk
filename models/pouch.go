package models

import "time"

type Pouch struct {
	Username    string
	Name        string
	MakePrivate bool  // Shape
	Encrypt     bool  // Shape
	SnipCount   int64  // Fullness (flag number)

	// Trend
	RunCount      int64 // Sum of runs of snippets? = Order
	ViewCount     int64 // Sum of views of snippets? = Order
	// Absolute last used (Yello
	LastUsed      int64  // Last time any snippet was used. = Brightness

	//
	BrokeRate     int64  // broke snippets/total snippets. = Greenness/Redness

	SharedWith []string
	Modified   time.Time
	PouchId    string
	UnOpened   int64
}
