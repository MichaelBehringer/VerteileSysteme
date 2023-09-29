package main

import "fmt"

// Creating our own stack
// no default stack in golang

func newStack() *Stack {
	return &Stack{}
}

func (s *Stack) push(item int) {
	s.items = append(s.items, item)
}

func (s *Stack) pop() (int, error) {
	if s.isEmpty() {
		return 0, fmt.Errorf("Stack is empty")
	}
	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]
	return item, nil
}

func (s *Stack) isEmpty() bool {
	return len(s.items) == 0
}

func initStack() {
	for i := 0; i < 30; i++ {
		stack.push(i)
	}
}
