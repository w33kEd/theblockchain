package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/w33ked/theblockchain/crypto"
)

func TestAccountStateTransferFail(t *testing.T) {
	state := NewAccountState()

	from := crypto.GeneratePrivateKey().PublicKey().Address()
	to := crypto.GeneratePrivateKey().PublicKey().Address()
	amount := uint64(90)

	assert.NotNil(t, state.Transfer(from, to, amount))

}

func TestAccountStateTransferSuccess(t *testing.T) {
	state := NewAccountState()
	from := crypto.GeneratePrivateKey().PublicKey().Address()

	state.AddBalance(from, uint64(100))

	to := crypto.GeneratePrivateKey().PublicKey().Address()
	amount := uint64(90)

	assert.Nil(t, state.Transfer(from, to, amount))

}
