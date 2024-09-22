package main

import (
	"fmt"
	"log"
	"machine"
	"strconv"
	"time"

	"github.com/yhlooo/ns-sports-golf-clubs/pkg/golfclubs"
)

func main() {
	time.Sleep(2 * time.Second)
	log.Print("ggg")

	clubs := golfclubs.New(
		machine.GPIO2,
		machine.GPIO3,
		machine.GPIO4,
	)
	if err := clubs.Configure(golfclubs.Config{}); err != nil {
		log.Fatalf("configure golf clubs error: %v", err)
	}

	reversed := false
	for {
		line := readLine(machine.Serial)

		switch line {
		case "reverse":
			reversed = !reversed
			clubs.SetReverse(reversed)
		default:
			speed, err := strconv.ParseInt(line, 10, 8)
			if err != nil {
				log.Printf("ERROR parse speed %q error: %v", line, err)
				continue
			}
			if speed <= 0 {
				log.Printf("swing speed is 0, skipped")
				continue
			}
			if speed >= 100 {
				speed = 100
			}
			log.Printf("swing at %d speed", speed)
			clubs.Swing(uint8(speed))
			log.Printf("swing done")
		}
	}
}

func readLine(s machine.Serialer) string {
	var line []byte
	for {
		c, err := s.ReadByte()
		if err != nil {
			continue
		}
		if c == '\n' || c == '\r' {
			fmt.Print("\n")
			return string(line)
		}
		fmt.Print(string(c))
		line = append(line, c)
	}
}
