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
		Host:            "192.168.0.21",
		ApplicationID:   "go-samsung-tv",
		ApplicationName: "Ninja Sphere         ", // XXX: Padding is required...
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

```
