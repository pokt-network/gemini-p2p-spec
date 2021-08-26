package ring

import (
	"encoding/binary"
	"math/big"
)

type (
	Ring interface {
		GetDistance([]byte, []byte) uint64
	}

	GeminiRing struct {
		Order int
		Ring  big.Int
	}
)

func NewGeminiRing(order int) *GeminiRing {
	var ring big.Int
	ring.Exp(big.NewInt(2), big.NewInt(int64(order)), nil)

	return &GeminiRing{
		Order: order,
		Ring:  ring,
	}
}

func (gr *GeminiRing) GetDistance(a, b []byte) uint64 {
	var distBA, rA, rB big.Int

	(&rA).SetBytes(a)
	(&rB).SetBytes(b)

	(&distBA).Sub(&rB, &rA)
	(&distBA).Mod(&distBA, &gr.Ring)

	return binary.BigEndian.Uint64(distBA.Bytes())
}
