package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Printf("%d\n", runtime.Getgid())
	ch := make(chan struct{})
	go func() {
		fmt.Printf("%d\n", runtime.Getgid())
		ch <- struct{}{}
	}()
	go func() {
		fmt.Printf("%d\n", runtime.Getgid())
		ch <- struct{}{}
	}()
	go func() {
		fmt.Printf("%d\n", runtime.Getgid())
		ch <- struct{}{}
	}()
	<-ch
	<-ch
	<-ch
	return nil
}
