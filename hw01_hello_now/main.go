package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const ntpServer = "ntp3.stratum2.ru"

func main() {
	currentTime := time.Now()

	ntpTime, err := ntp.Time(ntpServer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("current time:", currentTime.Format("2006-01-02 15:04:05 -0700 MST"))
	fmt.Println("exact time:", ntpTime.Round(time.Second))
}
