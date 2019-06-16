package datastructures

type DLLNode struct {
	key      struct{}
	value    struct{}
	previous *DLLNode
	next     *DLLNode
}

func NewDLLNode(key, value struct{}) *DLLNode {
	return &DLLNode{key: key, value: value}
}

func newdllEmptyNode() *DLLNode {
	return &DLLNode{}
}

func (d *DLLNode) SetValue(value struct{}) {
	d.value = value
}

func (d *DLLNode) GetKey() struct{} {
	return d.key
}

type DoublyLinkedList interface {
	AddToHead(node *DLLNode)
	DeleteNode(node *DLLNode)
	LastNode() *DLLNode
}

type doublyLinkedList struct {
	head *DLLNode
	tail *DLLNode
}

func NewDoublyLinkedList() DoublyLinkedList {
	head := newdllEmptyNode()
	tail := newdllEmptyNode()
	head.next = tail
	tail.previous = head
	return &doublyLinkedList{head, tail}
}

func (d *doublyLinkedList) AddToHead(node *DLLNode) {
	temp := d.head.next
	d.head.next = node
	node.previous = d.head
	node.next = temp
	temp.previous = node
}

func (d *doublyLinkedList) DeleteNode(node *DLLNode) {
	before := node.previous
	after := node.next
	before.next = after
	after.previous = before
}

func (d *doublyLinkedList) LastNode() *DLLNode {
	return d.tail.previous
}
