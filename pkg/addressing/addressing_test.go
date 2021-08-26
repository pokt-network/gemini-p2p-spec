package addressing

import (
	"testing"
)

func TestNewAddress(t *testing.T) {
	addr := NewAddress("10.10.123.123")

	if addr.GetStatus() != Raw {
		t.Log("Faulty Address Status")
		t.Fail()
	}

	addr.Hash()

	if addr.GetStatus() != Hashed {
		t.Log("Faulty Address Status")
		t.Fail()
	}

	if addr.GetHash() == nil {
		t.Log("Faulty Address Hash")
		t.Fail()
	}

	if len(addr.GetHash()) != 64 {
		t.Log("Faulty Address Hash Length")
		t.Fail()
	}
}

func TestNewAddressWithImmediateHashing(t *testing.T) {
	addr := NewAddress("10.10.123.123", true) // true

	if addr.GetStatus() != Hashed {
		t.Log("Faulty Address Status")
		t.Fail()
	}

	if addr.GetHash() == nil {
		t.Log("Faulty Address Hash")
		t.Fail()
	}

	if len(addr.GetHash()) != 64 {
		t.Log("Faulty Address Hash Length")
		t.Fail()
	}
}
