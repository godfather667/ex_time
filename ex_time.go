// ex_time.go - An Demo Program with an example of golang testing.
//
//   This progam  counts seconds for three hours. Each Second it prints
//   the tag "Tick" unless it is a minute period and then prints the tag "Tock".
//   Finally, at hour periods it prints the tag "Bong".
//
//   An additional freature it allows changing the "tag" values by setting your
//   browser to localhost:3000 and using the tags as endpoint for changing
//	 Example:  localhost:3000/tick/new_tag_word.
//
//   The tag words may be any case, however the new_tag_word is displayed as is.
//   The tags used for endpoints do not change, the tag "Tick" is always used
//   to modify the Tick tag.
//
//   This program is almost fully tested with "ex_time_test.go".
//
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

// Output Message Silencing for Testing
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
// This was the orginal main function to improve testability.
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

// msgHandler -- This function interacts with the browser to update
//    tag words.  It also displays on the browser a change notice.
//    It also display an error and help message if tag is incorrect.
func msgHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Path
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	fields := strings.FieldsFunc(msg, f)
	if len(fields) >= 2 {
		endp := strings.ToLower(fields[0]) // Converts tag to lowercase
		val := fields[1]                   // val is the new tag word
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

			// The follwing message is displayed when no valid tag is isolated.
		default:
			fmt.Fprintf(w, "Invalid Endpoint: Must be \"tick\" or \"tock\" or \"bong\"\n\n")
			fmt.Fprintf(w, "For example: localhost:3000/tick/new_word_here!")
		}
		// The follwing message is displayed when less than two fields isolated.
	} else {
		fmt.Fprintf(w, "Invalid Endpoint: Must be \"tick\" or \"tock\" or \"bong\"\n\n")
		fmt.Fprintf(w, "For example: localhost:3000/tick/new_word_here!")
	}
}

// HealthCheckHandler -- Is a simple check of HTTP Processing.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"alive": true}`)
}

// Execcutive Function:
// 	  Executes the "clock" Function and then sets up the HTTP message handlers.
func main() {
	go clock() // Call Clock Function
	// HTTP Handers, Processing Loop
	http.HandleFunc("/", msgHandler)
	http.HandleFunc("/health-check", HealthCheckHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
