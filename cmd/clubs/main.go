package main

import (
	"fmt"
	Addressing "gemelos/pkg/addressing"

	RandomData "github.com/Pallinder/go-randomdata"
)

const NetworkNodesCount = 6000
const AddressCaseLength = 3
const HatClubSize = 6000/2 ^ 3
const BootClubSize = 6000/2 ^ 3

type (
	Stats struct {
		HatClubs       map[string][]*Addressing.Address
		BootClubs      map[string][]*Addressing.Address
		HatClubsCount  int
		BootClubsCount int
	}
)

func GetRandomIp() string {
	return RandomData.IpV4Address()
}

func GetAddress() *Addressing.Address {
	return Addressing.NewAddress(GetRandomIp(), true)
}

func Categorize(stats *Stats, addr *Addressing.Address) {
	addrBinHash := addr.GetBinaryHash()

	hatcase := addrBinHash[0:AddressCaseLength]
	bootcase := addrBinHash[len(addrBinHash)-1-AddressCaseLength : len(addrBinHash)-1]

	if _, exists := stats.HatClubs[string(hatcase)]; exists {
		stats.HatClubs[string(hatcase)] = append(stats.HatClubs[string(hatcase)], addr)
	} else {
		stats.HatClubs[string(hatcase)] = append(
			make([]*Addressing.Address, 0, HatClubSize),
			addr,
		)
		stats.HatClubsCount++
	}

	if _, exists := stats.BootClubs[string(bootcase)]; exists {
		stats.BootClubs[string(bootcase)] = append(stats.BootClubs[string(bootcase)], addr)
	} else {
		stats.BootClubs[string(bootcase)] = append(
			make([]*Addressing.Address, 0, BootClubSize),
			addr,
		)
		stats.BootClubsCount++
	}
}

func PrintStats(stats *Stats) {
	fmt.Println("Stats:")
	fmt.Println("*) HatClubs:")
	fmt.Println("*---> Count:", stats.HatClubsCount)
	fmt.Println("*---> Clubs/Length/Values:")
	fmt.Println(" ")
	for k, v := range stats.HatClubs {
		fmt.Println("--------------------------------------")
		fmt.Printf("\n\nClub[%s]: Length=%d\n", k, len(v))
		fmt.Println("--------------------------------------")
		//for _, a := range v {
		//	fmt.Printf("\tValues: %v \t", a.Raw)
		//}
	}

	fmt.Printf("\n\n*) BootClubs:\n")
	fmt.Println("*---> Count:", stats.BootClubsCount)
	fmt.Println("*---> Clubs/Length:")
	fmt.Println(" ")
	for k, v := range stats.BootClubs {
		fmt.Println("--------------------------------------")
		fmt.Printf("\n\nClub[%s]: Length=%d\n", k, len(v))
		fmt.Println("--------------------------------------")
		//for _, a := range v {
		//	fmt.Printf("\tValues: %v \t", a.Raw)
		//}
	}
}

func GetStatsObj() *Stats {
	return &Stats{
		HatClubs:       make(map[string][]*Addressing.Address),
		BootClubs:      make(map[string][]*Addressing.Address),
		HatClubsCount:  0,
		BootClubsCount: 0,
	}
}

func isAddressUnique(as []*Addressing.Address, addr *Addressing.Address) bool {
	for i := 0; i < len(as); i++ {
		if as[i].Raw == addr.Raw {
			return false
		}
	}
	return true
}

func main() {
	stats := GetStatsObj()

	addressCount := 0
	addressPool := make([]*Addressing.Address, 0, 6000)

	for addressCount < NetworkNodesCount {
		addr := GetAddress()
		if isAddressUnique(addressPool, addr) {
			addressPool = append(addressPool, addr)
			addressCount++
			Categorize(stats, addr)
		}
	}

	fmt.Println("6000 adresses reached")

	PrintStats(stats)

	//fmt.Println("======= Addressess ========")
	//for _, v := range addressPool {
	//	fmt.Printf("%s\t", v.Raw)
	//}
	fmt.Println("Done")
}
