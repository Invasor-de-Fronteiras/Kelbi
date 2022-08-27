package stringstack

import (
	"errors"
)

// StringStack is a basic LIFO "stack" for storing strings.
type StringStack struct {
	Locked bool     `json:"locked"`
	Stack  []string `json:"stack"`
}

// New creates a new instance of StringStack
func New() *StringStack {
	return &StringStack{Locked: false}
}

// Set sets up a new StringStack
func (s *StringStack) Set(v string) {
	s.Stack = []string{v}
}

// Lock freezes the StringStack
func (s *StringStack) Lock() {
	if !s.Locked {
		s.Locked = true
	}
}

// Unlock unfreezes the StringStack
func (s *StringStack) Unlock() {
	if s.Locked {
		s.Locked = false
	}
}

// Push pushes a string onto the stack.
func (s *StringStack) Push(v string) {
	s.Stack = append(s.Stack, v)
}

// Pop pops a string from the stack.
func (s *StringStack) Pop() (string, error) {
	if len(s.Stack) == 0 {
		return "", errors.New("no items on stack")
	}

	x := s.Stack[len(s.Stack)-1]
	s.Stack = s.Stack[:len(s.Stack)-1]

	return x, nil
}
