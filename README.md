**ex_time**

This a demo of **Golang Testing**.

The program is in **ex_time.go**

This program sets up a tick clock and prints message depending on the following:

* tickClock  -- Tick Duration = 1 Second

* Prints "Tock" for each 60 second period (Minutes)
* Prints "Bong" for each 3,600 seconds (hours)
* Prints "Tick" for each second that does not interfer with the other messages.

The program allows changing tag messages via browser at **localhost:3000**

The allowable tags are **tick, tock, or bong**. 

	**To change tags use the tag as an endpoint followed by the new tag value:
	
	localhost:3000/tick/new_tag_word<cr>**

If an invalid end point is entered, the following message is displayed:

         **Invalid Endpoint: Must be "tick" or "tock" or "bong"
       
       For Example: localhost:3000/tick/new_word_here!**

The testing code is in **test_time_test.go**

It provides 87% Coverage

It has two test functions:  TestTickClock and TestClock
