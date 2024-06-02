package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println("Failed to open gpio:", err)
		return
	}
	defer rpio.Close()

	// Pin definitions
	touchPin := rpio.Pin(17) // GPIO17
	ledPin := rpio.Pin(18)   // GPIO18

	// Set touchPin as input
	touchPin.Input()

	// Set ledPin as output
	ledPin.Output()

	ledState := rpio.Low // Initialize LED state to Low (off)

	fmt.Println("Touch sensor is ready. Touch to toggle LED.")

	for {
		// Wait for a touch event
		if touchPin.Read() == rpio.High {
			// Debounce: wait a short period to confirm touch
			time.Sleep(300 * time.Millisecond)

			if touchPin.Read() == rpio.High {
				// Toggle LED state
				if ledState == rpio.Low {
					ledState = rpio.High
					ledPin.High()
					fmt.Println("LED is now ON")
				} else {
					ledState = rpio.Low
					ledPin.Low()
					fmt.Println("LED is now OFF")
				}
			}

			// Wait for touch to end before continuing
			for touchPin.Read() == rpio.High {
				time.Sleep(50 * time.Millisecond)
			}
		}
	}
}
