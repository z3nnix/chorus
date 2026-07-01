package chorus

import "github.com/fatih/color"

var (
	cyan   = color.New(color.FgCyan).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

func init() {
	color.NoColor = false
}
