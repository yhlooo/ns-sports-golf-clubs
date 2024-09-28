package menu

import (
	"image/color"

	"tinygo.org/x/drivers"
	"tinygo.org/x/tinyfont"
)

// GraphicsDisplay 基于图形显示器的菜单用户交互界面输出
type GraphicsDisplay struct {
	// 显示设备
	Display drivers.Displayer
	// 显示字体
	Font tinyfont.Fonter
	// 前景色
	ForegroundColor color.RGBA
	// 背景色
	BackgroundColor color.RGBA
	// 显示位置左上角 X 坐标
	X int16
	// 显示位置左上角 Y 坐标
	Y int16
	// 显示高度
	Height int16
	// 显示宽度
	Width int16
	// 字符上方填充空白
	PaddingTop int16
	// 字符下方填充空白
	PaddingBottom int16
	// 字符左侧填充空白
	PaddingLeft int16
}

var _ UIOutput = (*GraphicsDisplay)(nil)

// Show 显示菜单当前状态
func (g *GraphicsDisplay) Show(m *Menu) {
	// 计算显示区域
	dispWidth, dispHeight := g.Display.Size()
	height := g.Height
	if height == 0 {
		height = dispHeight - g.Y
	}
	width := g.Width
	if width == 0 {
		width = dispWidth - g.X
	}
	lineHeight := int16(g.Font.GetYAdvance()) + g.PaddingTop + g.PaddingBottom
	midLineY := (height - lineHeight) / 2

	names, selected := m.ItemNames()

	// 中间显示选中行
	for y := midLineY; y < midLineY+lineHeight; y++ {
		for x := int16(0); x < width; x++ {
			g.Display.SetPixel(x+g.X, y+g.Y, g.ForegroundColor)
		}
	}
	tinyfont.WriteLine(
		g.Display, g.Font,
		g.X+g.PaddingLeft, midLineY+lineHeight+g.Y-g.PaddingBottom-1,
		names[selected], g.BackgroundColor,
	)

	// 显示上方行
	minY := midLineY
	for i := int16(1); midLineY-i*lineHeight >= 0 && int16(selected)-i >= 0; i++ {
		for y := midLineY - i*lineHeight; y < midLineY-(i-1)*lineHeight; y++ {
			for x := int16(0); x < width; x++ {
				g.Display.SetPixel(x+g.X, y+g.Y, g.BackgroundColor)
			}
		}
		y := midLineY - i*lineHeight
		minY = y
		tinyfont.WriteLine(
			g.Display, g.Font,
			g.X+g.PaddingLeft, y+lineHeight+g.Y-g.PaddingBottom-1,
			names[int16(selected)-i], g.ForegroundColor,
		)
	}
	for y := int16(0); y < minY; y++ {
		for x := int16(0); x < width; x++ {
			g.Display.SetPixel(x+g.X, y+g.Y, g.BackgroundColor)
		}
	}

	// 显示下方行
	maxY := midLineY + lineHeight
	for i := int16(1); midLineY+(i+1)*lineHeight < height && int16(selected)+i < int16(len(names)); i++ {
		for y := midLineY + i*lineHeight; y < midLineY+(i+1)*lineHeight; y++ {
			for x := int16(0); x < width; x++ {
				g.Display.SetPixel(x+g.X, y+g.Y, g.BackgroundColor)
			}
		}
		y := midLineY + i*lineHeight
		maxY = y + lineHeight
		tinyfont.WriteLine(
			g.Display, g.Font,
			g.X+g.PaddingLeft, y+lineHeight+g.Y-g.PaddingBottom-1,
			names[int16(selected)+i], g.ForegroundColor,
		)
	}
	for y := maxY; y < height; y++ {
		for x := int16(0); x < width; x++ {
			g.Display.SetPixel(x+g.X, y+g.Y, g.BackgroundColor)
		}
	}

	_ = g.Display.Display()
}
