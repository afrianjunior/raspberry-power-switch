package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	rpio "github.com/stianeikeland/go-rpio/v4"
)

const (
	gpioPin        = 17
	shutdownTime   = 3 // seconds
	longPressTime  = 3 // seconds
	debouncePeriod = 50 * time.Millisecond
)

var lastPressTime time.Time

func main() {
	// Open and map memory to access GPIO, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println("Failed to open GPIO:", err)
		return
	}
	defer rpio.Close()

	// Set GPIO pin mode to input
	button := rpio.Pin(gpioPin)
	button.Input()

	fmt.Println("Waiting for button press...")

	for {
		// Wait for a rising edge (button press)
		if button.Read() == rpio.High {
			fmt.Println("Button pressed!")
			if time.Since(lastPressTime) < debouncePeriod {
				// Ignore button press if within debounce period
				continue
			}
			lastPressTime = time.Now()

			// Check if the Raspberry Pi is already running
			if isRaspberryPiRunning() {
				// If already running, initiate shutdown
				fmt.Println("Initiating shutdown...")
				shutdown()
			} else {
				// If not running, power on the Raspberry Pi
				fmt.Println("Raspberry Pi is off. Powering on...")
				powerOn()
			}
		}
		time.Sleep(debouncePeriod)
	}
}

func isRaspberryPiRunning() bool {
	cmd := exec.Command("hostnamectl", "status")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to get hostname status:", err)
		return false
	}
	return !bytes.Contains(output, []byte("off"))
}

func powerOn() {
	pin := rpio.Pin(gpioPin)
	pin.Output()
	pin.Low()

	// Wait for a short duration to ensure the wake-up signal is detected
	time.Sleep(100 * time.Millisecond)

	// Set the pin back to input mode to avoid holding it low indefinitely
	pin.Input()

	fmt.Println("Wake-up signal sent.")
}

func shutdown() {
	// Execute the shutdown command
	cmd := exec.Command("sudo", "shutdown", "-h", "now")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to initiate shutdown:", err)
	}
}
