package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	currentTime := time.Now()
	fmt.Printf("current time: %s\n", currentTime.Format("2006-01-02 15:04:05 -0700 MST"))
	NTPtime, err := ntp.Time("ntp3.stratum2.ru")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("exact time: %s\n", NTPtime.Round(time.Second))
}
