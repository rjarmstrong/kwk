package out

import (
	"bitbucket.com/sharingmachine/kwkcli/src/style"
	"bitbucket.com/sharingmachine/types"
	"bitbucket.com/sharingmachine/types/vwrite"
	"fmt"
	"io"
	"time"
)

func Version(i types.AppInfo) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprint(w, style.Fmt256(colors.RecentPouch, "kwk version: "))
		fmt.Fprintf(w, "%s", i.Version)
	})
}

func Upgraded(i types.AppInfo) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintln(w, "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Fprintln(w, "   kwk successfully upgraded!  ")
		fmt.Fprintln(w, "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Fprintf(w, "%s released %s\n", i.Version, style.Time(time.Unix(i.Time, 0)))
	})
}
