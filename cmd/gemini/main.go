package main

import (
	"fmt"

	Addressing "gemelos/pkg/addressing"
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

func Categorize(stats *Stats, p *Gemini.Geminus) *Stats {
	addrLength := p.Addr.GetBinBitLength()
	addrBinHash := p.Addr.GetBinaryHash()

	fmt.Println(addrLength, addrBinHash, p.Params.AddrLength)
	hatcase := addrBinHash[0:p.Params.HatLength]
	bootcase := addrBinHash[addrLength-1-p.Params.BootLength : addrLength-1]

	if _, exists := stats.HatClubs[string(hatcase)]; exists {
		stats.HatClubs[string(hatcase)] = append(stats.HatClubs[string(hatcase)], p)
	} else {
		stats.HatClubs[string(hatcase)] = append(
			make([]*Gemini.Geminus, 0, p.Params.ClubSize[Gemini.Hat]),
			p,
		)
		stats.HatClubsCount++
	}

	if _, exists := stats.BootClubs[string(bootcase)]; exists {
		stats.BootClubs[string(bootcase)] = append(stats.BootClubs[string(bootcase)], p)
	} else {
		stats.BootClubs[string(bootcase)] = append(
			make([]*Gemini.Geminus, 0, p.Params.ClubSize[Gemini.Boot]),
			p,
		)
		stats.BootClubsCount++
	}

	return stats
}

func findInPeerPool(pp []*Gemini.Geminus, addr Addressing.Addr) *Gemini.Geminus {
	var peer *Gemini.Geminus = nil
	for _, p := range pp {
		if p.Addr.GetRaw() == addr.GetRaw() {
			peer = p
		}
	}
	return peer
}

func main() {
	peerCount := 0
	peerPool := make([]*Gemini.Geminus, 0, NetworkNodesCount)
	clubStats := GetStatsObj()

	var peer *Gemini.Geminus
	fmt.Println("Generating peers.")
	for peerCount < NetworkNodesCount {
		gParams := Gemini.NewGeminiConfig(NetworkNodesCount, 160, 5, 3)
		peer = Gemini.NewGeminus(GetRandomIp(), gParams)
		if isPeerUnique(peerPool, peer) {
			peer.Init()
			peerPool = append(peerPool, peer)
			peerCount++
			Categorize(clubStats, peer)
		}
	}

	fmt.Println("Seeding...")
	for _, k := range clubStats.HatClubs {
		fmt.Printf(".")
		for i := 0; i < len(k); i++ {
			for j := 0; j < len(k); j++ {
				if i != j {
					k[i].SetState(k[j].Addr.GetRaw())
				}
			}
		}
	}

	for _, k := range clubStats.BootClubs {
		fmt.Printf(".")
		for i := 0; i < len(k); i++ {
			for j := 0; j < len(k); j++ {
				if i != j {
					k[i].SetState(k[j].Addr.GetRaw())
				}
			}
		}
	}

	target := pickRandomPeer(peerPool)
	//	for i := 0; i < len(peerPool); i++ {
	//		if target.Addr.GetRaw() != peerPool[i].Addr.GetRaw() {
	//			target.SetState(peerPool[i].Addr.GetRaw())
	//		}
	//	}

	destinations := pickRandomPeers(peerPool, 600)

	stats := make([]struct {
		rs        Gemini.RoutingStatus
		hopsCount int
	}, 0, 6000)

	fmt.Println("Playing routing scenario...")

	for _, dest := range destinations {
		foundAddr, status := target.Route(dest.Addr.GetRaw())

		switch status {
		case Gemini.HatRoute:
			stats = append(
				stats, struct {
					rs        Gemini.RoutingStatus
					hopsCount int
				}{
					rs:        Gemini.HatRoute,
					hopsCount: 1,
				})
			break

		case Gemini.BootForward:
			hopsCount := 1
			if foundAddr.GetRaw() == dest.Addr.GetRaw() {
				hopsCount = 2
			} else {
				for foundPeer := findInPeerPool(peerPool, foundAddr); foundPeer.Addr.GetRaw() != dest.Addr.GetRaw(); foundPeer = findInPeerPool(peerPool, foundAddr) {
					foundAddr, status = foundPeer.Route(dest.Addr.GetRaw())
					if hopsCount > 20 {
						break
					}
					hopsCount++
				}
			}

			stats = append(
				stats, struct {
					rs        Gemini.RoutingStatus
					hopsCount int
				}{
					rs:        Gemini.BootForward,
					hopsCount: hopsCount,
				})

			break

		case Gemini.RandomForward:
			//fmt.Println("Since it is a boot forward, they gonna have the same hat case")
			//fmt.Println("foundAddr hatcase", foundAddr.GetBinaryHash()[0:target.Params.HatLength-1])
			//fmt.Println("dest addr hatcase", dest.Addr.GetBinaryHash()[0:target.Params.HatLength-1])

			//foundPeer := findInPeerPool(peerPool, foundAddr)
			//for _, v := range foundPeer.Clubs[Gemini.Hat] {
			//	fmt.Println(v.GetRaw())
			//	if dest.Addr.GetRaw() == v.GetRaw() {
			//		// fmt.Println(dest.Addr.GetRaw(), v.GetRaw())
			//		fmt.Println("=====> found it!", dest.Addr.GetRaw(), v.GetRaw())
			//	}
			//}

			//d := foundPeer.SearchState(Gemini.Hat, dest.Addr.GetRaw())
			//if d != nil {
			//	fmt.Println("found it again!")
			//} else {
			//	fmt.Println("Right mate, riiiiight")
			//}
			//panic("ss")
			hopsCount := 1
			if foundAddr.GetRaw() == dest.Addr.GetRaw() {
				hopsCount = 2
			} else {
				for foundPeer := findInPeerPool(peerPool, foundAddr); foundPeer.Addr.GetRaw() != dest.Addr.GetRaw(); foundPeer = findInPeerPool(peerPool, foundAddr) {
					foundAddr, status = foundPeer.Route(dest.Addr.GetRaw())
					if hopsCount > 20 {
						break
					}
					hopsCount++
				}
			}

			stats = append(
				stats, struct {
					rs        Gemini.RoutingStatus
					hopsCount int
				}{
					rs:        Gemini.RandomForward,
					hopsCount: hopsCount,
				})
			break

		case Gemini.Undefined:
			stats = append(
				stats, struct {
					rs        Gemini.RoutingStatus
					hopsCount int
				}{
					rs:        Gemini.Undefined,
					hopsCount: 0,
				})
			break

		default:
			fmt.Println("RoutingStatus unrecognized")
			break
		}
	}

	PrintStats(clubStats)
	fmt.Println("Stats")
	fmt.Println("For a random node who has been given 500 random address to route:")
	fmt.Printf("\n Hat Club Size :%v, Boot Club Size: %v\n", len(target.Clubs[Gemini.Hat]), len(target.Clubs[Gemini.Boot]))

	routingStatsSummary := struct {
		hf    int
		bf    []int
		rf    []int
		undef int
	}{
		hf:    0,
		bf:    make([]int, 22, 22),
		rf:    make([]int, 22, 22),
		undef: 0,
	}

	for _, v := range stats {
		if v.rs == Gemini.HatRoute {
			routingStatsSummary.hf++
		} else if v.rs == Gemini.BootForward {
			routingStatsSummary.bf[v.hopsCount]++
		} else if v.rs == Gemini.RandomForward {
			routingStatsSummary.rf[v.hopsCount]++
		} else if v.rs == Gemini.Undefined {
			routingStatsSummary.undef++
		}
	}

	fmt.Printf("\n%v have been found in 1 hop (hat)\n", routingStatsSummary.hf)
	for i, v := range routingStatsSummary.bf {
		if v != 0 {
			fmt.Printf("\n%v have been found in %v hop(s) (boot)\n", v, i)
		}
	}
	for i, v := range routingStatsSummary.rf {
		if v != 0 {
			fmt.Printf("\n%v have been found in %v hop(s) (random-boot)\n", v, i)
		}
	}
	fmt.Printf("\n%v cause undefined behavior (cuz no wait until seed n stuff)\n", routingStatsSummary.undef)
	fmt.Println("Done")
}
