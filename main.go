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
		if button.Read() == rpio.Low {
			// fmt.Println("Button pressed!")
			// if time.Since(lastPressTime) < debouncePeriod {
			// 	// Ignore button press if within debounce period
			// 	continue
			// }
			// lastPressTime = time.Now()

			// Check if the Raspberry Pi is already running
			// if isRaspberryPiRunning() {
			// 	// If already running, initiate shutdown
			// 	fmt.Println("Initiating shutdown...")
			// 	shutdown()
			// } else {
			// 	// If not running, power on the Raspberry Pi
			// 	fmt.Println("Raspberry Pi is off. Powering on...")
			// 	powerOn()
			// }

			fmt.Println(rpio.High, "HIGH")
			fmt.Println(rpio.Low, "LOW")

			// fmt.Println("Raspberry Pi is off. Powering on...")
			powerOn()
			time.Sleep(debouncePeriod)
		}
		time.Sleep(100 * time.Millisecond)
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
	cmd := exec.Command("sudo", "rtcwake", "-m", "on", "-s", "0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to trigger wake-up:", err)
	}
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
