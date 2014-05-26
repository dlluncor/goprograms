package spoj

import (
	"fmt"
	"time"
)

func print(val string) {
	fmt.Println(val)
}

type Hi struct {
	val string
}

func (h *Hi) print() string {
	return fmt.Sprintf("%v\n", h.val)
}

func Akiva() string {
	hi := &Hi{"yo"}
	return hi.print()
}

// Learning about how channels work in go. Essentially put each
// function you want to execute in parallel inside of a go func()
// call it, put a value on a channel, and then pull from that channel
// similar to tp.WaitAll() in Python.
func Concurrency() {
	// Allocate a channel where 4 values can be stored on it before
	// another function gets blocked.
	c := make(chan error, 4)
	b := time.Now()
	for i, _ := range []string{"", "", "", ""} {
		val := fmt.Sprintf("hi %d", i)
		index := i
		// This is a function which runs in its own non-blocking thread.
		go func() {
			time.Sleep(4 * time.Second)
			print(val)
			// c <- is how we communicate values to a channel.
			if index != 2 {
				c <- nil
			} else {
				c <- fmt.Errorf("%d ff'ed up", index)
			}
		}()
	}
	fmt.Println("Outside chan.")
	// We can get all the values from a channel, checking to see
	// if there were any errors.
	for _ = range []string{"", "", "", ""} {
		x, ok := <-c
		a := time.Now()
		fmt.Printf("This is done. %v %v %v\n", x, ok, a.Sub(b))
	}
}
