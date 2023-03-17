package ui

import (
	_ "embed"

	"github.com/getlantern/systray"
	"github.com/gotk3/gotk3/glib"
	"github.com/hsfzxjy/zunagi/internal/lifecycle"
)

//go:embed assets/tray-connected.ico
var trayConnectedIcon []byte

func setupSystray() {
	systray.Register(func() {
		systray.SetIcon(trayConnectedIcon)
		systray.SetTooltip("Zunagi")
		settingsItem := systray.AddMenuItem("Settings", "")
		_ = settingsItem
		exitItem := systray.AddMenuItem("Exit", "")
		_ = exitItem
		go func() {
			for {
				select {
				case <-settingsItem.ClickedCh:
					glib.IdleAdd(func() {
						SettingsWindow.Show()
					})
				case <-exitItem.ClickedCh:
					lifecycle.Shutdown()
				}
			}
		}()
	}, nil)
}
