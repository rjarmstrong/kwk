package out

import (
	"fmt"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk/src/style"
	"time"
)

var (
	AppName        = fmt.Sprintf("%s kwk super snippets", style.Fmt256(style.ColorPouchCyan, style.IconSnippet))
	AppDescription = "A smart & friendly snippet manager for the CLI"
)

func Version(i types.AppInfo) string {
	return fmt.Sprintf("\n\n%s Version : %s\n%s Revision: %s\n%s Released: %s\n",
		style.Margin, i.Version, style.Margin, i.Build, style.Margin,
		time.Unix(i.Time, 0).Format(time.RFC822))

}
