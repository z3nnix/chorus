package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"

	"z3nnix/chorus/pkg/chorus"
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

func main() {
	startTime := time.Now()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Printf("\n%s\n", red("interrupt"))
		os.Exit(1)
	}()

	config, err := chorus.LoadConfig("chorus.build")
	if err != nil {
		exitWithError(err)
	}

	if err := config.Build(os.Args[1:]...); err != nil {
		exitWithError(err)
	}

	fmt.Printf("\n%s %s\n", green("done"), time.Since(startTime).Round(time.Millisecond))
}

func exitWithError(err error) {
	log.Fatal(red(err.Error()))
}
