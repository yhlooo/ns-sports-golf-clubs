package main

import (
	"context"
	"fmt"
	"log"
	"machine"
	"strconv"
	"time"

	"github.com/yhlooo/ns-sports-golf-clubs/pkg/encoder"
	"github.com/yhlooo/ns-sports-golf-clubs/pkg/golfclubs"
	"github.com/yhlooo/ns-sports-golf-clubs/pkg/menu"
)

func main() {
	time.Sleep(2 * time.Second)

	clubs := golfclubs.New(
		machine.GPIO2,
		machine.GPIO3,
		machine.GPIO4,
	)
	if err := clubs.Configure(golfclubs.Config{}); err != nil {
		log.Fatalf("configure golf clubs error: %v", err)
	}

	machine.GPIO8.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	enc := &encoder.Encoder{
		APin:    machine.GPIO6,
		BPin:    machine.GPIO7,
		Reverse: false,
	}
	if err := enc.Configure(); err != nil {
		log.Fatalf("configure encoder error: %v", err)
	}

	m := &menu.Menu{}
	m.AddItems(
		(&menu.Item{Name: "1a"}).AddSubItems(
			menu.BackItem("Back"),
			(&menu.Item{Name: "full", Run: func(m *menu.Menu) {
				time.Sleep(10 * time.Second)
				fmt.Println("full")
			}}).AddSubItems(
				menu.BackItem("Back"),
			),
			(&menu.Item{Name: "2/3", Run: func(m *menu.Menu) {
				fmt.Println("2/3")
			}}).AddSubItems(
				menu.BackItem("Back"),
			),
			(&menu.Item{Name: "1/3", Run: func(m *menu.Menu) {
				fmt.Println("1/3")
			}}).AddSubItems(
				menu.BackItem("Back"),
			),
			(&menu.Item{Name: "custom", Run: func(m *menu.Menu) {
				fmt.Println("custom")
			}}).AddSubItems(
				menu.BackItem("Back"),
			),
		),
		(&menu.Item{Name: "2"}).AddSubItems(menu.BackItem("Back")),
		(&menu.Item{Name: "3"}).AddSubItems(menu.BackItem("Back")),
		(&menu.Item{Name: "4a"}).AddSubItems(menu.BackItem("Back")),
		(&menu.Item{Name: "4b"}).AddSubItems(menu.BackItem("Back")),
	)
	serialUI := &menu.Serial{Serial: machine.Serial}
	encoderUI := &menu.Encoder{
		Encoder:   enc,
		ButtonPin: machine.GPIO8,
	}
	m.AddOutputs(serialUI)
	m.AddInputs(
		serialUI,
		encoderUI,
	)

	m.HandleInputs(context.Background())

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
