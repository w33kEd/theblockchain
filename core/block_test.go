package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/w33ked/theblockchain/types"
)

func randomBlock(height uint32) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: types.RandomHash(),
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	tx := Transaction{
		Data: []byte("foo"),
	}

	return NewBlock(header, []Transaction{tx})
}
func TestHashBlock(t *testing.T) {
	b := randomBlock(0)
	fmt.Println(b.Hash(BlockHasher{}))
}
