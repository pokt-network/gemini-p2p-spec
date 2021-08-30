package gemini

import (
	bytes "bytes"
	"errors"
	Addressing "gemelos/pkg/addressing"
	Ring "gemelos/pkg/ring"
	Tools "gemelos/pkg/tools"
	"math"
)

type Club string

const (
	Hat          Club = "Hat"
	Boot              = "Boot"
	Unrecognized      = "Unrecognized"
)

type RoutingStatus string

const (
	HatRoute      RoutingStatus = "HatRoute"
	BootForward                 = "BootForward"
	RandomForward               = "RandomForward"
	Undefined                   = "Undefined"
)

type (
	Gemini interface {
		Init() error
		GetState() []Addressing.Addr
		SetState(addr string) (Club, error)
		AddInClub(Club, Addressing.Addr) error
		GetClub(Club) ([]Addressing.Addr, error)
		HaveSameClub(Club, Addressing.Addr, Addressing.Addr) (bool, error)
		BelongsInClub(Club, string) (bool, error)
		GetAddrDistance(string) int
		Route(string) (Addressing.Addr, RoutingStatus)
		SearchState(Club, string) Addressing.Addr
	}

	GeminiConfig struct {
		Ring       *Ring.GeminiRing
		AddrLength int
		HatLength  int
		BootLength int
		ClubSize   map[Club]int
	}

	Geminus struct {
		Params *GeminiConfig
		Addr   Addressing.Addr
		Clubs  map[Club][]Addressing.Addr
	}
)

func NewGeminiConfig(networkCapacity, networkOrder, hatLength, bootLength int) *GeminiConfig {
	return &GeminiConfig{
		Ring:       Ring.NewGeminiRing(networkOrder),
		AddrLength: networkOrder,
		HatLength:  hatLength,
		BootLength: bootLength,
		ClubSize: map[Club]int{
			Hat:  int(float64(networkCapacity) / math.Pow(2, float64(hatLength))),
			Boot: int(float64(networkCapacity) / math.Pow(2, float64(hatLength))),
		},
	}
}

func NewGeminus(addr string, gParams *GeminiConfig) *Geminus {
	gAddr := Addressing.NewAddress(addr)

	return &Geminus{
		Params: gParams,
		Addr:   gAddr,
		Clubs: map[Club][]Addressing.Addr{
			Hat:  make([]Addressing.Addr, 0, gParams.ClubSize[Hat]),
			Boot: make([]Addressing.Addr, 0, gParams.ClubSize[Boot]),
		},
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
	return []Addressing.Addr(append(g.Clubs[Hat], g.Clubs[Boot]...))
}

func (g *Geminus) SetState(addr string) (Club, error) {
	var club Club

	if c, _ := g.BelongsInClub(Hat, addr); c {
		club = Hat
	} else if c, _ := g.BelongsInClub(Boot, addr); c {
		club = Boot
	} else {
		return Unrecognized, errors.New("Unrecognized club/case")
	}

	haddr := Addressing.NewAddress(addr, true)

	g.AddInClub(club, haddr)

	return club, nil
}

func (g *Geminus) GetClub(club Club) ([]Addressing.Addr, error) {
	if club != Hat && club != Boot {
		return nil, errors.New("Unrecognized club/case")
	}
	return g.Clubs[club], nil
}

func (g *Geminus) AddInClub(club Club, v Addressing.Addr) error {
	if club != Hat && club != Boot {
		return errors.New("Unrecognized club/case")
	}

	g.Clubs[club] = append(g.Clubs[club], v)

	return nil
}

func (g *Geminus) HaveSameClub(club Club, haddrA Addressing.Addr, haddrB Addressing.Addr) (bool, error) {
	var caseA, caseB []byte

	if club == Hat {
		caseStart, caseEnd := 0, g.Params.HatLength-1
		caseA = haddrA.GetBinaryHash()[caseStart:caseEnd]
		caseB = haddrB.GetBinaryHash()[caseStart:caseEnd]
	} else if club == Boot {
		haddrALength, haddrBLength := haddrA.GetBinBitLength(), haddrB.GetBinBitLength()

		haddrACaseStart, haddrACaseEnd := (haddrALength-1)-g.Params.BootLength, (haddrALength - 1)
		haddrBCaseStart, haddrBCaseEnd := (haddrBLength-1)-g.Params.BootLength, (haddrBLength - 1)

		caseA = haddrA.GetBinaryHash()[haddrACaseStart:haddrACaseEnd]
		caseB = haddrB.GetBinaryHash()[haddrBCaseStart:haddrBCaseEnd]
	} else {
		return false, errors.New("Unrecognized club/case")
	}

	return bytes.Compare(caseA, caseB) == 0, nil

}

func (g *Geminus) BelongsInClub(club Club, addr string) (bool, error) {
	haddr := Addressing.NewAddress(addr, true)
	return g.HaveSameClub(club, g.Addr, haddr)
}

func (g *Geminus) SearchState(c Club, needle string) Addressing.Addr {
	var item Addressing.Addr

	club, err := g.GetClub(c)
	if err != nil {
		panic(err)
	}

	hneedle := Addressing.NewAddress(needle, true)

	for _, v := range club {
		if bytes.Compare(v.GetBinaryHash(), hneedle.GetBinaryHash()) == 0 {
			item = v
		}
	}

	return item
}

func (g *Geminus) Route(destination string) (Addressing.Addr, RoutingStatus) {
	var foundAddr Addressing.Addr
	var status RoutingStatus

	// a lot of functions like this hash the given string address
	// we should hash a given address once (from the arguments)
	// to avoid unnecessary resource waste

	// TODO: add numerical distance routing
	if belongs, _ := g.BelongsInClub(Hat, destination); belongs {
		foundAddr = g.SearchState(Hat, destination)
		status = HatRoute
	}

	if foundAddr == nil {
		haddr := Addressing.NewAddress(destination, true)
		bootClub, _ := g.GetClub(Boot)
		for _, baddr := range bootClub {
			haveSameHatClub, err := g.HaveSameClub(Hat, haddr, baddr)
			if err != nil {
				panic(err)
			}
			if haveSameHatClub {
				foundAddr = baddr
				status = BootForward
				break
			}
		}
	}

	if foundAddr == nil {
		hatClub, _ := g.GetClub(Hat)
		hatClubSize := len(hatClub)

		haddr := Addressing.NewAddress(destination, true)

		haveSameBootCase := true
		for hatClubSize > 0 && haveSameBootCase {
			foundAddr = hatClub[Tools.PickRandom(1, hatClubSize+1)]
			haveSameBootCase, _ = g.HaveSameClub(Boot, foundAddr, haddr)
		}

		status = RandomForward
	}

	if foundAddr == nil {
		status = Undefined
	}

	return foundAddr, status
}
