package ui

import (
	"sync"

	_ "embed"

	"github.com/gotk3/gotk3/gtk"
	"github.com/hsfzxjy/zunagi/internal/config"
	"github.com/hsfzxjy/zunagi/internal/errors"
)

type _SettingsWindow struct {
	initOnce         sync.Once
	window           *gtk.Window
	entryHostAddress *gtk.Entry
}

var SettingsWindow _SettingsWindow

//go:embed glades/settings.glade
var settingsGlade string

func (s *_SettingsWindow) init() {
	s.initOnce.Do(func() {
		b, err := gtk.BuilderNewFromString(settingsGlade)
		errors.Check(err)
		b.ConnectSignals(map[string]any{
			"window:settings-delete": func() {
				s.window.Hide()
			},
			"window:settings-show": func() {
				s.entryHostAddress.SetText(config.HostAddress.Current())
			},
			"window:settings-apply": func() {
				text, err := s.entryHostAddress.GetText()
				errors.Check(err)
				config.HostAddress.Send(text)
				s.window.Hide()
			},
		})
		window, err := b.GetObject("window:settings")
		errors.Check(err)
		s.window = window.(*gtk.Window)
		entry, err := b.GetObject("entry:settings.host_address")
		errors.Check(err)
		s.entryHostAddress = entry.(*gtk.Entry)
	})
}

func (s *_SettingsWindow) Show() {
	s.init()
	s.window.ShowAll()
}
