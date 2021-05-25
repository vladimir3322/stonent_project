package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WaitSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals
	fmt.Println("Got signal for exiting", sig)
}
