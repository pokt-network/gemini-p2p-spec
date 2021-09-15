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
		HatCase  string
		BootCase string
	}

	Network struct {
		HatLength  int
		BootLength int
		Ring       *Ring
		Nodes      []Node
		HatMap     map[string][]Node
		BootMap    map[string][]Node
	}

	Route struct {
		targetID      *ID
		destinationID *ID
		Status        string
		Hops          int
		Routed        bool
	}

	Stats struct {
		HatLength           int
		BootLength          int
		HatClubsCount       int
		BootClubsCount      int
		AverageHatClubSize  int
		AverageBootClubSize int
		HatClubCoverage     float64
		BootClubCoverage    float64

		UniqueHatClubItems  []Node
		UniqueBootClubItems []Node

		LonelyIslands map[string]int

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
	return fmt.Sprintf("%0128b", &ID)
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

func NewNode(id *ID, h, b, hcSize, bcSize int, ring *Ring) *Node {
	return &Node{
		MyRing:   ring,
		ID:       id,
		HatClub:  make([]Node, 0, hcSize),
		BootClub: make([]Node, 0, bcSize),
		HatCase:  Hatcase(id.IntRep, h),
		BootCase: Bootcase(id.IntRep, b),
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
		Routed:        false,
	}
}

func NewStats() *Stats {
	return &Stats{
		HatLength:           0,
		BootLength:          0,
		HatClubsCount:       0,
		BootClubsCount:      0,
		AverageHatClubSize:  0,
		AverageBootClubSize: 0,
		HatClubCoverage:     0.0,
		BootClubCoverage:    0.0,
		LonelyIslands:       make(map[string]int),
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
	return GetBinary(ID)[:length]
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
			node := NewNode(id, h, b, hcSize, bcSize, network.Ring)
			network.Nodes = append(network.Nodes, *node)
		}
	}
}

func seedNetwork(network *Network) {
	hatMap := make(map[string][]Node)
	bootMap := make(map[string][]Node)

	for i := 0; i < len(network.Nodes); i++ {
		iHat := network.Nodes[i].HatCase
		iBoot := network.Nodes[i].BootCase

		_, exists := hatMap[iHat]
		if !exists {
			hatMap[iHat] = []Node{network.Nodes[i]}
		} else {
			hatMap[iHat] = append(hatMap[iHat], network.Nodes[i])
		}

		_, exists = bootMap[iBoot]
		if !exists {
			bootMap[iBoot] = []Node{network.Nodes[i]}
		} else {
			bootMap[iBoot] = append(bootMap[iBoot], network.Nodes[i])
		}
	}

	for i := 0; i < len(network.Nodes); i++ {
		iHat := network.Nodes[i].HatCase
		iBoot := network.Nodes[i].BootCase

		network.Nodes[i].HatClub = hatMap[iHat]
		network.Nodes[i].BootClub = bootMap[iBoot]
	}

	network.BootMap = bootMap
	network.HatMap = hatMap
}

func surveyNetwork(stats *Stats, network *Network) {
	stats.HatClubsCount = len(network.HatMap)
	stats.BootClubsCount = len(network.BootMap)

	allHatClubsElements := make([]Node, 0, cap(network.Nodes)*2)
	allBootClubsElements := make([]Node, 0, cap(network.Nodes)*2)

	uniqueHatClubsElements := make([]Node, 0, cap(network.Nodes))
	uniqueBootClubsElements := make([]Node, 0, cap(network.Nodes))

	for _, v := range network.HatMap {
		stats.AverageHatClubSize = stats.AverageHatClubSize + len(v)
		allHatClubsElements = append(allHatClubsElements, v...)
	}

	for _, v := range allHatClubsElements {
		if isNodeUnique(uniqueHatClubsElements, v) {
			uniqueHatClubsElements = append(uniqueHatClubsElements, v)
		}
	}

	for _, v := range network.BootMap {
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

	for _, bootClub := range network.BootMap {
		bootLen := len(bootClub)

		if bootLen == 1 {
			hatClub := Hatcase(bootClub[0].ID.IntRep, stats.HatLength)
			hatLen := len(network.HatMap[hatClub])
			if hatLen == 1 {
				stats.LonelyIslands["lonely-hat-boot"]++
			} else {
				stats.LonelyIslands["lonely-boot"]++
			}
		}
	}
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

func commonClub(nodeA, nodeB Node, network *Network, found *Node) string {
	if nodeA.HatCase == nodeB.HatCase {
		return "Hat"
	} else {
		fmt.Println("not Hat", nodeA.HatCase, nodeB.HatCase)
	}

	for _, e := range network.BootMap[nodeA.BootCase] {
		if e.HatCase == nodeB.HatCase {
			*found = e
			return "HatInBoot"
		} else {
			fmt.Println("not HatInBoot", e.HatCase, nodeB.HatCase)
		}
	}

	f, tries := false, 0
	for !f && tries < len(network.HatMap[nodeA.HatCase]) {
		random := pickRandomNode(network.HatMap[nodeA.HatCase])
		if random.BootCase != nodeB.BootCase {
			f = true
			*found = random
			return "ABootInHat"
		}
		tries++
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

func Router(target Node, destination Node, network *Network, h, b int) (Node, string) {
	found := Node{ID: nil, MyRing: nil, HatClub: nil, BootClub: nil}
	haveSame := commonClub(target, destination, network, &found)

	fmt.Println("destination and target have the same", haveSame)

	switch haveSame {
	case "Hat":
		fmt.Println("We are in hat case")
		distances := append(make([]Node, 0, len(network.HatMap[target.HatCase])), network.HatMap[target.HatCase]...)

		sortByClosestDistance(destination, distances)

		if len(distances) > 0 {
			numericallyClosest := distances[0]
			fmt.Println("Success, found numericaly closest id in the hat club")
			return numericallyClosest, "Hat"
		}
		fmt.Println("For some reason, no numerically closest node was picked")
		break

	case "HatInBoot":
		fmt.Println("We are in boot case")
		if found.ID != nil {
			fmt.Println("Success, found a boot club item with same hat case")
			return found, "InBoot"
		}
		fmt.Println("For some reason, no id from hat club with same boot club as destination was found")
		break

	case "ABootInHat":
		fmt.Println("We are in hat case looking for a differnt boot")
		if found.ID != nil {
			fmt.Println("Success, found a different boot club with same hat case")
			return found, "ABootInHat"
		}
		fmt.Println("For some reason, no id from hat club with different boot club as destination was found")
		break

	default:
		return Node{ID: nil, MyRing: nil, HatClub: nil, BootClub: nil}, "Undefined"
	}

	return Node{ID: nil, MyRing: nil, HatClub: nil, BootClub: nil}, "Undefined"
}

func simulateRouting(stats *Stats, network *Network) {
	targetNode := pickRandomNode(network.Nodes)
	destinationNodes := pickRandomNodes(network.Nodes, 2000)

	var currentTarget Node
	undefineds := 0
	for _, destinationNode := range destinationNodes {
		currentTarget = targetNode
		if destinationNode.ID.IntRep.Cmp(&currentTarget.ID.IntRep) != 0 {
			route := NewRoute(currentTarget.ID, destinationNode.ID)

			for !route.Routed && route.Hops < 200 {
				nextHop, status := Router(currentTarget, destinationNode, network, network.HatLength, network.BootLength)

				if status == "Undefined" || nextHop.ID == nil {
					route.Status = fmt.Sprintf("%s,%s", route.Status, status)
					route.Hops++
					fmt.Println(status, route.Hops)
					undefineds++
					break
				}

				if destinationNode.ID.IntRep.Cmp(&nextHop.ID.IntRep) == 0 {
					route.Routed = true
				} else {
					currentTarget = nextHop
				}

				route.Status = fmt.Sprintf("%s,%s", route.Status, status)
				route.Hops++
			}
			fmt.Println("Done with a route!")

			stats.Routes = append(stats.Routes, *route)
		}
	}
	fmt.Println("undefineds", undefineds)
}

func printStats(stats *Stats, details bool) {
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
	fmt.Println("*) Lonely Islands")
	if len(stats.LonelyIslands) > 0 {
		for k, v := range stats.LonelyIslands {
			fmt.Println("*--->", k, ":", v, "elements")
		}
	} else {
		fmt.Println("*---> None")
	}

	fmt.Printf("\n")
	fmt.Println("Routing stats:")
	fmt.Println(len(stats.Routes), "Requests for routing performed")

	hops := map[int]int{}
	for _, r := range stats.Routes {
		if r.Routed {
			hops[r.Hops]++
		}
	}

	for k, h := range hops {
		fmt.Println(h, "routes happened in ", k, "hops")
	}

	nhops := map[int]int{}
	for _, r := range stats.Routes {
		if !r.Routed {
			nhops[r.Hops]++
		}
	}

	for k, h := range nhops {
		fmt.Println(h, "routes did not route after", k, "hops")
	}
}

func simulateDistribution(networkSize, h, b int) *Stats {
	idLength := 128

	IDPool := make([]ID, 0, networkSize)
	network := NewNetwork(h, b, idLength, networkSize)
	statistics := NewStats()

	fmt.Println("Populating the network...")
	populateNetwork(IDPool, network, networkSize, h, b, idLength)

	fmt.Println("Seeding the network...")
	seedNetwork(network)

	fmt.Println("Surveying the network...")
	surveyNetwork(statistics, network)

	// fmt.Println("Simulating the routing...")
	// simulateRouting(statistics, network)

	return statistics
}

func main() {
	//networkSize, _ := StrConv.Atoi(os.Args[1])
	//h, _ := StrConv.Atoi(os.Args[2])
	//b, _ := StrConv.Atoi(os.Args[3])

	//averageFull := 0
	//averagePartially := 0

	//var stats *Stats
	//for i := 0; i < 10; i++ {
	//	stats = simulateDistribution(networkSize, h, b)
	//	if stats.LonelyIslands["lonely-hat-boot"] != 0 {
	//		averageFull += stats.LonelyIslands["lonely-hat-boot"]
	//	}

	//	if stats.LonelyIslands["lonely-boot"] != 0 {
	//		averagePartially += stats.LonelyIslands["lonely-boot"]
	//	}
	//}

	//printStats(stats, true)
	//fmt.Println("*) Average lonely islands stats on 10 runs:")
	//averageFull = averageFull / 10
	//averagePartially = averagePartially / 10

	//fmt.Println("*---> Average Fully Isolated Elements:", averageFull)
	//fmt.Println("*---> Average Partially Isolated Elements:", averagePartially)

	getRingPerspectives()
}

func getRingPerspectives() {
	ringMin := IntRing(0)
	ringMax := IntRing(128)
	// id := Id(&r)

	fmt.Println("Int representation:", &ringMin)
	fmt.Println("Binary representation:", GetBinary(ringMin))

	fmt.Println("Int representation:", &ringMax)
	fmt.Println("Binary representation:", GetBinary(ringMax))
}
