package gemini

import (
	bytes "bytes"
	"errors"
	Addressing "gemelos/pkg/addressing"
	Ring "gemelos/pkg/ring"
	Tools "gemelos/pkg/tools"
	"sort"
)

type Case string

const (
	Hat     Case = "Hat"
	Boot    Case = "Boot"
	Foreign Case = "Foreign"
)

type RoutingStatus string

const (
	HatFind   RoutingStatus = "Hatfind"
	BootFind  RoutingStatus = "Bootfind"
	Forward                 = "Forward"
	Undefined               = "Undefined"
)

type (
	Gemini interface {
		Init() error
		GetState() []Addressing.Addr
		SetState(addr string) Case
		SetInHatClub(Addressing.Addr)
		SetInBootClub(Addressing.Addr)
		GetHatClub() []Addressing.Addr
		GetBootClub() []Addressing.Addr
		IsInHatClub(string) bool
		IsInBootClub(string) bool
		GetAddrDistance(string) int
		Route(string) (Addressing.Addr, RoutingStatus)
		SearchState(Case, string) Addressing.Addr
	}

	GeminiConfig struct {
		Ring       Ring.Ring
		AddrLength int
		HatLength  int
		BootLength int
		ClubSize   struct {
			Hat  int
			Boot int
		}
	}

	Range struct {
		Start int
		End   int
	}

	Geminus struct {
		Params   GeminiConfig
		Addr     Addressing.Addr
		HatClub  []Addressing.Addr
		BootClub []Addressing.Addr
	}
)

func NewGeminiConfig(networkCapacity, networkOrder, peerOrderA, peerOrderB int) *GeminiConfig {
	return &GeminiConfig{
		Ring:       Ring.NewGeminiRing(networkOrder),
		AddrLength: networkOrder,
		HatLength:  peerOrderA,
		BootLength: peerOrderB,
		ClubSize: struct {
			Hat  int
			Boot int
		}{
			Hat:  networkCapacity/2 ^ peerOrderA,
			Boot: networkCapacity/2 ^ peerOrderB,
		},
	}
}

func NewGeminus(addr string, networkCapacity, networkOrder, peerOrderA, peerOrderB int) *Geminus {
	gAddr := Addressing.NewAddress(addr)
	gParams := NewGeminiConfig(networkCapacity, networkOrder, peerOrderA, peerOrderB)

	return &Geminus{
		Params:   *gParams,
		Addr:     gAddr,
		HatClub:  make([]Addressing.Addr, 0, gParams.ClubSize.Hat),
		BootClub: make([]Addressing.Addr, 0, gParams.ClubSize.Boot),
	}
}

func (g *Geminus) Init() error {
	g.Addr.Hash()
	if len(g.Addr.GetHash()) != g.Params.AddrLength/8 {
		return errors.New("Wrong Gemini Address Length Param or Faulty Hash Function")
	}
	return nil
}

func (g *Geminus) GetState() []Addressing.Addr {
	return []Addressing.Addr(append(g.HatClub, g.BootClub...))
}

func (g *Geminus) SetState(addr string) Case {
	isInHatClub := g.IsInHatClub(addr)
	isInBootClub := g.IsInBootClub(addr)

	haddr := Addressing.NewAddress(addr, true)
	// what do we do with distance?
	// distance := g.Params.Ring.GetDistance(g.Addr.GetHash(), haddr.GetHash())

	if isInHatClub {
		g.SetInHatClub(haddr)
		return Hat
	} else if isInBootClub {
		g.SetInBootClub(haddr)
		return Boot
	}
	return Foreign
}

func (g *Geminus) GetHatClub() []Addressing.Addr {
	return g.HatClub
}

func (g *Geminus) GetBootClub() []Addressing.Addr {
	return g.BootClub
}

func (g *Geminus) SetInHatClub(v Addressing.Addr) {
	g.HatClub = append(g.HatClub, v)
}

func (g *Geminus) SetInBootClub(v Addressing.Addr) {
	g.BootClub = append(g.BootClub, v)
}

func (g *Geminus) IsInHatClub(addr string) bool {
	haddr := Addressing.NewAddress(addr, true)
	hatStart, hatEnd := 0, g.Params.HatLength-1
	myHatCase, addrHatCase := g.Addr.GetBinaryHash()[hatStart:hatEnd], haddr.GetBinaryHash()[hatStart:hatEnd]

	return bytes.Compare(myHatCase, addrHatCase) == 0
}

func (g *Geminus) IsInBootClub(addr string) bool {
	haddr := *Addressing.NewAddress(addr, true)

	myAddrLength, haddrLength := g.Addr.GetBinBitLength(), haddr.GetBinBitLength()

	haddrBootStart, haddrBootEnd := (haddrLength-1)-g.Params.BootLength, (haddrLength - 1)
	myBootStart, myBootEnd := (myAddrLength-1)-g.Params.BootLength, (myAddrLength - 1)

	myBootCase, addrBootCase := g.Addr.GetBinaryHash()[myBootStart:myBootEnd], haddr.GetBinaryHash()[haddrBootStart:haddrBootEnd]

	return bytes.Compare(myBootCase, addrBootCase) == 0
}

func (g *Geminus) SearchState(c Case, needle string) Addressing.Addr {
	var haystack []Addressing.Addr

	hneedle := Addressing.NewAddress(needle, true)

	if c == Hat {
		haystack = g.GetHatClub()
	} else if c == Boot {
		haystack = g.GetBootClub()
	}

	// initial thoughts, we can perhaps find by numerical distance aka o(1)
	// think about this later with otto
	el := sort.Search(len(haystack), func(i int) bool {
		return bytes.Compare(haystack[i].GetHash(), hneedle.GetHash()) == 0
	})

	if el < len(haystack) {
		return haystack[el]
	}

	return nil
}

func (g *Geminus) Route(destination string) (Addressing.Addr, RoutingStatus) {
	var foundAddr Addressing.Addr
	var status RoutingStatus

	// a lot of functions like this hash the given []byte address
	// we should hash a given address once (from the arguments)
	// to avoid unnecessary resource waste

	// TODO: add numerical distance routing
	if g.IsInHatClub(destination) {
		foundAddr = g.SearchState(Hat, destination)
		if foundAddr != nil {
			return foundAddr, HatFind
		}
	}

	if g.IsInBootClub(destination) {
		foundAddr = g.SearchState(Boot, destination)
		if foundAddr != nil {
			return foundAddr, BootFind
		}
	}

	if foundAddr == nil {
		bootClubSize := len(g.BootClub)
		// temporary if for testing purposes
		// in real life, a node will wait til it has seeded before it starts routing
		if bootClubSize > 0 {
			return g.GetBootClub()[Tools.PickRandom(1, bootClubSize+1)], Forward
		} else {
			status = Undefined
		}
	}

	return foundAddr, status
}
