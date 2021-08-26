package ring

import (
	"crypto/sha1"
	"encoding/binary"
	"testing"
)

func hash(str []byte) []byte {
	h := sha1.New()
	h.Write(str)
	return h.Sum(nil)
}

func TestNewAddress(t *testing.T) {
	gRing := NewGeminiRing(160)

	if gRing.Order != 160 {
		t.Log("Faulty Gemini Ring Order")
		t.Fail()
	}

	uint64v := binary.BigEndian.Uint64(gRing.Ring.Bytes())

	if uint64v == 0 {
		t.Log("Wrong Ring Type/Value")
		t.Fail()
	}
}

func TestGetDistance(t *testing.T) {
	gRing := NewGeminiRing(160)

	dist := gRing.GetDistance(
		hash([]byte("10.0.0.1")),
		hash([]byte("10.0.0.2")),
	)

	if dist == 0 {
		t.Log("Wrong distance calculation")
		t.Fail()
	}
}
