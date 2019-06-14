package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/training_project/handler"
	"github.com/training_project/handler/mq"
)

func main() {
	fmt.Println("App Running")
	start(mq.New())

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case <-term:
		log.Println("Application terminated")
	}
}
func start(handlers ...handler.Handler) {
	for _, h := range handlers {
		h.Start()
	}
}
