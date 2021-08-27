package main

import (
	"fmt"

	Gemini "gemelos/pkg/gemini"

	RandomData "github.com/Pallinder/go-randomdata"
)

type (
	Stats struct {
		HatClubs            map[string][]*Gemini.Geminus
		BootClubs           map[string][]*Gemini.Geminus
		HatClubsCount       int
		BootClubsCount      int
		AverageHatClubSize  int
		AverageBootClubSize int
		HatClubCoverage     float64
		BootClubCoverage    float64
	}
)

const NetworkNodesCount = 6000

func PrintStats(stats *Stats) {
	allHatClubsElements := make([]*Gemini.Geminus, 0, 6000*2)
	allBootClubsElements := make([]*Gemini.Geminus, 0, 6000*2)

	uniqueHatClubsElements := make([]*Gemini.Geminus, 0, 6000)
	uniqueBootClubsElements := make([]*Gemini.Geminus, 0, 6000)

	for _, v := range stats.HatClubs {
		stats.AverageHatClubSize = stats.AverageHatClubSize + len(v)
		allHatClubsElements = append(allHatClubsElements, v...)
	}

	for _, v := range allHatClubsElements {
		if isPeerUnique(uniqueHatClubsElements, v) {
			uniqueHatClubsElements = append(uniqueHatClubsElements, v)
		}
	}
	for _, v := range stats.BootClubs {
		stats.AverageBootClubSize = stats.AverageBootClubSize + len(v)
		allBootClubsElements = append(allBootClubsElements, v...)
	}
	for _, v := range allBootClubsElements {
		if isPeerUnique(uniqueBootClubsElements, v) {
			uniqueBootClubsElements = append(uniqueBootClubsElements, v)
		}
	}

	stats.AverageHatClubSize = stats.AverageHatClubSize / stats.HatClubsCount
	stats.AverageBootClubSize = stats.AverageBootClubSize / stats.BootClubsCount

	stats.HatClubCoverage = float64(100 * (NetworkNodesCount / len(uniqueHatClubsElements)))
	stats.BootClubCoverage = float64(100 * (NetworkNodesCount / len(uniqueBootClubsElements)))

	fmt.Println("Stats:")
	fmt.Println("*) HatClubs:")
	fmt.Println("*---> Count:", stats.HatClubsCount)
	fmt.Println("*---> Average (Actual) Clubs Size:", stats.AverageHatClubSize)
	fmt.Println("*---> Hat Clubs Network Coverage (%):", stats.HatClubCoverage)
	fmt.Println("*---> Number of Nodes Covered by Hat Clubs:", len(uniqueHatClubsElements))

	fmt.Printf("\n")
	fmt.Println("*) BootClubs:")
	fmt.Println("*---> Count:", stats.BootClubsCount)
	fmt.Println("*---> Average (Actual) Clubs Size:", stats.AverageBootClubSize)
	fmt.Println("*---> Boot Clubs Network Coverage (%):", stats.BootClubCoverage)
	fmt.Println("*---> Number of Nodes Covered by Boot Clubs:", len(uniqueBootClubsElements))
}

func GetStatsObj() *Stats {
	return &Stats{
		HatClubs:            make(map[string][]*Gemini.Geminus),
		BootClubs:           make(map[string][]*Gemini.Geminus),
		HatClubsCount:       0,
		BootClubsCount:      0,
		AverageHatClubSize:  0,
		AverageBootClubSize: 0,
		HatClubCoverage:     0.0,
		BootClubCoverage:    0.0,
	}
}

func GetRandomIp() string {
	return RandomData.IpV4Address()
}

func isPeerUnique(gs []*Gemini.Geminus, g *Gemini.Geminus) bool {
	for i := 0; i < len(gs); i++ {
		if gs[i].Addr.GetRaw() == g.Addr.GetRaw() {
			return false
		}
	}
	return true
}

func pickRandomPeer(pp []*Gemini.Geminus) *Gemini.Geminus {
	randomIndex := RandomData.Number(0, len(pp)-1)
	return pp[randomIndex]
}

func pickRandomPeers(pp []*Gemini.Geminus, count int) []*Gemini.Geminus {
	stack := make([]*Gemini.Geminus, 0, count)
	for i := 0; i < count; i++ {
		stack = append(stack, pickRandomPeer(pp))
	}
	return stack
}

func Categorize(stats *Stats, p *Gemini.Geminus) {
	addrLength := p.Addr.GetBinBitLength()
	addrBinHash := p.Addr.GetBinaryHash()

	fmt.Println(addrLength, addrBinHash, p.Params.AddrLength)
	hatcase := addrBinHash[0:p.Params.HatLength]
	bootcase := addrBinHash[addrLength-1-p.Params.BootLength : addrLength-1]

	if _, exists := stats.HatClubs[string(hatcase)]; exists {
		stats.HatClubs[string(hatcase)] = append(stats.HatClubs[string(hatcase)], p)
	} else {
		stats.HatClubs[string(hatcase)] = append(
			make([]*Gemini.Geminus, 0, p.Params.ClubSize.Hat),
			p,
		)
		stats.HatClubsCount++
	}

	if _, exists := stats.BootClubs[string(bootcase)]; exists {
		stats.BootClubs[string(bootcase)] = append(stats.BootClubs[string(bootcase)], p)
	} else {
		stats.BootClubs[string(bootcase)] = append(
			make([]*Gemini.Geminus, 0, p.Params.ClubSize.Boot),
			p,
		)
		stats.BootClubsCount++
	}
}

func main() {
	peerCount := 0
	peerPool := make([]*Gemini.Geminus, 0, NetworkNodesCount)
	clubStats := GetStatsObj()

	var peer *Gemini.Geminus
	fmt.Println("Generating peers.")
	for peerCount < NetworkNodesCount {
		peer = Gemini.NewGeminus(GetRandomIp(), NetworkNodesCount, 160, 5, 3)
		if isPeerUnique(peerPool, peer) {
			peer.Init()
			peerPool = append(peerPool, peer)
			peerCount++
			Categorize(clubStats, peer)
		}
		fmt.Println(peerCount)
	}

	target := pickRandomPeer(peerPool)
	fmt.Println("Seeding...")
	for i := 0; i < len(peerPool); i++ {
		if target.Addr.GetRaw() != peerPool[i].Addr.GetRaw() {
			target.SetState(peerPool[i].Addr.GetRaw())
		}
	}

	// destinations := pickRandomPeers(peerPool, 3000)

	stats := struct {
		hf    int
		bf    int
		fwd   int
		undef int
	}{hf: 0, bf: 0, fwd: 0, undef: 0}

	fmt.Println("Playing routing scenario...")
	for _, dest := range peerPool {
		fmt.Println("Asking peer:", target.Addr.GetRaw(), "for:", dest.Addr.GetRaw())
		_, status := target.Route(dest.Addr.GetRaw())
		fmt.Println(status)
		switch status {
		case Gemini.HatFind:
			stats.hf++
			break

		case Gemini.BootFind:
			stats.bf++
			break

		case Gemini.Forward:
			stats.fwd++
			break

		case Gemini.Undefined:
			stats.undef++
			break

		default:
			fmt.Println("RoutingStatus unrecognized")
			break
		}
	}

	PrintStats(clubStats)
	fmt.Println("Stats")
	fmt.Println("For a random node who has been given 6000 random address to route:")
	fmt.Printf("\n Hat Club Size :%v, Boot Club Size: %v\n", len(target.HatClub), len(target.BootClub))
	fmt.Printf("\n%v have been found in 1 hop (hat)\n", stats.hf)
	fmt.Printf("\n%v have been found in 1 hop (boot)\n", stats.bf)
	fmt.Printf("\n%v have been found in >2 hop (fwd)\n", stats.fwd)
	fmt.Printf("\n%v cause undefined behavior (cuz no wait until seed n stuff)\n", stats.undef)
	fmt.Println("Done")
}
