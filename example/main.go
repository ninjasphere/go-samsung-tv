package main

import (
	"log"
	"time"

	"github.com/ninjasphere/go-samsung-tv"
)

func main() {
	samsung.EnableLogging = true

	tv := samsung.TV{
		Host:            "192.168.0.21",
		ApplicationID:   "go-samsung-tv",
		ApplicationName: "Ninja Sphere         ", // XXX: Need extra spaces or bits of "samsung remote" gets added on the end??
	}

	tv.OnPowerChange(time.Second*5, func(online bool) {
		if online {
			log.Println("TV came online!")

			if err := tv.SendCommand("KEY_VOLUP"); err != nil {
				log.Printf("Failed to send command. Error: %s", err)
			}

		} else {
			log.Println("TV went offline!")
		}
	})

}
