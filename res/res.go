package res

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed placeholder.png
var b []byte

// Placeholder returns a placeholder resource.
func Placeholder() *fyne.StaticResource {
	return fyne.NewStaticResource("placeholder", b)
}
