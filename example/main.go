package main

import (
	"log"
	"time"

	"github.com/ninjasphere/go-samsung-tv"
)

func main() {
	samsung.EnableLogging = true

	tv := samsung.TV{
		Host:            "192.168.1.21",
		ApplicationID:   "go-samsung-tv",
		ApplicationName: "Ninja Sphere         ", // XXX: Need extra spaces or bits of "samsung remote" gets added on the end??
	}

	// Once-off check if tv is online (timeout after 2 seconds)
	if tv.Online(time.Second * 2) {
		log.Println("TV is online!")
	} else {
		log.Println("TV is offline!")
	}

	// Continuous updates as TV goes online and offline
	for online := range tv.OnlineState(time.Second * 5) {

		if online {
			log.Println("TV came online!")

			// Turn the volume up when it comes online
			if err := tv.SendCommand("KEY_VOLUP"); err != nil {
				log.Printf("Failed to send command. Error: %s", err)
			}
		} else {
			log.Println("TV went offline!")
		}

	}

}
