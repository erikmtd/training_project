package main

import (
	"fmt"

	"github.com/training_project/handler"
	"github.com/training_project/handler/web"
)

func main() {
	fmt.Println("App Running")
	start(web.New())
	/*
		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)

		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigs
			fmt.Println()
			fmt.Println(sig)
			done <- true
		}()

		fmt.Println("awaiting signal")
		<-done
		db.New().Close()
		fmt.Println("exiting")*/
}
func start(handlers ...handler.Handler) {
	for _, h := range handlers {
		h.Start()
	}
}
