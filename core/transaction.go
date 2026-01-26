package core

import "github.com/w33ked/theblockchain/crypto"

type Transaction struct {
	Data      []byte
	PublicKey crypto.PublicKey
	Signature *crypto.Signature
}
