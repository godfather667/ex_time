// ex_time_test.go -- Test Functions for ex_time.go.
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Notice the three counts add up to "top_limit"
const tickCount = 10620
const tockCount = 177
const bongCount = 3

// Test Reponses for testMsgHandler Function.
var testCode = []string{"/tick/test", "Tick Changed to:  test\n",
	"/tock/test", "Tock Changed to:  test\n",
	"/bong/test", "Bong Changed to:  test\n",
	"/Fail/Test", "Invalid Endpoint: Must be \"tick\" or \"tock\" or \"bong\"\n\nFor example: localhost:3000/tick/new_word_here!",
	"/Short", "Invalid Endpoint: Must be \"tick\" or \"tock\" or \"bong\"\n\nFor example: localhost:3000/tick/new_word_here!"}

// Global Initialization Values.
func init() {
	// Set Initial State for Messages
	bong = "Bong"
	tock = "Tock"
	tick = "Tick"
	// Set standard Tick Duration (1 Second)
	timeDuration = msec
	// Set silent running for testing (So only the test generated messages are present).
	// The outputs of ex_time.go are still available to the testing functions, but are
	// not displayed on the console. This makes for a cleaner looking test.
	silent = true
}

// TestTickClock -- Test the "tickClock.
func TestTickClock(t *testing.T) {

	// Set Test Tick Time Period of "1 msec" for testing (Improves run time significantly).
	//   NOTE: The test period could less thab 1 msec - but this might lead to latency issues.
	timeDuration = msec                    // 1 Millisecond Period.
	result := make([]string, 0, top_limit) // Clock result buffer

	sec := make(chan string) // Channel for reporting Ticks form tickClock.
	stop := make(chan bool)  // Channel for reporting end of Ticks
	go tickClock(sec, stop)  // Launch tickClock Function

	for {
		select {
		// Handle Message Processing
		case msg := <-sec: // Acquire Tick Message Value.
			if !silent {
				fmt.Println(msg) // Suppress message to console -
			}
			result = append(result, msg) //append msg to result buffer

		// Handle Stop Ticks Channel and compute and evaluate results.
		case <-stop:
			l := len(result)
			if !silent {
				fmt.Println("len = ", l)
			}
			if l != top_limit {
				t.Errorf("Number of ticks incorrect, got: %d, want: %d", l, top_limit)
			}
			// Compute Freq of Messages to verify correct number of each message.
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
			// Verify that records number of each tag message are correct.
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

// TestClock -- Test the "clock" Function
func TestClock(t *testing.T) {

	//	timeDuration = 1 millisecond   // Tick time Period
	timeDuration = msec
	go clock() // Launch Clock Function
}

// TestHealthCheckHandler - Checks Health of HTTP Connection
// NOTE - In Future Refactor testMsgHandler to perform this test also.
func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestMsgHandler -- Test Endpoint Processing

func TestMsgHandler(t *testing.T) {
	for i := 0; i < 10; i += 2 {
		// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", testCode[i], nil)
		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(msgHandler)

		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		if rr.Body.String() != testCode[i+1] {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), testCode[i+1])
		}
	}
}
