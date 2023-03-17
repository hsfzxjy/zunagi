package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/hsfzxjy/zunagi/internal/errors"
	"github.com/hsfzxjy/zunagi/internal/log"
	"golang.design/x/clipboard"
)

func (c Client) handleConnection(pctx context.Context, conn net.Conn) {
	defer func() {
		if err, ok := recover().(error); ok && err != nil {
			log.Warnf("disconnected: %s", err)
		}
		conn.Close()
	}()
	ctx, cancel := context.WithCancel(pctx)
	defer cancel()
	go func() {
		<-ctx.Done()
		conn.SetDeadline(time.Now())
		conn.Close()
	}()
	receiveAndPaste(conn)
}

func receiveAndPaste(conn net.Conn) {
	var sizeBuffer [8]byte
	var buffer []byte

	for {
		_, err := io.ReadFull(conn, sizeBuffer[:])
		errors.Check(err)
		size := binary.BigEndian.Uint64(sizeBuffer[:])

		// handle heartbeat
		if size == 0 {
			conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
			_, err = conn.Write(sizeBuffer[:])
			errors.Check(err)
			continue
		}

		const MAX_SIZE = 1 << 30 // 1 Gb
		if size > MAX_SIZE {
			panic(fmt.Errorf("size too large, size=%d bytes", size))
		}
		if uint64(cap(buffer)) < size {
			buffer = make([]byte, size)
		}
		_, err = io.ReadFull(conn, buffer[:size])
		errors.Check(err)
		log.Infof("received image[%d bytes] from host", size)
		<-clipboard.Write(clipboard.FmtImage, buffer[:size])
		log.Infof("written image[%d bytes] to clipboard", size)
	}
}
