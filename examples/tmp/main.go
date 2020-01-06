package main

import (
	"sync"
)

type Struct struct{ sync.Mutex }

func main() {
	var s Struct
	s.A()
}

func (s *Struct) A() {
	s.Lock()
	s.B()
	s.Unlock()
}

func (s *Struct) B() {
	s.Lock()
	s.Unlock()
}
