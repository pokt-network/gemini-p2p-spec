package tools

import (
	bytes "bytes"
	rand "math/rand"
	sort "sort"
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

func BinarySearch(s [][]byte, t []byte) []byte {
	el := sort.Search(len(s), func(i int) bool {
		return bytes.Compare(s[i], t) == 0
	})

	return s[el]
}

func PickRandom(sob [][]byte) []byte {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(len(sob) - 1)
	return sob[randomInt]
}
