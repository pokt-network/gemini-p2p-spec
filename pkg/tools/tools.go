package tools

import (
	"encoding/binary"
	"math/big"
	rand "math/rand"
	"time"
)

func PickRandom(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(max - min)
	return randomInt
}

func BigIntToUint64(b big.Int) uint64 {
	return binary.BigEndian.Uint64(b.Bytes())
}
