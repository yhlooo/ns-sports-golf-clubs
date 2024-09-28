package menu

import (
	"context"
	"sync"
)

// Menu 菜单
type Menu struct {
	lock   sync.RWMutex
	root   *Item
	cursor uint32

	inputs        []UIInput
	outputs       []UIOutput
	operationChan chan Operation
}

// AddItems 添加菜单项
func (m *Menu) AddItems(items ...*Item) {
	m.lock.Lock()
	if m.root == nil {
		m.root = &Item{}
	}
	m.root.AddSubItems(items...)
	m.lock.Unlock()
	m.Show()
}

// AddInputs 添加菜单输入
func (m *Menu) AddInputs(inputs ...UIInput) {
	m.inputs = append(m.inputs, inputs...)
}

// AddOutputs 添加菜单输出
func (m *Menu) AddOutputs(outputs ...UIOutput) {
	m.outputs = append(m.outputs, outputs...)
	m.Show()
}

// HandleInputs 开始处理输入，阻塞直到 ctx 结束
func (m *Menu) HandleInputs(ctx context.Context) {
	if m.operationChan == nil {
		m.operationChan = make(chan Operation)
	}
	for _, input := range m.inputs {
		input.StartReceiving(m.operationChan)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case op := <-m.operationChan:
			switch {
			case op.NextN != nil:
				m.NextN(op.NextN.N)
			case op.Enter != nil:
				m.Enter()
			case op.Back != nil:
				m.Back()
			}
		}
	}
}

// Show 显示菜单
func (m *Menu) Show() {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, output := range m.outputs {
		output.Show(m)
	}
}

// NextN 选择下 n 项，若 n 是负数表示上 -n 项
func (m *Menu) NextN(n int32) {
	m.lock.Lock()
	if n < 0 {
		m.cursor -= uint32(-n)
	} else {
		m.cursor += uint32(n)
	}
	m.lock.Unlock()
	m.Show()
}

// Enter 进入当前项
func (m *Menu) Enter() {
	m.lock.Lock()
	if m.root == nil || len(m.root.children) == 0 {
		m.lock.Unlock()
		return
	}
	i := m.cursor % uint32(len(m.root.children))
	m.root = m.root.children[i]
	m.cursor = 0
	m.lock.Unlock()
	if m.root.Run != nil {
		m.root.Run(m)
	}
	m.Show()
}

// Back 返回
func (m *Menu) Back() {
	m.lock.Lock()
	if m.root == nil || m.root.parent == nil {
		m.lock.Unlock()
		return
	}
	m.root = m.root.parent
	m.cursor = 0
	m.lock.Unlock()
	m.Show()
}

// ItemNames 返回选项名和当前所选项序号
func (m *Menu) ItemNames() (names []string, selected uint32) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.root == nil || len(m.root.children) == 0 {
		return
	}
	selected = m.cursor % uint32(len(m.root.children))
	for _, item := range m.root.children {
		names = append(names, item.Name)
	}
	return
}

// Item 菜单项
type Item struct {
	// 项名
	Name string
	// 进入菜单项执行
	Run func(m *Menu)

	parent   *Item
	children []*Item
}

// AddSubItems 添加子项
func (item *Item) AddSubItems(items ...*Item) *Item {
	for _, subItem := range items {
		subItem.parent = item
		item.children = append(item.children, subItem)
	}
	return item
}

// BackItem 返回上一级菜单选项
func BackItem(name string) *Item {
	if name == "" {
		name = "Back"
	}
	return &Item{Name: name, Run: BackFn}
}

// BackFn 返回上一级菜单
func BackFn(m *Menu) {
	m.Back()
	m.Back()
}
