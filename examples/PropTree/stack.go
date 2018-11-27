package main

import "container/list"

type Stack struct {
	l *list.List
}

func NewStack() *Stack {
	return &Stack{l: list.New()}
}

func (s *Stack) Push(v interface{}) {
	s.l.PushBack(v)
}

func (s *Stack) Pop() interface{} {
	if e := s.l.Back(); e != nil {
		s.l.Remove(e)
		return e.Value
	}
	return nil
}

func (s *Stack) Poll(i int) interface{} {
	if e := s.l.Back(); e != nil {
		for ; i > 0 && e != nil; i-- {
			e = e.Prev()
		}
		if e != nil {
			return e.Value
		}
	}
	return nil
}

func (s *Stack) Depth() int {
	return s.l.Len()
}

func (s *Stack) Reset() {
	s.l.Init()
}
