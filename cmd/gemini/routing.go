package main

import (
	CryptoRand "crypto/rand"
	"fmt"
	RandomData "github.com/Pallinder/go-randomdata"
	"math"
	"math/big"
	"math/rand"
	"sort"
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
		MyRing         *Ring
		ID             *ID
		HatClub        []Node
		BootClub       []Node
		HatClubString  string
		BootClubString string
	}

	Network struct {
		HatLength  int
		BootLength int
		Ring       *Ring
		Nodes      []Node
		Hatmap     map[string][]Node
		BootMap    map[string][]Node
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

func NewNode(id *ID, hcSize, bcSize int, ring *Ring, h, b int) *Node {
	return &Node{
		MyRing:         ring,
		ID:             id,
		HatClub:        make([]Node, 0, hcSize),
		BootClub:       make([]Node, 0, bcSize),
		HatClubString:  Hatcase(id.IntRep, h),
		BootClubString: Bootcase(id.IntRep, b),
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
	res := GetBinary(ID)
	return res[:length]
}

func Bootcase(ID big.Int, length int) string {
	idBin := GetBinary(ID)
	idLength := len(idBin)
	return idBin[idLength-length:]
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
			node := NewNode(id, hcSize, bcSize, network.Ring, h, b)
			network.Nodes = append(network.Nodes, *node)
		}
	}
}

func seedNetwork(network *Network) {

	hatmap := make(map[string][]Node)
	bootMap := make(map[string][]Node)

	//preload all boot/hat cases combinations

	for i := 0; i < len(network.Nodes); i++ {
		iHat := network.Nodes[i].HatClubString
		iBoot := network.Nodes[i].BootClubString

		_, exists := hatmap[iHat]
		if !exists {
			hatmap[iHat] = []Node{network.Nodes[i]}
		} else {
			hatmap[iHat] = append(hatmap[iHat], network.Nodes[i])
		}

		_, exists = bootMap[iBoot]
		if !exists {
			bootMap[iBoot] = []Node{network.Nodes[i]}
		} else {
			bootMap[iBoot] = append(bootMap[iBoot], network.Nodes[i])
		}
	}

	for i := 0; i < len(network.Nodes); i++ {
		iHat := network.Nodes[i].HatClubString
		iBoot := network.Nodes[i].BootClubString

		network.Nodes[i].HatClub = hatmap[iHat]
		network.Nodes[i].BootClub = bootMap[iBoot]
	}

	network.BootMap = bootMap
	network.Hatmap = hatmap
}

func surveyNetwork(stats *Stats, network *Network) {
	var emptybootcount = 0

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
		if len(v.BootClub) == 0 {
			emptybootcount++
		}
	}

	fmt.Println("Empty BOOTS", emptybootcount)

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

func commonClub(idA, idB Node, h, b int) string {
	aHat := idA.HatClubString
	bHat := idB.HatClubString

	//fmt.Println(aHat, bHat)

	if aHat == bHat {
		return "Hat"
	}

	aBoot := idA.BootClubString
	bBoot := idB.BootClubString

	//fmt.Println(aBoot, bBoot)
	if aBoot == bBoot {
		return "Boot"
	}

	for x, node := range idA.BootClub {
		if node.HatClubString == idB.HatClubString {
			fmt.Println("index for inboot", x)
			return "InBoot"
		}
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

func Router(target Node, destination Node, h, b int, route *Route, network *Network) (Node, string) {
	if target.ID.IntRep.Cmp(&destination.ID.IntRep) == 0 {
		return target, "hat"
	}

	route.Hops++

	haveSame := commonClub(target, destination, h, b)
	fmt.Println("destination and target have the same", haveSame)

	switch haveSame {
	case "Hat":
		fmt.Println("We are in hat case")

		for _, e := range network.Hatmap[target.HatClubString] {
			//fmt.Println(e.ID.IntRep.String(), destination.ID.IntRep.String())
			if e.ID.IntRep.Cmp(&destination.ID.IntRep) == 0 {
				return Router(e, destination, h, b, route, network)
			}
		}
		break
	case "Boot":
		fmt.Println("We are in boot case")

		for _, e := range network.BootMap[target.BootClubString] {
			//fmt.Println(e.ID.IntRep.String(), destination.ID.IntRep.String())
			if e.ID.IntRep.Cmp(&destination.ID.IntRep) == 0 {
				return Router(e, destination, h, b, route, network)
			}
		}
		break
	case "InBoot":
		fmt.Println("We are in hat in boot case scenario")
		for _, e := range network.BootMap[target.BootClubString] {
			if commonClub(e, destination, h, b) == "Hat" {
				fmt.Println("e", e.ID.IntRep.String())
				fmt.Println("Success, found a boot club item with same hat case")
				return Router(e, destination, h, b, route, network)
			}
		}
		fmt.Println("For some reason, no id from hat club with same boot club as destination was found")
		break
	case "Neither":
		var e Node
		for i := 0; i < len(network.Hatmap[target.HatClubString]); i++ {
			v := rand.Intn(len(network.Hatmap[target.HatClubString])-1) + 0
			if network.Hatmap[target.HatClubString][v].BootClubString != destination.BootClubString {
				e = network.Hatmap[target.HatClubString][v]
				fmt.Println("Success, found a random hat club item with a different boot club")
				return Router(e, destination, h, b, route, network)
			}
		}

		fmt.Println("For some reason, could not pick a random hat club item that has a different boot club", len(target.HatClub))
		break
	default:
		return Node{ID: nil, MyRing: nil, HatClub: nil, BootClub: nil}, "Undefined"
	}

	return Node{ID: nil, MyRing: nil, HatClub: nil, BootClub: nil}, "Nomatch"
}

func simulateRouting(stats *Stats, network *Network) {
	origin := pickRandomNode(network.Nodes)
	destinationNodes := pickRandomNodes(network.Nodes, 2000)

	var currentTarget Node
	for _, destinationNode := range destinationNodes {
		currentTarget = origin
		if destinationNode.ID.IntRep.Cmp(&currentTarget.ID.IntRep) != 0 {
			route := NewRoute(currentTarget.ID, destinationNode.ID)
			routed := false
			for !routed {
				nextHop, status := Router(currentTarget, destinationNode, network.HatLength, network.BootLength, route, network)

				fmt.Println("Next Hop=", nextHop.ID.IntRep.String())
				fmt.Println("Destination=", destinationNode.ID.IntRep.String())

				if nextHop.ID.IntRep.Cmp(&destinationNode.ID.IntRep) == 0 {
					fmt.Println("Path from ", origin.ID.IntRep.String(), "to", destinationNode.ID.IntRep.String(), "EXISTS")
					fmt.Println(status, nextHop.ID.IntRep.String(), route.Hops)
					stats.Routes = append(stats.Routes, *route)
					routed = true
				} else {
					currentTarget = nextHop
					route.Status = fmt.Sprintf("%s,%s", route.Status, status)
				}
			}

			fmt.Println("Done with a route!")
		}
	}
	fmt.Println("Amount of destinations", len(destinationNodes))
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
	// h = b = 5 seems to be the highest to work on 6k nodes
	networkSize := 10000
	h := 5
	b := 5

	idLength := 128

	IDPool := make([]ID, 0, networkSize)
	network := NewNetwork(h, b, idLength, networkSize)
	statistics := NewStats()

	fmt.Println("Populating the network...")
	populateNetwork(IDPool, network, networkSize, h, b, idLength)

	fmt.Println("Seed the network...")
	seedNetwork(network)

	fmt.Println("Surveying the network...")
	surveyNetwork(statistics, network)

	fmt.Println("Simulating the routing...")
	simulateRouting(statistics, network)

	printStats(statistics)

}
