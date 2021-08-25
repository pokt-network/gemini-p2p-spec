package gemini

import (
	bytes "bytes"
	"errors"
	Addressing "gemelos/pkg/addressing"
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
		SetInHatClub(int, Addressing.Addr)
		SetInBootClub(int, Addressing.Addr)
		GetHatClub() *[]Addressing.Addr
		GetBootClub() *[]Addressing.Addr
		IsInHatClub(string) bool
		IsInBootClub(string) bool
		GetAddrDistance(string) int
		Route(string, payload []byte) (Addressing.Addr, RoutingStatus)
		SearchState(Case, string) Addressing.Addr
	}

	GeminiParams struct {
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
		State    []Addressing.Addr
		HatClub  Range
		BootClub Range
	}
)

func NewGeminus(addr string) *Geminus {
	gAddr := Addressing.NewAddress(addr)
	geminiParams := GeminiParams{
		AddrLength: 64,
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
		Params: geminiParams,
		Addr:   gAddr,
		// this allocation should be dynamic per need
		State: make([]Addressing.Addr, geminiParams.ClubSize.Hat+geminiParams.ClubSize.Boot, geminiParams.ClubSize.Hat+geminiParams.ClubSize.Boot),
		HatClub: Range{
			Start: 0,
			End:   geminiParams.ClubSize.Hat - 1,
		},
		BootClub: Range{
			Start: geminiParams.ClubSize.Hat,
			End:   (geminiParams.ClubSize.Boot + geminiParams.ClubSize.Hat) - 1,
		},
	}
}

func (g *Geminus) Init() error {
	g.Addr.Hash()
	if len(g.Addr.GetHash()) != g.Params.AddrLength {
		return errors.New("Wrong Gemini Address Length Param or Faulty Hash Function")
	}
	return nil
}

func (g *Geminus) GetState() []Addressing.Addr {
	return g.State
}

func (g *Geminus) SetState(addr string) int {
	isInHatClub := g.IsInBootClub(addr)
	isInBootClub := g.IsInBootClub(addr)

	if isInHatClub {
		// sort
		// add to proper position (according to numerical distance)
		haddr := Addressing.NewAddress(addr, true)
		d, _ := Tools.GetLSDistance(g.Addr.GetHash(), haddr.GetHash()) // not sure if this is the right "numerical distance" fn
		g.SetInHatClub(d, haddr)
		return 0
	} else if isInBootClub {
		// sort
		// add to proper position (according to numerical distance)
		haddr := Addressing.NewAddress(addr, true)
		d, _ := Tools.GetLSDistance(g.Addr.GetHash(), haddr.GetHash()) // not sure if this is the right "numerical distance" fn
		g.SetInBootClub(d, haddr)
		return 0
	}
	return 1
}

func (g *Geminus) GetHatClub() []Addressing.Addr {
	return g.State[g.HatClub.Start:g.HatClub.End]
}

func (g *Geminus) GetBootClub() []Addressing.Addr {
	return g.State[g.BootClub.Start:g.BootClub.End]
}

func (g *Geminus) SetInHatClub(d int, v Addressing.Addr) {
	g.State[g.HatClub.Start+d] = v
}

func (g *Geminus) SetInBootClub(d int, v Addressing.Addr) {
	g.State[g.BootClub.Start+d] = v
}

func (g *Geminus) IsInHatClub(addr string) bool {
	haddr := Addressing.NewAddress(addr, true)
	hatStart, hatEnd := 0, g.Params.HatLength
	myHatCase, addrHatCase := g.Addr.GetHash()[hatStart:hatEnd], haddr.GetHash()[hatStart:hatEnd]

	return bytes.Compare(myHatCase, addrHatCase) == 0
}

func (g *Geminus) IsInBootClub(addr string) bool {
	haddr := Addressing.NewAddress(addr, true)
	bootStart, bootEnd := (g.Params.AddrLength-1)-g.Params.BootLength, (g.Params.AddrLength - 1)
	myBootCase, addrBootCase := g.Addr.GetHash()[bootStart:bootEnd], haddr.GetHash()[bootStart:bootEnd]

	return bytes.Compare(myBootCase, addrBootCase) == 0
}

func (g *Geminus) GetAddrDistance(addr string) int {
	haddr := Addressing.NewAddress(addr, true)
	d, err := Tools.GetLSDistance(g.Addr.GetHash(), haddr.GetHash())
	if err != nil {
		return -1
	}

	return d
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
