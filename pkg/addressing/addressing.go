package addressing

import (
	"crypto/sha1"
	"fmt"
	"math/big"

	gUUID "github.com/google/uuid"
)

const Delimiter string = "@"

type Status string

const (
	Hashed Status = "Hashed"
	Raw    Status = "Organic"
)

type (
	Addr interface {
		Hash()
		GetHash() []byte
		GetBinaryHash() []byte
		GetRaw() string
		GetUUID() []byte
		GetStatus() Status
		String() string
		GetBitLength() int
		GetBinBitLength() int
	}

	Address struct {
		ID     []byte
		Raw    string
		Hashed []byte
		Status Status
	}
)

func NewAddress(addr string, hash ...bool) *Address {
	a := &Address{
		ID:     gUUID.NodeID(),
		Raw:    addr,
		Hashed: nil,
		Status: Raw,
	}

	if len(hash) > 0 && hash[0] == true {
		a.Hash()
	}

	return a
}

func (a *Address) Hash() {
	if a.Status == Hashed {
		return
	}

	// TODO: decouple this using an interface
	h := sha1.New()
	h.Write(a.ID)
	h.Write([]byte(Delimiter))
	h.Write([]byte(a.Raw))
	a.Hashed = h.Sum(nil)
	a.Status = Hashed
}

func (a *Address) GetHash() []byte {
	if a.Status == Hashed {
		return a.Hashed
	} else {
		return nil
	}
}

func (a *Address) GetStatus() Status {
	return a.Status
}

func (a *Address) GetRaw() string {
	return a.Raw
}

func (a *Address) GetUUID() []byte {
	return a.ID
}

func (a *Address) String() string {
	return fmt.Sprintf("Raw: %s, Hash: %b", a.Raw, a.Hashed)
}

func (a *Address) GetBinaryHash() []byte {
	var binRep big.Int
	(&binRep).SetBytes(a.Hashed)
	return []byte(fmt.Sprintf("%b", &binRep))
}

func (a *Address) GetBitLength() int {
	if a.Status == Hashed {
		return len(a.Hashed) * 8
	}
	return -1
}

func (a *Address) GetBinBitLength() int {
	if a.Status == Hashed {
		return len(a.GetBinaryHash())
	}
	return -1
}
