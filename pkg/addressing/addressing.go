package addressing

import (
	"crypto/sha512"
	"fmt"
)

type Status string

const (
	Hashed Status = "Hashed"
	Raw    Status = "Organic"
)

type (
	Addr interface {
		Hash()
		GetHash() []byte
		GetRaw() string
		GetStatus() Status
		String() string
	}

	Address struct {
		Raw    string
		Hashed []byte
		Status Status
	}
)

func NewAddress(addr string, hash ...bool) *Address {
	a := &Address{
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

	h := sha512.New()
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

func (a *Address) String() string {
	return fmt.Sprintf("Raw: %s, Hash: %b", a.Raw, a.Hashed)
}
