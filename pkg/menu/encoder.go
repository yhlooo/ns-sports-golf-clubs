package menu

import (
	"fmt"
	"machine"
	"time"

	"github.com/yhlooo/ns-sports-golf-clubs/pkg/encoder"
)

// Encoder 基于编码器的菜单用户交互界面输入源的实现
type Encoder struct {
	// 编码器
	Encoder *encoder.Encoder
	// 按钮针脚
	ButtonPin machine.Pin
}

var _ UIInput = (*Encoder)(nil)

// StartReceiving 开始接收操作，并将操作输入到 ch
func (e *Encoder) StartReceiving(ch chan<- Operation) {
	btnPress := false
	btnLastPress := time.Now()
	if err := e.ButtonPin.SetInterrupt(machine.PinFalling, func(pin machine.Pin) {
		btnPress = true
		btnLastPress = time.Now()
	}); err != nil {
		panic(fmt.Errorf("set interrupt error: %w", err))
	}

	go func() {
		for {
			ch <- Operation{} // 发送一个空操作检查菜单是否可操作
			value := e.Encoder.Value()
			time.Sleep(time.Millisecond)
			if btnPress {
				now := time.Now()
				switch {
				case now.Sub(btnLastPress) > 3*time.Second:
					// 过期了
				case e.ButtonPin.Get():
					// 短按进入
					ch <- Operation{Enter: &Enter{}}
				case now.Sub(btnLastPress) > 500*time.Millisecond:
					// 长按退出
					ch <- Operation{Back: &Back{}}
				default:
					// 可能是长按，再等等
					continue
				}
				btnPress = false
			}
			if cur := e.Encoder.Value(); cur != value {
				ch <- Operation{NextN: &NextN{N: cur - value}}
				continue
			}
		}
	}()
}
