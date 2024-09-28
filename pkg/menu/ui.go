package menu

// UIOutput 用户交互界面输出
type UIOutput interface {
	// Show 显示菜单当前状态
	Show(m *Menu)
}

// UIInput 用户交互界面输入源
type UIInput interface {
	// StartReceiving 开始接收操作，并将操作输入到 ch
	StartReceiving(ch chan<- Operation)
}

// Operation 菜单操作
// NOTE: 以下各成员仅可指定其中一项
type Operation struct {
	// 选择下或上 n 项操作
	NextN *NextN
	// 进入菜单操作
	Enter *Enter
	// 返回操作
	Back *Back
}

// NextN 选择下或上 n 项操作
type NextN struct {
	N int32
}

// Enter 进入菜单操作
type Enter struct{}

// Back 返回操作
type Back struct{}
