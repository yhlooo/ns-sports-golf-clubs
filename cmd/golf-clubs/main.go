package main

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"machine"
	"strconv"
	"time"

	"tinygo.org/x/drivers/sh1106"
	"tinygo.org/x/tinyfont/proggy"

	"github.com/yhlooo/ns-sports-golf-clubs/pkg/encoder"
	"github.com/yhlooo/ns-sports-golf-clubs/pkg/golfclubs"
	"github.com/yhlooo/ns-sports-golf-clubs/pkg/menu"
)

func main() {
	time.Sleep(2 * time.Second)

	// 初始化高尔夫球杆
	clubs := golfclubs.New(
		machine.GPIO2,
		machine.GPIO3,
		machine.GPIO4,
	)
	if err := clubs.Configure(golfclubs.Config{}); err != nil {
		log.Fatalf("configure golf clubs error: %v", err)
	}

	// 初始化显示器
	i2c := machine.I2C1
	if err := i2c.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
		SDA:       machine.GPIO10,
		SCL:       machine.GPIO11,
	}); err != nil {
		log.Fatalf("configure i2c error: %v", err)
	}
	display := sh1106.NewI2C(i2c)
	display.Configure(sh1106.Config{Width: 128, Height: 32})
	display.ClearDisplay()

	// 初始化编码器
	machine.GPIO8.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	enc := &encoder.Encoder{
		APin: machine.GPIO6,
		BPin: machine.GPIO7,
	}
	if err := enc.Configure(); err != nil {
		log.Fatalf("configure encoder error: %v", err)
	}

	// 初始化菜单
	settingsNode := &menu.BaseNode{NodeName: "Settings"}
	settingsNode.AddChildren(
		menu.NewBackNode("Back"),
		menu.NewBoolValueNode("Reverse", false, true, func(reverse bool) {
			clubs.SetReverse(reverse)
		}),
	)
	root := &menu.BaseNode{NodeName: "Root"}
	root.AddChildren(
		&menu.ValueNode{
			BaseNode: menu.BaseNode{NodeName: "Custom"},
			FormatValue: func(value int32) string {
				v := value % 100
				if v < 0 {
					v += 100
				}
				return strconv.FormatInt(int64(v)+1, 10)
			},
			OnEnter: func(node *menu.ValueNode) {
				speed := node.Value() % 100
				if speed < 0 {
					speed += 100
				}
				log.Printf("swing at %d speed", speed)
				clubs.Swing(uint8(speed))
				log.Printf("swing done")
			},
		},
		settingsNode,
	)
	m := &menu.Menu{}
	m.SetRoot(root)

	serialUI := &menu.Serial{Serial: machine.Serial}
	encoderUI := &menu.Encoder{
		Encoder:   enc,
		ButtonPin: machine.GPIO8,
	}
	displayUI := &menu.GraphicsDisplay{
		Display:         &display,
		Font:            &proggy.TinySZ8pt7b,
		ForegroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.RGBA{A: 255},
		PaddingLeft:     1,
		PaddingTop:      -1,
		PaddingBottom:   1,
		Width:           40,
		Height:          32,
	}
	m.AddOutputs(serialUI, displayUI)
	m.AddInputs(serialUI, encoderUI)

	m.HandleInputs(context.Background())
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
