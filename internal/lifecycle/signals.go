package lifecycle

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Setup() {
	ch := make(chan os.Signal, 8)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-ch
		signal.Stop(ch)
		log.Printf("signal %v received, exiting\n", s)
		Shutdown()
	}()
}
