package main

import (
	"context"
	"net"
	"sync"

	"github.com/hsfzxjy/zunagi/internal/config"
	"github.com/hsfzxjy/zunagi/internal/errors"
	"github.com/hsfzxjy/zunagi/internal/lifecycle"
	"github.com/hsfzxjy/zunagi/internal/log"
)

type Server struct{}

func (s *Server) runServer(ctx context.Context, wg *sync.WaitGroup, address string) {
	var err error

	listener, err := net.Listen("tcp", address)
	errors.Check(err)

	go func() {
		<-ctx.Done()
		listener.Close()
		wg.Done()
	}()

	for {
		log.Infof("waiting for guest to connect at %s", address)
		conn, err := listener.Accept()
		if err != nil {
			log.Warnf("cannot accept connection from guest: %s", err)
			return
		}
		log.Infof("accepted connection from %s", conn.RemoteAddr())
		s.handleConnection(ctx, conn)
	}
}

func (s *Server) Run() {
	observer, cancel := config.HostAddress.Listen()
	defer cancel()
	go func() {
		<-lifecycle.DiedCtx().Done()
		cancel()
	}()
	var (
		serverCancel = func() {}
		serverCtx    context.Context
		wg           sync.WaitGroup
	)

	cancelAndWait := func() {
		if serverCancel != nil {
			serverCancel()
			wg.Wait()
		}
	}

	for address := range observer {
		cancelAndWait()
		serverCtx, serverCancel = context.WithCancel(lifecycle.DiedCtx())
		wg.Add(1)
		go s.runServer(serverCtx, &wg, address)
	}
	cancelAndWait()
}
