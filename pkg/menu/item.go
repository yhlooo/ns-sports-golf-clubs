package menu

import "strconv"

// Node 菜单节点
type Node interface {
	// Name 返回当前节点名
	Name() string
	// Back 退出当前节点，返回退出后的节点
	Back() Node
	// Enter 进入当前节点所选项，返回进入后的节点
	Enter() Node
	// NextN 选择下 n 项，若 n 是负数表示上 -n 项
	NextN(n int32)
	// Items 返回当前节点的子项和所选项序号
	Items() (names []string, selected int32)

	// AddChildren 添加子节点
	AddChildren(children ...Node)
	// SetParent 设置父节点
	SetParent(parent Node)
}

// BaseNode Node 的一个基础实现
type BaseNode struct {
	NodeName string

	parent   Node
	cursor   int32
	children []Node
}

var _ Node = (*BaseNode)(nil)

// Name 返回当前节点名
func (node *BaseNode) Name() string {
	return node.NodeName
}

// Back 退出当前节点，返回退出后的节点
func (node *BaseNode) Back() Node {
	if node.parent == nil {
		return node
	}
	return node.parent
}

// Enter 进入当前节点，返回进入后的节点
func (node *BaseNode) Enter() Node {
	if len(node.children) == 0 {
		return node
	}
	i := node.cursor % int32(len(node.children))
	if i < 0 {
		i += int32(len(node.children))
	}
	return node.children[i]
}

// NextN 选择下 n 项，若 n 是负数表示上 -n 项
func (node *BaseNode) NextN(n int32) {
	node.cursor += n
}

// Items 返回当前节点的子项和所选项序号
func (node *BaseNode) Items() (names []string, selected int32) {
	if len(node.children) == 0 {
		return
	}
	selected = node.cursor % int32(len(node.children))
	if selected < 0 {
		selected += int32(len(node.children))
	}
	for _, child := range node.children {
		names = append(names, child.Name())
	}
	return
}

// AddChildren 添加子节点
func (node *BaseNode) AddChildren(children ...Node) {
	for _, child := range children {
		child.SetParent(node)
		node.children = append(node.children, child)
	}
}

// SetParent 设置父节点
func (node *BaseNode) SetParent(parent Node) {
	node.parent = parent
}

// NewBoolValueNode 创建存储布尔值的 *ValueNode
func NewBoolValueNode(name string, val bool, nameWithValue bool, onEnter func(bool)) *ValueNode {
	node := &ValueNode{
		BaseNode: BaseNode{NodeName: name},
		FormatValue: func(value int32) string {
			if value%2 != 0 {
				return "true"
			} else {
				return "false"
			}
		},
		OnEnter: func(node *ValueNode) {
			if onEnter != nil {
				onEnter(node.Value()%2 != 0)
			}
		},
	}
	if val {
		node.SetValue(1)
	}
	if nameWithValue {
		node.NodeName = nodeNameWithValue
	}
	return node
}

// ValueNode 值节点， Node 的实现
type ValueNode struct {
	BaseNode
	// 节点名
	NodeName func(node *ValueNode) string
	// 将值格式化为字符串
	FormatValue func(value int32) string
	// 进入当前节点所选项时执行
	OnEnter func(node *ValueNode)
}

var _ Node = (*ValueNode)(nil)

// Name 返回当前节点名
func (node *ValueNode) Name() string {
	if node.NodeName != nil {
		return node.NodeName(node)
	}
	return node.BaseNode.Name()
}

// Enter 进入当前节点，返回进入后的节点
func (node *ValueNode) Enter() Node {
	if node.OnEnter != nil {
		node.OnEnter(node)
	}
	// 执行完后返回父节点
	return node.Back()
}

// Items 返回当前节点的子项和所选项序号
func (node *ValueNode) Items() (names []string, selected int32) {
	if node.FormatValue != nil {
		return []string{node.FormatValue(node.cursor)}, 0
	}
	return []string{strconv.FormatInt(int64(node.cursor), 10)}, 0
}

// AddChildren 添加子节点
func (node *ValueNode) AddChildren(_ ...Node) {}

// SetValue 设置值
func (node *ValueNode) SetValue(v int32) {
	node.cursor = v
}

// Value 获取值
func (node *ValueNode) Value() int32 {
	return node.cursor
}

// nodeNameWithValue 返回包含值的节点名
func nodeNameWithValue(node *ValueNode) string {
	return node.BaseNode.NodeName + ": " + node.FormatValue(node.Value())
}

// ActionNode 动作节点， Node 的实现
type ActionNode struct {
	BaseNode
	// 进入当前节点所选项时执行
	OnEnter func(node *ActionNode)
}

var _ Node = (*ActionNode)(nil)

// Enter 进入当前节点，返回进入后的节点
func (node *ActionNode) Enter() Node {
	if node.OnEnter != nil {
		node.OnEnter(node)
	}
	// 执行完后返回父节点
	return node.Back()
}

// AddChildren 添加子节点
func (node *ActionNode) AddChildren(_ ...Node) {}

// NewBackNode 创建 *BackNode
func NewBackNode(name string) *BackNode {
	return &BackNode{BaseNode: BaseNode{NodeName: name}}
}

// BackNode 返回节点
type BackNode struct {
	BaseNode
}

var _ Node = (*BackNode)(nil)

// Enter 进入当前节点，返回进入后的节点
func (node *BackNode) Enter() Node {
	// 返回父的父节点
	return node.Back()
}

// AddChildren 添加子节点
func (node *BackNode) AddChildren(_ ...Node) {}
