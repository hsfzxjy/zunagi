package ui

import "github.com/gotk3/gotk3/gtk"

func Setup() {
	gtk.Init(nil)
	setupSystray()
}
