package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/psewda/pie/app"
)

var q chan os.Signal

func main() {
	fmt.Printf("*********************************************************************\n")
	fmt.Printf("------------------------- PIE SESSION STORE -------------------------\n")
	fmt.Printf("*********************************************************************\n")
	fmt.Printf("Version:     %s\n", app.Version)
	fmt.Printf("Golang:      %s\n", app.Golang)
	fmt.Printf("Git-Commit:  %s\n", app.GitCommit)
	fmt.Printf("Built:       %s\n", app.Built)
	fmt.Printf("OS/Arch:     %s\n", app.OsArch)
	fmt.Printf("*********************************************************************\n")
	fmt.Printf("\n")

	app := app.NewApp()
	app.Run(8800)

	done := make(chan bool)
	q = handleOsSignal(app, done)
	<-done
	log.Printf("process terminated gracefully, have a wonderful day")
}

func handleOsSignal(app app.App, done chan bool) chan os.Signal {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		sig := <-quit
		log.Printf("caught os signal, terminating process => [%s]", sig)
		app.Dispose()
		close(done)
	}()
	return quit
}
