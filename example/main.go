package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var errStop = errors.New("just quit...")

func chanTest() {
	var (
		errc = make(chan error, 1)
		resc = make(chan int, 1)
		interruptc = make(chan struct{}, 1)
	)

	go func() {
		errc <- errStop
		fmt.Printf("1. errc: %d\n", len(errc))
		time.Sleep(2*time.Second)
	}()

	// time.Sleep(1*time.Second)

	select {
	case err := <-errc:
		time.Sleep(1*time.Second)
		fmt.Printf("x catch err[%v]\n", err.Error())
	case res := <-resc:
		fmt.Printf("x catch resc signal[%v]\n", res)
	case c := <-interruptc:
		fmt.Printf("x catch interrupt signal[%v]\n", c)
		//default:
		//	time.Sleep(5*time.Second)
		//	fmt.Printf("x just do nothing...%d\n", len(errc))
	}
}

func contextTest() {
	c1 := context.Background()
	c2 := context.Background()
	c3 := context.Background()
	c4 := context.TODO()
	c5 := context.TODO()
	fmt.Printf("%p\n", c1)
	fmt.Printf("%p\n", c2)
	fmt.Printf("%p\n", c3)
	fmt.Printf("%p\n", c4)
	fmt.Printf("%p\n", c5)
}

func main() {
	contextTest()
}

