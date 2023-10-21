package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func cleanup() {

}
func main() {
	fmt.Println("Starting Aggregator")

	// cleanup handler to quit when told to do so
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	go func() {
		<-ch
		cleanup()
		os.Exit(1)
	}()
}
