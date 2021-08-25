package tools

import (
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
