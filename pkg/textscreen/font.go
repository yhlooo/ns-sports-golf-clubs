package textscreen

import (
	"image/color"
	"io"
	"strings"

	"tinygo.org/x/drivers"
	"tinygo.org/x/tinyfont"
)

// NewTextScreen 创建一个显示文本的屏幕
func NewTextScreen(
	display drivers.Displayer,
	x, y, maxWidth, maxHeight int16,
	f tinyfont.Fonter,
	fg, bg color.RGBA,
	lineSpace int16,
) io.Writer {
	return &textScreen{
		display:   display,
		x:         x,
		y:         y,
		maxHeight: maxHeight,
		maxWidth:  maxWidth,
		font:      f,
		fg:        fg,
		bg:        bg,
		lineSpace: lineSpace,
	}
}

// textScreen 显示文本的屏幕
type textScreen struct {
	display   drivers.Displayer
	x         int16
	y         int16
	maxWidth  int16
	maxHeight int16

	font      tinyfont.Fonter
	fg        color.RGBA
	bg        color.RGBA
	lineSpace int16

	bufferLines []string
}

var _ io.Writer = &textScreen{}

// Write 往屏幕写入数据
func (s *textScreen) Write(p []byte) (int, error) {
	// 更新 buffer
	var bufferLines []string
	if len(s.bufferLines) == 0 {
		bufferLines = WrapText(s.font, string(p), s.maxWidth)
	} else {
		bufferLines = append(
			s.bufferLines[:len(s.bufferLines)-1],
			WrapText(s.font, s.bufferLines[len(s.bufferLines)-1]+string(p), s.maxWidth)...,
		)
	}

	// 裁剪上方显示不下的行
	for len(bufferLines) > 0 && int16(len(bufferLines))*(int16(s.font.GetYAdvance())+s.lineSpace) > s.maxHeight {
		bufferLines = bufferLines[1:]
	}

	// 显示
	oldLen := len(s.bufferLines)
	y := s.y + int16(s.font.GetYAdvance())
	for i, line := range bufferLines {
		if i < oldLen {
			// 清除原有数据
			tinyfont.WriteLine(s.display, s.font, s.x, y, s.bufferLines[i], s.bg)
		}
		tinyfont.WriteLine(s.display, s.font, s.x, y, line, s.fg)
		y += int16(s.font.GetYAdvance()) + s.lineSpace
	}
	if err := s.display.Display(); err != nil {
		return 0, err
	}

	s.bufferLines = bufferLines
	return len(p), nil
}

// WrapText 给文本换行
func WrapText(f tinyfont.Fonter, str string, maxWidth int16) []string {
	var lines []string
	for _, l := range strings.Split(str, "\n") {
		if l == "" {
			lines = append(lines, "")
			continue
		}
		for l != "" {
			cnt := uint32(len(l))
			for _, width := tinyfont.LineWidth(f, l[:cnt]); width > uint32(maxWidth); {
				if newCnt := cnt * uint32(maxWidth) / width; newCnt < cnt {
					cnt = newCnt
				} else {
					cnt--
				}
				_, width = tinyfont.LineWidth(f, l[:cnt])
			}
			lines = append(lines, l[:cnt])
			l = l[cnt:]
		}
	}
	return lines
}

// WriteLines 输出文本行到显示器
func WriteLines(
	display drivers.Displayer,
	font tinyfont.Fonter,
	x, y, maxWidth, lineSpace int16,
	str string,
	c color.RGBA,
) {
	y += int16(font.GetYAdvance())
	lines := strings.Split(str, "\n")
	for _, l := range lines {
		if l == "" {
			y += int16(font.GetYAdvance()) + lineSpace
			continue
		}
		for l != "" {
			cnt := uint32(len(l))
			for _, width := tinyfont.LineWidth(font, l[:cnt]); width > uint32(maxWidth); {
				if newCnt := cnt * uint32(maxWidth) / width; newCnt < cnt {
					cnt = newCnt
				} else {
					cnt--
				}
				_, width = tinyfont.LineWidth(font, l[:cnt])
			}
			tinyfont.WriteLine(display, font, x, y, l[:cnt], c)
			y += int16(font.GetYAdvance()) + lineSpace
			l = l[cnt:]
		}
	}
}
