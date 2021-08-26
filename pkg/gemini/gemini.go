package gemini

import (
	bytes "bytes"
	"errors"
	"fmt"
	Addressing "gemelos/pkg/addressing"
	Ring "gemelos/pkg/ring"
	Tools "gemelos/pkg/tools"
	"sort"
)

type Case string

const (
	Hat  Case = "Hat"
	Boot Case = "Boot"
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
		SetState(addr string) int
		SetInHatClub(Addressing.Addr)
		SetInBootClub(Addressing.Addr)
		GetHatClub() []Addressing.Addr
		GetBootClub() []Addressing.Addr
		IsInHatClub(string) bool
		IsInBootClub(string) bool
		GetAddrDistance(string) int
		Route(string, payload []byte) (Addressing.Addr, RoutingStatus)
		SearchState(Case, string) Addressing.Addr
	}

	GeminiParams struct {
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
		Params   GeminiParams
		Addr     Addressing.Addr
		HatClub  []Addressing.Addr
		BootClub []Addressing.Addr
	}
)

func NewGeminus(addr string) *Geminus {
	gAddr := Addressing.NewAddress(addr)
	gRing := Ring.NewGeminiRing(160)
	gParams := GeminiParams{
		Ring:       gRing,
		AddrLength: gRing.Order,
		HatLength:  5,
		BootLength: 5,
		ClubSize: struct {
			Hat  int
			Boot int
		}{
			Hat:  187,
			Boot: 187,
		},
	}

	return &Geminus{
		Params:   gParams,
		Addr:     gAddr,
		HatClub:  make([]Addressing.Addr, gParams.ClubSize.Hat, gParams.ClubSize.Hat),
		BootClub: make([]Addressing.Addr, gParams.ClubSize.Boot, gParams.ClubSize.Boot),
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

func (g *Geminus) SetState(addr string) int {
	isInHatClub := g.IsInHatClub(addr)
	isInBootClub := g.IsInBootClub(addr)

	haddr := Addressing.NewAddress(addr, true)
	distance := g.Params.Ring.GetDistance(g.Addr.GetHash(), haddr.GetHash())

	fmt.Printf("Found distance %d", distance)
	fmt.Printf("\nIs in Hat?: %t, Is in Boot?: %t (%s)", isInHatClub, isInBootClub, addr)

	if isInHatClub {
		g.SetInHatClub(haddr)
		return 0
	} else if isInBootClub {
		g.SetInBootClub(haddr)
		return 0
	}
	return 1
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

	fmt.Printf("My hash %b, their hash: %b, start: %d, end: %d\n", myHatCase, addrHatCase, hatStart, hatEnd)

	return bytes.Compare(myHatCase, addrHatCase) == 0
}

func (g *Geminus) IsInBootClub(addr string) bool {
	haddr := *Addressing.NewAddress(addr, true)
	bootStart, bootEnd := (g.Params.AddrLength-1)-g.Params.BootLength, (g.Params.AddrLength - 1)
	myBootCase, addrBootCase := g.Addr.GetBinaryHash()[bootStart:bootEnd], haddr.GetBinaryHash()[bootStart:bootEnd]

	fmt.Printf("My hash %b, their hash: %b, start: %d, end: %d\n", myBootCase, addrBootCase, bootStart, bootEnd)

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

	return haystack[el]
}

func (g *Geminus) Route(destination string, payload []byte) (Addressing.Addr, RoutingStatus) {
	var foundAddr Addressing.Addr
	var status RoutingStatus

	// a lot of functions like this hash the given []byte address
	// we should hash a given address once (from the arguments)
	// to avoid unnecessary resource waste

	// TODO: add numerical distance routing
	if g.IsInHatClub(destination) {
		foundAddr = g.SearchState(Hat, destination)
		status = HatFind
	} else if g.IsInBootClub(destination) {
		foundAddr = g.SearchState(Boot, destination)
		status = BootFind
	} else {
		foundAddr = g.GetBootClub()[Tools.PickRandom(1, g.Params.ClubSize.Boot)]
		status = Forward
	}

	if foundAddr == nil {
		return foundAddr, Undefined
	}

	return foundAddr, status
}
