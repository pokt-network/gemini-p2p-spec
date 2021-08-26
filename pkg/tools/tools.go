package tools

import (
	"encoding/binary"
	"math/big"
	rand "math/rand"
	"time"

	edlib "github.com/hbollon/go-edlib"
)

func GetLSDistance(s, t []byte) (int, error) {
	res, err := edlib.StringsSimilarity(string(s), string(t), edlib.Levenshtein)
	if err != nil {
		return 0, err
	}

	return int(res), err
}

func PickRandom(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(max - min)
	return randomInt
}

func BigIntToUint64(b big.Int) uint64 {
	return binary.BigEndian.Uint64(b.Bytes())
}
