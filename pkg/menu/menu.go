package menu

import (
	"context"
	"sync"
)

// Menu 菜单
type Menu struct {
	lock sync.RWMutex
	root Node

	inputs        []UIInput
	outputs       []UIOutput
	operationChan chan Operation
}

// SetRoot 设置菜单根节点
func (m *Menu) SetRoot(root Node) {
	m.lock.Lock()
	m.root = root
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
	if m.root == nil {
		m.lock.Unlock()
		return
	}
	m.root.NextN(n)
	m.lock.Unlock()
	m.Show()
}

// Enter 进入当前项
func (m *Menu) Enter() {
	m.lock.Lock()
	if m.root == nil {
		m.lock.Unlock()
		return
	}
	m.root = m.root.Enter().Entered()
	m.lock.Unlock()
	m.Show()
}

// Back 返回
func (m *Menu) Back() {
	m.lock.Lock()
	if m.root == nil {
		m.lock.Unlock()
		return
	}
	m.root = m.root.Back()
	m.lock.Unlock()
	m.Show()
}

// ItemNames 返回选项名和当前所选项序号
func (m *Menu) ItemNames() (names []string, selected int32) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.root == nil {
		return
	}
	return m.root.Items()
}
