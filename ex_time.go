// ex_time.go - Test Program for an example of golang testing.
package main

import (
	"fmt"
	"os"
	"time"
)

const top_limit = 10800 // Number of Seconds in 3 hours
const hourCount = 3600  // One Hour if tick duration = 1 second
const minCount = 60     // One Minute if tick duration = 1 second

const msec = 1000000 * time.Nanosecond // Tick Duration multiplier (1 millisecond)
const tsec = time.Second

// Set timeDuration -- Time between Ticks
var timeDuration = msec

// Default Time Marker variables
var bong = "Bong"
var tock = "Tock"
var tick = "tick"

// tickClock  -- Tick Duration = 1 Second
// Prints "Tock" for each 60 second period
// Prints "Bong" for each 3,600 seconds
// Prints "Tick" for each second that does not interfer with the other messages.
func tickClock(sec chan<- string, stop chan<- bool) {
	s := 0
	ticker := time.NewTicker(timeDuration)
	<-ticker.C // Sync on next Second Tick
	defer ticker.Stop()
	for {
		<-ticker.C
		s++
		if s%hourCount == 0 {
			sec <- bong
		} else if s%minCount == 0 {
			sec <- tock
		} else {
			sec <- (tick)
		}
		if s >= top_limit { // 3 hour count
			break
		}
	}
	stop <- true
}

// Clock manages the tickClock execution displaying message from tickClock function
func clock() {
	sec := make(chan string)
	stop := make(chan bool)
	go tickClock(sec, stop)
	for {
		select {
		case msg := <-sec:
			fmt.Println(msg)
		case <-stop:
			return
		}
	}
}

// Execcutive Function:
// 	  executes "clock" Function and then exits
func main() {
	clock()    // Call Clock Function
	os.Exit(0) // Exit on return
}
