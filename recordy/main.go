package main

import (
	"fmt"
	"os"
	"time"

	"recordy/internal/model"

	"golang.org/x/term"
)

func main() {
	first := true
	fmt.Println("Hello, World!")

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Error making raw mode:", err)
		return
	}

	events := []model.Event{}
	buf := make([]byte, 1)

	startTime := time.Now()
	lastEventTime := startTime
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		if n == 0 {
			continue
		}
		b := buf[0]

		now := time.Now()
		var dt int64
		if first {
			first = false
			dt = 0
		} else {
			dt = now.Sub(lastEventTime).Milliseconds()
		}

		lastEventTime = now

		str := mapByte(b)
		if str == "CTRL_C" {
			break
		}

		events = append(events, model.Event{
			T:    dt,
			Type: "KEY",
			Key:  str,
		})

		lastEventTime = time.Now()
	}

	term.Restore(fd, oldState)
	fmt.Println("\nRecorded events:")
	for _, e := range events {
		fmt.Printf("%+v\n", e)
	}
}

func mapByte(b byte) string {
	switch b {
	case 3:
		return "CTRL_C"
	case 13:
		return "ENTER"
	case 27:
		return "ESC"
	default:
		if b >= 32 && b <= 126 {
			return string(b)
		}
		return "UNKNOWN"
	}
}
