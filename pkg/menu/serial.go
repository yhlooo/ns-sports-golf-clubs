package menu

import (
	"fmt"
	"machine"
	"time"
)

// Serial 基于串口的菜单用户交互界面的实现
type Serial struct {
	// 接收输入和发送输出的串口
	Serial machine.Serialer
}

var _ UIOutput = (*Serial)(nil)
var _ UIInput = (*Serial)(nil)

// Show 显示菜单当前状态
func (s *Serial) Show(m *Menu) {
	content := "\x1b[100A\x1b[100D\x1b[2J"
	items, selected := m.ItemNames()
	for i, item := range items {
		if uint32(i) == selected {
			content += "\x1b[7m" + item + "\x1b[0m\r\n"
		} else {
			content += item + "\r\n"
		}
	}
	_, _ = fmt.Fprint(s.Serial, content)
}

// StartReceiving 开始接收操作
func (s *Serial) StartReceiving(ch chan<- Operation) {
	go func() {
		input := ""
		for {
			ch <- Operation{} // 发送一个空操作检查菜单是否可操作

			for {
				c, err := s.Serial.ReadByte()
				if err != nil {
					time.Sleep(time.Millisecond)
					continue
				}
				input += string(c)
				break
			}
			switch input {
			case "\x1b[A": // 上
				ch <- Operation{NextN: &NextN{N: -1}}
			case "\x1b[B": // 下
				ch <- Operation{NextN: &NextN{N: 1}}
			case "\x1b[D": // 左
				ch <- Operation{Back: &Back{}}
			case "\x1b[C", "\r": // 右、回车
				ch <- Operation{Enter: &Enter{}}
			case "\x1b", "\x1b[": // 输入一半
				continue
			default:
				_, _ = fmt.Fprintf(s.Serial, "%q\r\n", input)
			}
			input = ""
		}
	}()
}
