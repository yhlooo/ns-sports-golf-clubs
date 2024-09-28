package menu

// UIOutput 用户交互界面输出
type UIOutput interface {
	// Show 显示菜单当前状态
	Show(m *Menu)
}

// UIInput 用户交互界面输入源
type UIInput interface {
	// SetReceiveChan 设置接收输入操作的 channel
	SetReceiveChan(chan<- Operation)
}

// Operation 菜单操作
// NOTE: 以下各成员必须指定且仅可指定其中一项
type Operation struct {
	NextN *NextN
	Enter *Enter
	Back  *Back
}

// NextN 选择下或上 n 项操作
type NextN struct {
	N int32
}

// Enter 进入菜单操作
type Enter struct{}

// Back 返回操作
type Back struct{}
