package main

import (
	"fmt"

	Gemini "gemelos/pkg/gemini"

	RandomData "github.com/Pallinder/go-randomdata"
)

const NetworkNodesCount = 6000
const AddressCaseLength = 3
const HatClubSize = 6000/2 ^ 3
const BootClubSize = 6000/2 ^ 3

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

func main() {
	peerCount := 0
	peerPool := make([]*Gemini.Geminus, 0, NetworkNodesCount)

	var peer *Gemini.Geminus
	fmt.Println("Generating peers.")
	for peerCount < NetworkNodesCount {
		peer = Gemini.NewGeminus(GetRandomIp(), NetworkNodesCount, 160, 3)
		if isPeerUnique(peerPool, peer) {
			peer.Init()
			peerPool = append(peerPool, peer)
			peerCount++
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

	fmt.Println("Stats")
	fmt.Println("For a random node who has been given 300 random address to route:")
	fmt.Printf("\n Hat Club Size :%v, Boot Club Size: %v\n", len(target.HatClub), len(target.BootClub))
	fmt.Printf("\n%v have been found in 1 hop (hat)\n", stats.hf)
	fmt.Printf("\n%v have been found in 1 hop (boot)\n", stats.bf)
	fmt.Printf("\n%v have been found in >2 hop (fwd)\n", stats.fwd)
	fmt.Printf("\n%v cause undefined behavior (cuz no wait until seed n stuff)\n", stats.undef)
	fmt.Println("Done")
}
