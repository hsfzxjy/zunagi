package main

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/hsfzxjy/zunagi/internal/errors"
	"github.com/hsfzxjy/zunagi/internal/log"
	"golang.design/x/clipboard"
)

func (s *Server) handleConnection(pctx context.Context, conn net.Conn) {
	ctx, cancel := context.WithCancel(pctx)

	var state string
	const (
		HEARTBEAT = "heartbeat"
		SENDDATA  = "writing to guest"
	)

	defer func() {
		cause := recover().(error)
		log.Errorf("error on %s: %s", state, cause)
		conn.Close()
		cancel()
		log.Infof("guest connection %s closed", conn.LocalAddr())
	}()

	imageCh := clipboard.Watch(ctx, clipboard.FmtImage)

	var hbTimer *time.Timer
	resetTimer := func() {
		if hbTimer == nil {
			hbTimer = time.NewTimer(time.Second)
			return
		}
		hbTimer.Reset(time.Second)

	}

	for {
		resetTimer()
		select {
		case <-hbTimer.C:
			state = HEARTBEAT
			heartBeat(conn)
		case image := <-imageCh:

			state = SENDDATA
			writeImageToGuest(conn, image)
			if !hbTimer.Stop() {
				<-hbTimer.C
			}
		}
	}
}

func heartBeat(conn net.Conn) {
	var zeroBuffer [8]byte
	conn.SetDeadline(time.Now().Add(time.Second))
	_, err := conn.Write(zeroBuffer[:])
	errors.Check(err)
	_, err = io.ReadFull(conn, zeroBuffer[:])
	errors.Check(err)
}

func writeImageToGuest(conn net.Conn, image []byte) {
	var sizeBuffer [8]byte
	log.Infof("captured clipboard image[%d bytes]", len(image))
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	binary.BigEndian.PutUint64(sizeBuffer[:], uint64(len(image)))
	_, err := conn.Write(sizeBuffer[:])
	errors.Check(err)
	_, err = conn.Write(image)
	errors.Check(err)
	log.Infof("written image[%d bytes] to guest", len(image))
}
