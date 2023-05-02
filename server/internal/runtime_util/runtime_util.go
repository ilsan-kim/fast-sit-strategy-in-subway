package runtime_util

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

var GracefulShubdownJob chan struct{}

func RegisterSignal(stopFunc func()) (done chan struct{}) {
	done = make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)

		s := <-c

	B:
		for {
			switch s {
			case syscall.SIGHUP:
				fmt.Println("I got a SIGHUP signal")
			case syscall.SIGINT:
				fmt.Println("I got a SIGINT signal")
				break B
			case syscall.SIGTERM:
				fmt.Println("I got a SIGTERM signal")
				break B
			case syscall.SIGQUIT:
				fmt.Println("I got a SIGQUIT signal")
			default:
				fmt.Println("I got a signal ", s)
			}
			time.Sleep(time.Second)
		}

		stopFunc()
		done <- struct{}{}
	}()
	return done
}

func RunFuncSafely(f func()) {
	GracefulShubdownJob <- struct{}{}
	defer func() {
		r := recover()
		if r != nil {
			log.Println(string(debug.Stack()))
		}
		<-GracefulShubdownJob
	}()
	f()
}
