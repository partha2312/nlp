package datastructures

import datastructures "github.com/partha2312/nlp/datastructures/doublylinkedlist"

type LRU interface {
	Get(key struct{}) *datastructures.DLLNode
	Put(key, value struct{})
}

type lru struct {
	capacity int
	cache    map[struct{}]*datastructures.DLLNode
	dll      datastructures.DoublyLinkedList
}

func NewLRU(capacity int) LRU {
	cache := make(map[struct{}]*datastructures.DLLNode)
	dll := datastructures.NewDoublyLinkedList()
	return &lru{capacity, cache, dll}
}

func (l *lru) Get(key struct{}) *datastructures.DLLNode {
	node, ok := l.cache[key]
	if !ok {
		return nil
	}
	l.dll.DeleteNode(node)
	l.dll.AddToHead(node)
	return node
}

func (l *lru) Put(key, value struct{}) {
	node, ok := l.cache[key]
	if !ok {
		if l.capacity == len(l.cache) {
			nodeToDel := l.dll.LastNode()
			delete(l.cache, nodeToDel.GetKey())
			l.dll.DeleteNode(nodeToDel)
		}
		node = datastructures.NewDLLNode(key, value)
		l.cache[key] = node
		l.dll.AddToHead(node)
	} else {
		node.SetValue(value)
		l.cache[key] = node
		l.dll.DeleteNode(node)
		l.dll.AddToHead(node)
	}
}
