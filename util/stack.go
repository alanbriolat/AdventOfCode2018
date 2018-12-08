package util

import "github.com/cheekybits/genny/generic"

type Generic generic.Type

type GenericStack struct {
	Data []Generic
	Nil Generic
}

func NewGenericStack(initSize int) GenericStack {
	return GenericStack{Data: make([]Generic, 0, initSize)}
}

func (s *GenericStack) Count() int {
	return len(s.Data)
}

func (s *GenericStack) Push(x Generic) {
	s.Data = append(s.Data, x)
}

func (s *GenericStack) Pop() (Generic, bool) {
	if result, ok := s.Peek(); !ok {
		return s.Nil, false
	} else {
		s.Data = s.Data[:len(s.Data)-1]
		return result, true
	}
}

func (s *GenericStack) Peek() (Generic, bool) {
	last := len(s.Data) - 1
	if last < 0 {
		return s.Nil, false
	} else {
		return s.Data[last], true
	}
}

func (s *GenericStack) Top() Generic {
	if len(s.Data) > 0 {
		return s.Data[len(s.Data)-1]
	} else {
		return s.Nil
	}
}
