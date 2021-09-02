package main

import (
	CryptoRand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"sort"

	RandomData "github.com/Pallinder/go-randomdata"
)

type (
	ID struct {
		IntRep big.Int
		BinRep string
	}

	Ring struct {
		Order  int
		IntRep big.Int
	}

	Node struct {
		MyRing   *Ring
		ID       *ID
		HatClub  []Node
		BootClub []Node
	}

	Network struct {
		HatLength  int
		BootLength int
		Ring       *Ring
		Nodes      []Node
	}

	Route struct {
		targetID      *ID
		destinationID *ID
		Status        string
		Hops          int
	}

	Stats struct {
		HatClubs            map[string][]Node
		BootClubs           map[string][]Node
		HatClubsCount       int
		BootClubsCount      int
		AverageHatClubSize  int
		AverageBootClubSize int
		HatClubCoverage     float64
		BootClubCoverage    float64

		UniqueHatClubItems  []Node
		UniqueBootClubItems []Node

		Routes []Route
	}
)

func Id(ring *big.Int) *big.Int {
	var ID big.Int
	randInt, err := CryptoRand.Int(CryptoRand.Reader, ring)
	if err != nil {
		panic(err)
	}
	(&ID).Set(randInt)
	return &ID
}

func GetBinary(ID big.Int) string {
	return fmt.Sprintf("%b", &ID)
}

func IntRing(order int) big.Int {
	var ring big.Int
	ring.Exp(big.NewInt(2), big.NewInt(int64(order)), nil)

	return ring
}

func NewID(r *Ring) *ID {
	intID := Id(&(r.IntRep))
	binRep := GetBinary(*intID)
	return &ID{IntRep: *intID, BinRep: binRep}
}

func NewRing(order int) *Ring {
	ring := IntRing(order)
	return &Ring{
		Order:  order,
		IntRep: ring,
	}
}

func NewNode(id *ID, hcSize, bcSize int, ring *Ring) *Node {
	return &Node{
		MyRing:   ring,
		ID:       id,
		HatClub:  make([]Node, 0, hcSize),
		BootClub: make([]Node, 0, bcSize),
	}
}

func NewNetwork(hatl, bootl, order, netSize int) *Network {
	ring := NewRing(order)

	return &Network{
		HatLength:  hatl,
		BootLength: bootl,
		Ring:       ring,
		Nodes:      make([]Node, 0, netSize),
	}
}

func NewRoute(tID, destID *ID) *Route {
	return &Route{
		targetID:      tID,
		destinationID: destID,
		Status:        "",
		Hops:          0,
	}
}

func NewStats() *Stats {
	return &Stats{
		HatClubs:            make(map[string][]Node),
		BootClubs:           make(map[string][]Node),
		HatClubsCount:       0,
		BootClubsCount:      0,
		AverageHatClubSize:  0,
		AverageBootClubSize: 0,
		HatClubCoverage:     0.0,
		BootClubCoverage:    0.0,
		Routes:              make([]Route, 0, 2000),
	}
}

func getDistance(ring, a, b big.Int) int64 {
	var distBA, rA, rB big.Int

	(&rA).Set(&a)
	(&rB).Set(&b)

	(&distBA).Sub(&rB, &rA)
	(&distBA).Mod(&distBA, &ring)

	//buf := bytes.NewReader(distBA.Bytes())
	//binary.Read(buf, binary.BigEndian, &distance)

	return int64(math.Abs(float64(distBA.Int64())))
}

func Hatcase(ID big.Int, length int) string {
	return GetBinary(ID)[0 : length-1]
}

func Bootcase(ID big.Int, length int) string {
	idBin := GetBinary(ID)
	idLength := len(idBin)

	return idBin[idLength-1-length : idLength-1]
}

func isIDUnique(pool []ID, id ID) bool {
	unique := true
	for _, v := range pool {
		if v.IntRep.Cmp(&id.IntRep) == 0 {
			unique = false
			break
		}
	}
	return unique
}

func isNodeUnique(nodes []Node, node Node) bool {
	unique := true
	for _, v := range nodes {
		if v.ID.IntRep.Cmp(&(node.ID.IntRep)) == 0 {
			unique = false
			break
		}
	}
	return unique
}

func populateNetwork(idPool []ID, network *Network, netSize, h, b, idLength int) {
	hcSize := netSize / int(math.Pow(2, float64(h)))
	bcSize := netSize / int(math.Pow(2, float64(b)))

	for i := 0; i < netSize; i++ {
		id := NewID(network.Ring)
		if isIDUnique(idPool, *id) {
			idPool = append(idPool, *id)
			node := NewNode(id, hcSize, bcSize, network.Ring)
			network.Nodes = append(network.Nodes, *node)
		}
	}
}

func seedNetwork(network *Network) {
	for i := 0; i < len(network.Nodes); i++ {
		for j := 0; j < len(network.Nodes); j++ {
			iHat := Hatcase(network.Nodes[i].ID.IntRep, network.HatLength)
			jHat := Hatcase(network.Nodes[j].ID.IntRep, network.HatLength)

			iBoot := Bootcase(network.Nodes[i].ID.IntRep, network.BootLength)
			jBoot := Bootcase(network.Nodes[j].ID.IntRep, network.BootLength)

			if iHat == jHat {
				network.Nodes[i].HatClub = append(network.Nodes[i].HatClub, network.Nodes[j])
				network.Nodes[j].HatClub = append(network.Nodes[j].HatClub, network.Nodes[i])
			}

			if iBoot == jBoot {
				network.Nodes[i].BootClub = append(network.Nodes[i].BootClub, network.Nodes[j])
				network.Nodes[j].BootClub = append(network.Nodes[j].BootClub, network.Nodes[i])
			}
		}
	}
}

func surveyNetwork(stats *Stats, network *Network) {
	for _, v := range network.Nodes {
		vHat := Hatcase(v.ID.IntRep, network.HatLength)
		vBoot := Bootcase(v.ID.IntRep, network.BootLength)

		if _, exists := stats.HatClubs[string(vHat)]; exists {
			stats.HatClubs[string(vHat)] = append(stats.HatClubs[string(vHat)], v)
		} else {
			stats.HatClubs[string(vHat)] = append(
				make([]Node, 0, int(float64(cap(network.Nodes))/math.Pow(2, float64(network.HatLength)))),
				v,
			)
			stats.HatClubsCount++
		}

		if _, exists := stats.BootClubs[string(vBoot)]; exists {
			stats.BootClubs[string(vBoot)] = append(stats.BootClubs[string(vBoot)], v)
		} else {
			stats.BootClubs[string(vBoot)] = append(
				make([]Node, 0, int(float64(cap(network.Nodes))/math.Pow(2, float64(network.BootLength)))),
				v,
			)
			stats.BootClubsCount++
		}
	}

	allHatClubsElements := make([]Node, 0, cap(network.Nodes)*2)
	allBootClubsElements := make([]Node, 0, cap(network.Nodes)*2)

	uniqueHatClubsElements := make([]Node, 0, cap(network.Nodes))
	uniqueBootClubsElements := make([]Node, 0, cap(network.Nodes))

	for _, v := range stats.HatClubs {
		stats.AverageHatClubSize = stats.AverageHatClubSize + len(v)
		allHatClubsElements = append(allHatClubsElements, v...)
	}

	for _, v := range allHatClubsElements {
		if isNodeUnique(uniqueHatClubsElements, v) {
			uniqueHatClubsElements = append(uniqueHatClubsElements, v)
		}
	}

	for _, v := range stats.BootClubs {
		stats.AverageBootClubSize = stats.AverageBootClubSize + len(v)
		allBootClubsElements = append(allBootClubsElements, v...)
	}

	for _, v := range allBootClubsElements {
		if isNodeUnique(uniqueBootClubsElements, v) {
			uniqueBootClubsElements = append(uniqueBootClubsElements, v)
		}
	}

	stats.AverageHatClubSize = stats.AverageHatClubSize / stats.HatClubsCount
	stats.AverageBootClubSize = stats.AverageBootClubSize / stats.BootClubsCount

	stats.HatClubCoverage = float64(100 * (cap(network.Nodes) / len(uniqueHatClubsElements)))
	stats.BootClubCoverage = float64(100 * (cap(network.Nodes) / len(uniqueBootClubsElements)))

	stats.UniqueHatClubItems = uniqueHatClubsElements
	stats.UniqueBootClubItems = uniqueBootClubsElements

}

func pickRandomNode(nodePool []Node) Node {
	randomIndex := RandomData.Number(0, len(nodePool)-1)
	return nodePool[randomIndex]
}

func pickRandomNodes(nodePool []Node, count int) []Node {
	stack := make([]Node, 0, count)
	for i := 0; i < count; i++ {
		node := pickRandomNode(nodePool)
		if isNodeUnique(stack, node) {
			stack = append(stack, node)
		}
	}
	return stack
}

func commonClub(idA ID, idB ID, h, b int) string {
	aHat := Hatcase(idA.IntRep, h)
	bHat := Hatcase(idB.IntRep, h)

	fmt.Println(aHat, bHat)

	if aHat == bHat {
		return "Hat"
	}

	aBoot := Bootcase(idA.IntRep, b)
	bBoot := Bootcase(idB.IntRep, b)

	fmt.Println(aBoot, bBoot)
	if aBoot == bBoot {
		return "Boot"
	}

	return "Neither"
}

func sortByClosestDistance(ref Node, list []Node) {
	sort.Slice(list, func(i, j int) bool {
		distA := getDistance(ref.MyRing.IntRep, ref.ID.IntRep, list[i].ID.IntRep)
		distB := getDistance(ref.MyRing.IntRep, ref.ID.IntRep, list[j].ID.IntRep)

		if distA < distB {
			return true
		} else {
			return false
		}
	})
}

func Router(target Node, destination Node, h, b int) (Node, string) {
	haveSame := commonClub(*target.ID, *destination.ID, h, b)

	fmt.Println("destination and target have the same", haveSame)

	switch haveSame {
	case "Hat":
		fmt.Println("We are in hat case")
		distances := append(make([]Node, 0, len(target.HatClub)), target.HatClub...)

		sortByClosestDistance(destination, distances)

		if len(distances) > 0 {
			numericallyClosest := distances[0]
			fmt.Println("Success, found numericaly closest id in the hat club")
			return numericallyClosest, "Hat"
		}
		fmt.Println("For some reason, no numerically closest node was picked")

	case "Boot":
		fmt.Println("We are in boot case")
		for _, e := range target.BootClub {
			if commonClub(*e.ID, *destination.ID, h, b) == "Hat" {
				fmt.Println("e", *e.ID)
				fmt.Println("Success, found a boot club item with same hat case")
				return e, "Boot"
			}
		}
		fmt.Println("For some reason, no id from hat club with same boot club as destination was found")

	case "Neither":
		var e Node
		for i := 0; i < len(target.HatClub); i++ {
			if Bootcase(target.HatClub[i].ID.IntRep, b) != Bootcase(destination.ID.IntRep, b) {
				e = target.HatClub[i]
				fmt.Println("========")
				fmt.Println(e.ID)
				fmt.Println("========")
				break
			}
		}

		if e.ID != nil {
			fmt.Println("Success, found a random hat club item with a different boot club")
			return e, "Random"
		}
		fmt.Println("For some reason, could not pick a random hat club item that has a different boot club")

	default:
		return Node{ID: nil, MyRing: nil, HatClub: nil, BootClub: nil}, "Undefined"
	}

	return Node{ID: nil, MyRing: nil, HatClub: nil, BootClub: nil}, "Undefined"
}

func simulateRouting(stats *Stats, network *Network) {
	targetNode := pickRandomNode(network.Nodes)
	destinationNodes := pickRandomNodes(network.Nodes, 600)

	var currentTarget Node
	for _, destinationNode := range destinationNodes {
		currentTarget = targetNode
		if destinationNode.ID.IntRep.Cmp(&currentTarget.ID.IntRep) != 0 {
			route := NewRoute(currentTarget.ID, destinationNode.ID)
			routed := false

			for !routed {
				nextHop, status := Router(currentTarget, destinationNode, network.HatLength, network.BootLength)

				fmt.Println("Next Hop=", nextHop.ID.IntRep)
				fmt.Println("Destination=", destinationNode.ID.IntRep)

				if nextHop.ID.IntRep.Cmp(&destinationNode.ID.IntRep) == 0 {
					routed = true
				} else {
					currentTarget = nextHop
				}

				route.Status = fmt.Sprintf("%s,%s", route.Status, status)
				route.Hops++
				fmt.Println(status, *nextHop.ID, route.Hops)
			}
			fmt.Println("Done with a route!")

			stats.Routes = append(stats.Routes, *route)
		}
	}
}

func printStats(stats *Stats) {
	fmt.Println("Stats:")
	fmt.Println("*) HatClubs:")
	fmt.Println("*---> Count:", stats.HatClubsCount)
	fmt.Println("*---> Average (Actual) Clubs Size:", stats.AverageHatClubSize)
	fmt.Println("*---> Hat Clubs Network Coverage (%):", stats.HatClubCoverage)
	fmt.Println("*---> Number of Nodes Covered by Hat Clubs:", len(stats.UniqueHatClubItems))

	fmt.Printf("\n")
	fmt.Println("*) BootClubs:")
	fmt.Println("*---> Count:", stats.BootClubsCount)
	fmt.Println("*---> Average (Actual) Clubs Size:", stats.AverageBootClubSize)
	fmt.Println("*---> Boot Clubs Network Coverage (%):", stats.BootClubCoverage)
	fmt.Println("*---> Number of Nodes Covered by Boot Clubs:", len(stats.UniqueBootClubItems))

	fmt.Printf("\n")
	fmt.Println("Routing stats:")
	fmt.Println(len(stats.Routes), "Requests for routing performed")

	hops := map[int]int{}
	for _, r := range stats.Routes {
		hops[r.Hops]++
	}

	for k, h := range hops {
		fmt.Println(h, "routes happened in ", k, "hops")
	}
}

func main() {
	networkSize := 6000
	h := 5
	b := 5

	idLength := 128

	IDPool := make([]ID, 0, networkSize)
	network := NewNetwork(h, b, idLength, networkSize)
	statistics := NewStats()

	fmt.Println("Populating the network...")
	populateNetwork(IDPool, network, networkSize, h, b, idLength)

	fmt.Println("Surveying the network...")
	seedNetwork(network)

	fmt.Println("Surveying the network...")
	surveyNetwork(statistics, network)

	fmt.Println("Simulating the routing...")
	// simulateRouting(statistics, network)

	printStats(statistics)
}
