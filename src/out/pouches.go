package out

import (
	"github.com/rjarmstrong/kwk-types/vwrite"
	"github.com/rjarmstrong/kwk/src/style"
	"io"
)

func PouchNotDeleted(name string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Info(w, "%s was NOT deleted.\n", name)
	})
}

func PouchDeleted(name string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Info(w, "%s and contents deleted.\n", name)
	})
}

func PouchCheckDelete(name string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Warn(w, "Are you sure you want to delete pouch %s? [y/n] ", name)
	})
}

func PouchCreated(name string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Info(w, "Pouch %s created\n", name)
	})
}

func PouchRenamed(from string, to string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Info(w, "Pouch %s renamed to %s\n", from, to)
	})
}

func PouchLocked(name string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Info(w, "Pouch %s locked. All snippets inside are %s PRIVATE.\n", style.IconPrivatePouch, name)
	})
}

func PouchCheckUnLock(name string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Warn(w, "Are you sure you want pouch %s public ? [y/n] ", name)
	})
}

func PouchUnLocked(name string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Warn(w, "Pouch %s unlocked. All snippets inside are now %s PUBLIC.\n", style.IconPouch, name)
	})
}

func PouchNotUnLocked(name string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Warn(w, "Pouch %s  %s NOT unlocked.\n", style.IconPouch, name)
	})
}
