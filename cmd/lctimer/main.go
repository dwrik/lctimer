package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Time limits
const (
	etime = iota*10 + 10 // 10 mins
	mtime                // 20 mins
	htime                // 30 mins
)

// Settings type
type settings struct {
	time int
	diff string
}

// Current settings
var setting settings

// Set difficulty and time limit (mins) based on difficulty flag
func setDiffTime(easy, medium, hard bool) {
	switch {
	case easy:
		setting = settings{etime, "easy"}
	case medium:
		setting = settings{mtime, "medium"}
	case hard:
		setting = settings{htime, "hard"}
	}
}

// Wait for user intervention
func userInput(ch chan bool) {
	var temp string
	fmt.Scanln(&temp)
	ch <- true
}

// Send notification
func notify() {
	// notification properties
	sound := "Frog"
	body := "Time's up!"
	title := "LeetCode Timer"

	// build applescript
	script := fmt.Sprintf(
		"display notification \"%s\" with title \"%s\" sound name \"%s\"",
		body, title, sound,
	)

	// run applescript
	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send notification: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	// Define difficulty flags
	easy := flag.Bool("easy", false, "set timer for an easy problem (10 mins)")
	medium := flag.Bool("medium", false, "set timer for a medium problem (20 mins)")
	hard := flag.Bool("hard", false, "set timer for a hard problem (30 mins)")

	// Define custom usage
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s sets a timer based on the difficulty\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "of the problem and notifies when the timer ends\n")
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s <difficulty>\n", os.Args[0])
		flag.CommandLine.PrintDefaults()
		os.Exit(2)
	}

	// Parse flags
	flag.Parse()

	// Set difficulty and time limits
	setDiffTime(*easy, *medium, *hard)

	// Print usage and exit if
	// difficulty not provided
	if setting.diff == "" {
		flag.Usage()
		os.Exit(5)
	}

	// Create timer
	timer := time.NewTimer(time.Duration(setting.time) * time.Minute)
	fmt.Fprintf(os.Stdout, "[+] timer started\n")

	// Get user input
	userC := make(chan bool)
	go userInput(userC)
	fmt.Fprint(os.Stdout, "[i] press enter to stop the timer...")

	// Handle Time up or
	// User intervention
	select {
	case <-timer.C:
		notify()
		fmt.Fprintln(os.Stdout, "\n[!] time's up")
	case <-userC:
		if !timer.Stop() {
			<-timer.C
		}
		fmt.Fprintln(os.Stdout, "[x] timer stopped")
	}
}
