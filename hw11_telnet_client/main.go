package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"
)

var (
	timeout        time.Duration
	defaultTimeout = 10 * time.Second
)

func main() {
	// loggers for easy writing to stderr
	serviceLog := log.New(os.Stderr, "\n", 0)
	errorLog := log.New(os.Stderr, "\n", log.Llongfile)

	flag.DurationVar(
		&timeout,
		"timeout",
		defaultTimeout,
		"specify timeout for network operations")
	flag.Parse()

	if len(flag.Args()) < 2 {
		errorLog.Fatalf("specify host and port: got %+v", flag.Args())
	}

	host, port := flag.Arg(0), flag.Arg(1)

	// signal channel to catch SIGINT
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	// notification channel for exit
	stopCh := make(chan struct{})

	tc := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	defer tc.Close()

	if err := tc.Connect(); err != nil {
		errorLog.Fatalf("unable to connect to %s:%s: %s", host, port, err.Error())
	}

	go func() {
		err := tc.Send()
		if err != nil {
			errorLog.Println(err)
		} else {
			serviceLog.Println("...EOF")
		}
		stopCh <- struct{}{}
	}()

	go func() {
		err := tc.Receive()
		if err != nil {
			errorLog.Println(err)
		} else {
			serviceLog.Println("...Connection was closed by peer")
		}
		stopCh <- struct{}{}
	}()

	select {
	case <-stopCh:
	case <-sigCh:
	}

	serviceLog.Println("See ya!")
}
