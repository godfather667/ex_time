package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestClock(t *testing.T) {

	//	timeDuration = time.Second    // Tick time Period
	timeDuration = msecond * time.Nanosecond // Test Tick Time Period
	result := make([]string, top_limit)      // Clock result

	sec := make(chan string)
	stop := make(chan bool)
	go tickClock(sec, stop)
	for {
		select {
		case msg := <-sec:
			fmt.Println(msg)
		case <-stop:
			os.Exit(0)
		}
	}
	l := len(result)
	if l < top_limit {
		t.Errorf("Number of ticks incorrect, got: %d, want: %d", l, top_limit)
	}
}
