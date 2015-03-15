# go-samsung-tv

[![godoc](http://img.shields.io/badge/godoc-Reference-blue.svg)](https://godoc.org/github.com/ninjasphere/go-samsung-tv)
[![MIT License](https://img.shields.io/badge/license-MIT-yellow.svg)](LICENSE)
[![Ninja Sphere](https://img.shields.io/badge/built%20by-ninja%20blocks-lightgrey.svg)](http://ninjablocks.com)

---


A simple golang package to control Samsung TVs over IP.

A list of possible commands is available from http://wiki.samygo.tv/index.php5/Key_codes

### Example

This example will turn the volume up whenever the TV becomes contactable.

```go
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
		ApplicationName: "Ninja Sphere         ", // XXX: Currently needs padding
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


```
