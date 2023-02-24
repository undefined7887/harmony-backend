package main

import (
	"context"
	"fmt"
	stdlog "log"
	"time"
)

const (
	StartTimeout = time.Second * 15
	StopTimeout  = time.Second * 2
)

func main() {
	app := NewApp()

	startCtx, cancel := context.WithTimeout(context.Background(), StartTimeout)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		stdlog.Fatalln(fmt.Errorf("failed to start: %v", err))
	}

	signal := <-app.Done()
	fmt.Println()
	stdlog.Printf("received %s signal", signal.String())

	stopCtx, cancel := context.WithTimeout(context.Background(), StopTimeout)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		stdlog.Fatalln(fmt.Errorf("failed to stop: %v", err))
	}

	fmt.Println("bye")
}
