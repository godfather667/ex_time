package main

import (
	"fmt"
	"testing"
)

func TestClock(t *testing.T) {

	//	timeDuration = time.Second    // Tick time Period
	timeDuration = msec                 // Test Tick Time Period
	result := make([]string, top_limit) // Clock result

	sec := make(chan string)
	stop := make(chan bool)
	go tickClock(sec, stop)
	for {
		select {
		case msg := <-sec:
			fmt.Println(msg)
		case <-stop:
			l := len(result)
			fmt.Println("len = ", l)
			if l == top_limit {
				t.Errorf("Number of ticks incorrect, got: %d, want: %d", l, top_limit)
			}
			fmt.Println("TestClock break")
			return
		}
	}
}
