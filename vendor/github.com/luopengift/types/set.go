package types

import (
	"sync"
)

type Set struct {
	s   map[interface{}]bool
	mux *sync.RWMutex
}

func NewSet(e ...interface{}) *Set {
	set := Set{
		s:   make(map[interface{}]bool),
		mux: new(sync.RWMutex),
	}
	set.Add(e...)
	return &set
}

// 添加一个元素
func (self *Set) Add(e ...interface{}) *Set {
	self.mux.Lock()
	for _, v := range e {
		self.s[v] = true
	}
	self.mux.Unlock()
	return self
}

// 删除一个元素
func (self *Set) Remove(v interface{}) *Set {
	self.mux.Lock()
	delete(self.s, v)
	self.mux.Unlock()
	return self
}

// 清空所有元素
func (self *Set) Clear() error {
	self.mux.Lock()
	for k, _ := range self.s {
		self.Remove(k)
	}
	self.mux.Unlock()
	return nil
}

func (self *Set) Contains(v interface{}) bool {
	self.mux.RLock()
	defer self.mux.RUnlock()
	return self.s[v]
}

// 获取元素集合
func (self *Set) Elements() []interface{} {
	elements := []interface{}{}
	for k := range self.s {
		elements = append(elements, k)
	}
	return elements
}

func (self *Set) Len() int {
	return len(self.s)
}

// 是否和其他set一致
func (self *Set) Same(o *Set) bool {
	if self.Len() != o.Len() {
		return false
	}
	for k := range self.s {
		if !o.Contains(k) {
			return false
		}
	}
	return true
}

// 并集
func (self *Set) Union(o *Set) *Set {
	union := NewSet(self.Elements()...) //新创建一个集合,以免影响原集合
	for _, v := range o.Elements() {
		union.Add(v)
	}
	return union
}

// 交集
func (self *Set) Inter(o *Set) *Set {
	inter := NewSet()
	for _, v := range self.Elements() {
		if o.Contains(v) {
			inter.Add(v)
		}
	}
	return inter
}

// 差集
func (self *Set) Diff(o *Set) *Set {
	diff := NewSet()
	for _, v := range self.Elements() {
		if !o.Contains(v) {
			diff.Add(v)
		}
	}
	return diff
}
