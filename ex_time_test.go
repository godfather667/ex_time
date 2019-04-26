package main

import (
	"fmt"
	"testing"
)

// Notice the three counts add up to "top_limit"
const tickCount = 10620
const tockCount = 177
const bongCount = 3

func init() {
	// Set Initial State for Messages
	bong = "Bong"
	tock = "Tock"
	tick = "Tick"
	// Set standard Tick Duration (1 Second)
	timeDuration = tsec
}

func TestClock(t *testing.T) {

	//	timeDuration = time.Second    // Tick time Period
	timeDuration = msec                // Test Tick Time Period
	result := make([]string, 0, 10800) // Clock result

	sec := make(chan string)
	stop := make(chan bool)
	go tickClock(sec, stop)
	for {
		select {
		case msg := <-sec:
			fmt.Println(msg)
			result = append(result, msg)
		case <-stop:
			l := len(result)
			fmt.Println("len = ", l)
			if l != top_limit {
				t.Errorf("Number of ticks incorrect, got: %d, want: %d", l, top_limit)
			}
			// Compute Freq of Messages
			recTick := 0
			recTock := 0
			recBong := 0
			for _, v := range result {
				if v == "Tick" {
					recTick++
				}
				if v == "Tock" {
					recTock++
				}
				if v == "Bong" {
					recBong++
				}
			}
			if recTick != tickCount {
				t.Errorf("Number of ticks incorrect, got: %d, want: %d", recTick, tickCount)
			}
			if recTock != tockCount {
				t.Errorf("Number of ticks incorrect, got: %d, want: %d", recTock, tockCount)
			}
			if recBong != bongCount {
				t.Errorf("Number of ticks incorrect, got: %d, want: %d", recBong, bongCount)
			}
			return
		}
	}
}
