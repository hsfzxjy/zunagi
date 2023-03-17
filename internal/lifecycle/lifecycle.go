package lifecycle

import (
	"context"

	"github.com/gotk3/gotk3/gtk"
)

var (
	diedCtx, diedCancel = context.WithCancel(context.Background())
)

func DiedCtx() context.Context { return diedCtx }
func Shutdown() {
	gtk.MainQuit()
	diedCancel()
}
