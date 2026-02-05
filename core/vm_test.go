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
	// push	data := []byte{0x09, 0x0a, 0x02, 0x0a, 0x0b}stack
	data := []byte{0x09, 0x0a, 0x02, 0x0a, 0x0b}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())
	
	assert.Equal(t, byte(11), vm.stack[vm.sp])
}
