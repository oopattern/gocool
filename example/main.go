package main

import (
	"fmt"
	"time"
	"errors"
)

var errStop = errors.New("just quit...")

func main() {
	fmt.Println("process start")

	var (
		errc = make(chan error, 1)
		resc = make(chan int, 1)
		interruptc = make(chan struct{}, 1)
	)

	go func() {
		time.Sleep(time.Second)
		fmt.Println("catch stop signal")
		errc <- errStop
	}()

	select {
	case err := <-errc:
		fmt.Printf("catch err[%v]", err)
	case res := <-resc:
		fmt.Printf("catch resc signal[%v]", res)
	case c := <-interruptc:
		fmt.Printf("catch interrupt signal[%v]", c)
	}

	fmt.Println("process finish")
}

