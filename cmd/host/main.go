package main

import (
	"github.com/hsfzxjy/zunagi/internal/config"
	"github.com/hsfzxjy/zunagi/internal/errors"
	"github.com/hsfzxjy/zunagi/internal/lifecycle"
	"github.com/hsfzxjy/zunagi/internal/log"
	"github.com/hsfzxjy/zunagi/internal/role"
	"github.com/hsfzxjy/zunagi/internal/ui"
	"golang.design/x/clipboard"

	_ "embed"
)

func main() {
	log.Setup()
	config.Setup(role.Host)
	lifecycle.Setup()
	ui.Setup()

	go func() {
		err := clipboard.Init()
		errors.Check(err)
		var server Server
		server.Run()
	}()

	ui.RunMain()
}
