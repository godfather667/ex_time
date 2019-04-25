package main

import (
	"fmt"
	"os"
	"time"
)

const top_limit = 10800 // Number of Seconds in 3 hours
const hourCount = 3600  // One Hour if tick duration = 1 second
const minCount = 60     // One Minute if tick duration = 1 second

const msecond = 1000000 // Tick Duration multiplier for Nanoseconds (1 millisecond)
const tsec = time.Second

// Set timeDuration -- Time between Ticks
var timeDuration = tsec

// Default Time Marker variables
var bong = "Bong"
var tock = "Tock"
var tick = "tick"

func tickClock(sec chan<- string, stop chan<- bool) {
	s := 0
	ticker := time.NewTicker(timeDuration)
	<-ticker.C // Sync on next Second Tick
	defer ticker.Stop()
	for {
		<-ticker.C
		s++
		fmt.Print(s, ": ")
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
	//	fmt.Println(bong)
	stop <- true
}

func clock() {
	sec := make(chan string)
	stop := make(chan bool)
	go tickClock(sec, stop)
	for {
		select {
		case msg := <-sec:
			fmt.Println(msg)
		case <-stop:
			break
		}
	}
}

func main() {
	clock()    // Call Clock Function
	os.Exit(0) // Exit on return
}
