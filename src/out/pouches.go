package out

import (
	"fmt"
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"io"
)

func PouchNotDeleted(name string) vwrite.Handler {
	return Info(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s was NOT deleted.\n", name)
	}))
}

func PouchDeleted(name string) vwrite.Handler {
	return Info(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s and contents deleted.\n", name)
	}))
}

func PouchCheckDelete(name string) vwrite.Handler {
	return Warn(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Are you sure you want to delete pouch %s? [y/n] ", name)
	}))
}

func PouchCreated(name string) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Pouch %s created\n", name)
	}))
}

func PouchRenamed(from string, to string) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Pouch %s renamed to %s\n", from, to)
	}))
}

func PouchLocked(name string) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Pouch %s locked. All snippets inside are %s PRIVATE.\n", style.IconPrivatePouch, name)
	}))
}

func PouchCheckUnLock(name string) vwrite.Handler {
	return Success(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Are you sure you want pouch %s public ? [y/n] ", name)
	}))
}

func PouchUnLocked(name string) vwrite.Handler {
	return Warn(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Pouch %s unlocked. All snippets inside are now %s PUBLIC.\n", style.IconPouch, name)
	}))
}

func PouchNotUnLocked(name string) vwrite.Handler {
	return Info(vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Pouch %s  %s NOT unlocked.\n", style.IconPouch, name)
	}))
}
