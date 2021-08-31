package main

import (
	CryptoRand "crypto/rand"
	"fmt"
	"math/big"
)

func ID(ring big.Int) *big.Int {
	var ID big.Int
	randInt, err := CryptoRand.Int(CryptoRand.Reader, &ring)
	if err != nil {
		panic(err)
	}
	(&ID).Set(randInt)
	return &ID
}

func GetBinary(ID *big.Int) string {
	return fmt.Sprintf("%b", ID)
}

func Ring(order int) big.Int {
	var ring big.Int
	ring.Exp(big.NewInt(2), big.NewInt(int64(order)), nil)

	return ring
}

//func GetDistance(a, b big.Int) uint64 {
//	var distance uint64
//	var distBA, rA, rB big.Int
//
//	(&rA).Set(a)
//	(&rB).Set(b)
//
//	(&distBA).Sub(&rB, &rA)
//	(&distBA).Mod(&distBA, &gr.Ring)
//
//	//buf := bytes.NewReader(distBA.Bytes())
//	//binary.Read(buf, binary.BigEndian, &distance)
//
//	return distBA.Int64()
//}

func main() {
	ring := Ring(128)
	for i := 0; i < 20; i++ {
		id := ID(ring)
		fmt.Println(id, id.BitLen())
		fmt.Println(GetBinary(id), len(GetBinary(id)))
	}
}
