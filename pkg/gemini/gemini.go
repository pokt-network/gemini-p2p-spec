package gemini

import (
	bytes "bytes"
	Tools "gemelos/pkg/tools"
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
		GetState() [][]byte
		SetState(addr []byte) int
		GetHatClub() [][]byte
		GetBootClub() [][]byte
		IsInHatClub([]byte) bool
		IsInBootClub([]byte) bool
		GetAddrDistance([]byte) int
		Route(destination []byte, payload []byte) ([]byte, RoutingStatus)
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
		Addr     []byte
		State    [][]byte
		HatClub  Range
		BootClub Range
	}
)

func NewGeminus(addr []byte) *Geminus {
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
		Addr:   addr,
		State:  make([][]byte, geminiParams.ClubSize.Hat+geminiParams.ClubSize.Boot),
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

func (g *Geminus) GetState() [][]byte {
	return g.State
}

func (g *Geminus) SetState(addr []byte) int {
	isInHatClub := true
	isInBootClub := false

	if isInHatClub {
		// sort
		// add to proper position (according to numerical distance)
		return 0
	} else if isInBootClub {
		// sort
		// add to proper position (according to numerical distance)
		return 0
	}
	return 1
}

func (g *Geminus) GetHatClub() [][]byte {
	return g.State[g.HatClub.Start:g.HatClub.End]
}

func (g *Geminus) GetBootClub() [][]byte {
	return g.State[g.BootClub.Start:g.BootClub.End]
}

func (g *Geminus) IsInHatClub(addr []byte) bool {
	hatStart, hatEnd := 0, g.Params.HatLength
	myHatCase, addrHatCase := g.Addr[hatStart:hatEnd], addr[hatStart:hatEnd]
	return bytes.Compare(myHatCase, addrHatCase) == 0
}

func (g *Geminus) IsInBootClub(addr []byte) bool {
	bootStart, bootEnd := (g.Params.AddrLength-1)-g.Params.BootLength, (g.Params.AddrLength - 1)
	myBootCase, addrBootCase := g.Addr[bootStart:bootEnd], addr[bootStart:bootEnd]

	return bytes.Compare(myBootCase, addrBootCase) == 0
}

func (g *Geminus) GetAddrDistance(addr []byte) int {
	d, err := Tools.GetLSDistance(g.Addr, addr)
	if err != nil {
		return -1
	}

	return d
}

func (g *Geminus) Route(destination []byte, payload []byte) ([]byte, RoutingStatus) {
	var foundAddr []byte
	var status RoutingStatus

	if g.IsInHatClub(destination) {
		foundAddr = Tools.BinarySearch(g.GetHatClub(), destination)
		status = HatFind
	} else if g.IsInBootClub(destination) {
		foundAddr = Tools.BinarySearch(g.GetBootClub(), destination)
		status = BootFind
	} else {
		foundAddr = Tools.PickRandom(g.GetBootClub())
		status = Forward
	}

	if len(foundAddr) == 0 {
		return foundAddr, Undefined
	}

	return foundAddr, status
}
