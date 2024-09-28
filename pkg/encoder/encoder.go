package encoder

import (
	"fmt"
	"machine"
)

// Encoder 旋转编码器
type Encoder struct {
	// 接收编码器 A 相信号的针脚
	APin machine.Pin
	// 接收编码器 B 相信号的针脚
	BPin machine.Pin
	// 反转
	Reverse bool

	value int32
}

// Configure 配置编码器
func (e *Encoder) Configure() error {
	e.APin.Configure(machine.PinConfig{Mode: machine.PinInput})
	e.BPin.Configure(machine.PinConfig{Mode: machine.PinInput})
	if err := e.APin.SetInterrupt(machine.PinRising, func(_ machine.Pin) {
		if e.BPin.Get() != e.Reverse {
			e.value++
		} else {
			e.value--
		}
	}); err != nil {
		return fmt.Errorf("set interrupt error: %w", err)
	}
	return nil
}

// Value 当前值
func (e *Encoder) Value() int32 {
	return e.value
}

// SetValue 设置值
func (e *Encoder) SetValue(value int32) {
	e.value = value
}
