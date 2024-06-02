package main

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

func main() {
	// Initialize periph.io
	if _, err := host.Init(); err != nil {
		fmt.Println("Failed to initialize periph.io:", err)
		return
	}

	// Define GPIO pins
	touchPin := rpi.P1_11 // GPIO17
	ledPin := rpi.P1_12   // GPIO18

	// Set touchPin as input
	if err := touchPin.In(gpio.PullUp, gpio.BothEdges); err != nil {
		fmt.Println("Failed to set touchPin as input:", err)
		return
	}

	// Set ledPin as output
	if err := ledPin.Out(gpio.Low); err != nil {
		fmt.Println("Failed to set ledPin as output:", err)
		return
	}

	ledState := false

	fmt.Println("Touch sensor is ready. Touch to toggle LED.")

	for {
		// Wait for a touch event
		touchPin.WaitForEdge(-1)

		// Debounce: wait a short period to confirm touch
		time.Sleep(300 * time.Millisecond)

		if touchPin.Read() == gpio.High {
			// Toggle LED state
			ledState = !ledState
			if ledState {
				ledPin.Out(gpio.High)
				fmt.Println("LED is now ON")
			} else {
				ledPin.Out(gpio.Low)
				fmt.Println("LED is now OFF")
			}
		}

		// Wait for touch to end before continuing
		for touchPin.Read() == gpio.High {
			time.Sleep(50 * time.Millisecond)
		}
	}
}
