package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVM(t *testing.T) {
	// 1 + 2 = 3
	// 1
	// push stack
	//2
	// push stack
	// add
	// 3
	// push stack
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())
	// result := vm.stack.Pop()
	// assert.Equal(t, 11, result)
	value := vm.stack.Pop().([]byte)
	assert.Equal(t, "FOO", string(value))

	data = []byte{0x09, 0x0a, 0x02, 0x0a, 0x0e}
	vm = NewVM(data)
	assert.Nil(t, vm.Run())
	result := vm.stack.Pop().(int)
	assert.Equal(t, 7, result)

}

func TestStack(t *testing.T) {
	s := NewStack(128)

	s.Push(1)
	s.Push(2)

	value := s.Pop()
	assert.Equal(t, value, 1)

	value = s.Pop()
	assert.Equal(t, value, 2)

}
