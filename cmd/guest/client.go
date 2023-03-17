package main

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/hsfzxjy/zunagi/internal/config"
	"github.com/hsfzxjy/zunagi/internal/lifecycle"
	"github.com/hsfzxjy/zunagi/internal/log"
)

type Client struct{}

func (c Client) connectHost(ctx context.Context, wg *sync.WaitGroup, hostAddress string) {
	var (
		err    error
		conn   net.Conn
		dialer = net.Dialer{Timeout: 5 * time.Second}
	)

	defer wg.Done()

	for {
		log.Infof("connecting to host at %s", hostAddress)
		conn, err = dialer.DialContext(ctx, "tcp", hostAddress)
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err != nil {
			log.Warnf("cannot connect to host: %s, retry in 5 seconds", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Infof("connected to host at %s", hostAddress)
		c.handleConnection(ctx, conn)
	}
}

func (c Client) Run() {
	observer, cancel := config.HostAddress.Listen()
	defer cancel()
	go func() {
		<-lifecycle.DiedCtx().Done()
		cancel()
	}()
	var (
		clientCancel = func() {}
		clientCtx    context.Context
		wg           sync.WaitGroup
	)

	cancelAndWait := func() {
		if clientCancel != nil {
			clientCancel()
			wg.Wait()
		}
	}

	for address := range observer {
		cancelAndWait()
		clientCtx, clientCancel = context.WithCancel(lifecycle.DiedCtx())
		wg.Add(1)
		go c.connectHost(clientCtx, &wg, address)
	}
	cancelAndWait()
}
