// ex_time.go - Test Program for an example of golang testing.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"
)

const top_limit = 10800 // Number of Seconds in 3 hours
const hourCount = 3600  // One Hour if tick duration = 1 second
const minCount = 60     // One Minute if tick duration = 1 second

const msec = 1000000 * time.Nanosecond // Tick Duration multiplier (1 millisecond)
const tsec = time.Second

// Set timeDuration -- Time between Ticks
var timeDuration = tsec
var silent = false

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
			if !silent {
				fmt.Println(msg)
			}
		case <-stop:
			return
		}
	}
}

func msgHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Path
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	fields := strings.FieldsFunc(msg, f)
	if len(fields) >= 2 {
		endp := strings.ToLower(fields[0])
		val := fields[1]
		switch endp {
		case "tick":
			tick = val
			fmt.Fprintf(w, "Tick Changed to:  %s\n", val)
		case "tock":
			tock = val
			fmt.Fprintf(w, "Tock Changed to:  %s\n", val)
		case "bong":
			bong = val
			fmt.Fprintf(w, "Bong Changed to:  %s\n", val)
		default:
			fmt.Fprintf(w, "Invalid Endpoint: Must be \"tick\" or \"tock\" or \"bong\"\n\n")
			fmt.Fprintf(w, "For example: localhost:3000/tick/new_word_here!")
		}
	} else {
		fmt.Fprintf(w, "Invalid Endpoint: Must be \"tick\" or \"tock\" or \"bong\"\n\n")
		fmt.Fprintf(w, "For example: localhost:3000/tick/new_word_here!")
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

// Execcutive Function:
// 	  executes "clock" Function and then exits
func main() {
	go clock() // Call Clock Function
	http.HandleFunc("/", msgHandler)
	http.HandleFunc("/health-check", HealthCheckHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
	//	os.Exit(0) // Exit on return
}
