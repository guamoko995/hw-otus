package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func main() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "server connection timeout")
	flag.Parse()

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))
	tc := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)

	if err := tc.Connect(); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	go func() {
		if err := tc.Send(); err != nil {
			if errors.Is(err, net.ErrClosed) {
				fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
			}
		} else {
			fmt.Fprintln(os.Stderr, "...EOF")
		}
		tc.Close()
		cancel()
	}()

	go func() {
		tc.Receive()
		tc.Close()
	}()

	<-ctx.Done()
}
