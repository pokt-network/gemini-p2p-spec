package main

import (
	CryptoRand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"sort"
	StrConversion "strconv"

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
		MyRing *Ring
		ID     *ID
		Clubs  map[string][]Node
		Cases  map[string]string
	}

	Network struct {
		Ring  *Ring
		Nodes []Node

		CaseLengths map[string]int
		Maps        map[string]map[string][]Node
	}

	Route struct {
		targetID      *ID
		destinationID *ID
		Status        string
		Hops          int
		Routed        bool
	}

	Stats struct {
		ClubCounts        map[string]int
		ClubSizesAverages map[string]int
		ClubsCoverage     map[string]float64
		UniqueClubsItems  map[string][]Node

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

func NewNode(id *ID, h, t, hcSize, tcSize int, ring *Ring) *Node {
	clubs := map[string][]Node{
		"head":  make([]Node, 0, hcSize),
		"tail":  make([]Node, 0, tcSize),
		"rtail": make([]Node, 0, hcSize),
		"rhead": make([]Node, 0, tcSize),
	}

	cases := map[string]string{
		"head": Case("head", id.IntRep, h),
		"tail": Case("tail", id.IntRep, t),
	}

	return &Node{
		MyRing: ring,
		ID:     id,
		Clubs:  clubs,
		Cases:  cases,
	}
}

func NewNetwork(headL, tailL, order, netSize int) *Network {
	ring := NewRing(order)

	caseLengths := map[string]int{
		"head": headL,
		"tail": tailL,
	}

	maps := map[string]map[string][]Node{
		"head":  make(map[string][]Node),
		"tail":  make(map[string][]Node),
		"rhead": make(map[string][]Node),
		"rtail": make(map[string][]Node),
	}

	return &Network{
		CaseLengths: caseLengths,
		Maps:        maps,
		Ring:        ring,
		Nodes:       make([]Node, 0, netSize),
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
	clubCounts := map[string]int{
		"head":  0,
		"tail":  0,
		"rhead": 0,
		"rtail": 0,
	}

	clubSizesAverage := map[string]int{
		"head":  0,
		"tail":  0,
		"rhead": 0,
		"rtail": 0,
	}

	clubsCoverage := map[string]float64{
		"head":  0,
		"tail":  0,
		"rhead": 0,
		"rtail": 0,
	}

	uniqueClubsItems := map[string][]Node{
		"head":  make([]Node, 0),
		"tail":  make([]Node, 0),
		"rhead": make([]Node, 0),
		"rtail": make([]Node, 0),
	}

	return &Stats{
		ClubCounts:        clubCounts,
		ClubSizesAverages: clubSizesAverage,
		ClubsCoverage:     clubsCoverage,
		UniqueClubsItems:  uniqueClubsItems,

		Routes: make([]Route, 0, 2000),
	}
}

func getDistance(ring, a, b big.Int) int64 {
	var distBA, rA, rB big.Int

	(&rA).Set(&a)
	(&rB).Set(&b)

	(&distBA).Sub(&rB, &rA)
	(&distBA).Mod(&distBA, &ring)

	return int64(math.Abs(float64(distBA.Int64())))
}

func Case(targetCase string, ID big.Int, length int) string {
	bin := GetBinary(ID)

	switch targetCase {
	case "head":
		return bin[:length]

	case "tail":
		idLength := len(bin)
		return bin[idLength-length:]

	default:
		return ""
	}
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

func populateNetwork(idPool []ID, network *Network, netSize, h, t, idLength int) {
	hcSize := netSize / int(math.Pow(2, float64(h)))
	tcSize := netSize / int(math.Pow(2, float64(t)))

	for i := 0; i < netSize; i++ {
		id := NewID(network.Ring)
		if isIDUnique(idPool, *id) {
			idPool = append(idPool, *id)
			node := NewNode(id, h, t, hcSize, tcSize, network.Ring)
			network.Nodes = append(network.Nodes, *node)
		}
	}
}

func seedNetwork(network *Network) {
	for i := 0; i < len(network.Nodes); i++ {
		iHead := network.Nodes[i].Cases["head"]
		iTail := network.Nodes[i].Cases["tail"]

		_, hExists := network.Maps["head"][iHead]
		if !hExists {
			network.Maps["head"][iHead] = []Node{network.Nodes[i]}
		} else {
			network.Maps["head"][iHead] = append(network.Maps["head"][iHead], network.Nodes[i])
		}

		_, tExists := network.Maps["tail"][iTail]
		if !tExists {
			network.Maps["tail"][iTail] = []Node{network.Nodes[i]}
		} else {
			network.Maps["tail"][iTail] = append(network.Maps["tail"][iTail], network.Nodes[i])
		}
	}

	for i := 0; i < len(network.Nodes); i++ {
		// this is reversed on purpose
		iHead := network.Nodes[i].Cases["tail"]
		iTail := network.Nodes[i].Cases["head"]

		_, hExists := network.Maps["rhead"][iHead]
		if !hExists {
			network.Maps["rhead"][iHead] = []Node{network.Nodes[i]}
		} else {
			network.Maps["rhead"][iHead] = append(network.Maps["rhead"][iHead], network.Nodes[i])
		}

		_, tExists := network.Maps["rtail"][iTail]
		if !tExists {
			network.Maps["rtail"][iTail] = []Node{network.Nodes[i]}
		} else {
			network.Maps["rtail"][iTail] = append(network.Maps["rtail"][iTail], network.Nodes[i])
		}
	}

	for i := 0; i < len(network.Nodes); i++ {
		iHead := network.Nodes[i].Cases["head"]
		iTail := network.Nodes[i].Cases["tail"]

		network.Nodes[i].Clubs["head"] = network.Maps["head"][iHead]
		network.Nodes[i].Clubs["tail"] = network.Maps["tail"][iTail]
	}

	for i := 0; i < len(network.Nodes); i++ {
		iHead := network.Nodes[i].Cases["tail"]
		iTail := network.Nodes[i].Cases["head"]

		network.Nodes[i].Clubs["rhead"] = network.Maps["rhead"][iHead]
		network.Nodes[i].Clubs["rtail"] = network.Maps["rtail"][iTail]
	}
}

func surveyNetwork(stats *Stats, network *Network) {
	stats.ClubCounts["head"] = len(network.Maps["head"])
	stats.ClubCounts["tail"] = len(network.Maps["tail"])

	stats.ClubCounts["rhead"] = len(network.Maps["rhead"])
	stats.ClubCounts["rtail"] = len(network.Maps["rtail"])

	allHeadClubsElements := make([]Node, 0, cap(network.Nodes)*2)
	allTailClubsElements := make([]Node, 0, cap(network.Nodes)*2)

	allrHeadClubsElements := make([]Node, 0, cap(network.Nodes)*2)
	allrTailClubsElements := make([]Node, 0, cap(network.Nodes)*2)

	uniqueHeadClubsElements := make([]Node, 0, cap(network.Nodes))
	uniqueTailClubsElements := make([]Node, 0, cap(network.Nodes))

	uniquerHeadClubsElements := make([]Node, 0, cap(network.Nodes))
	uniquerTailClubsElements := make([]Node, 0, cap(network.Nodes))

	for _, k := range network.Maps["head"] {
		stats.ClubSizesAverages["head"] = stats.ClubSizesAverages["head"] + len(k)
		allHeadClubsElements = append(allHeadClubsElements, k...)
	}

	for _, v := range allHeadClubsElements {
		if isNodeUnique(uniqueHeadClubsElements, v) {
			uniqueHeadClubsElements = append(uniqueHeadClubsElements, v)
		}
	}

	for _, k := range network.Maps["tail"] {
		stats.ClubSizesAverages["tail"] = stats.ClubSizesAverages["tail"] + len(k)
		allTailClubsElements = append(allTailClubsElements, k...)
	}

	for _, v := range allTailClubsElements {
		if isNodeUnique(uniqueTailClubsElements, v) {
			uniqueTailClubsElements = append(uniqueTailClubsElements, v)
		}
	}

	for _, k := range network.Maps["rhead"] {
		stats.ClubSizesAverages["rhead"] = stats.ClubSizesAverages["rhead"] + len(k)
		allrHeadClubsElements = append(allrHeadClubsElements, k...)
	}

	for _, v := range allrHeadClubsElements {
		if isNodeUnique(uniquerHeadClubsElements, v) {
			uniquerHeadClubsElements = append(uniquerHeadClubsElements, v)
		}
	}

	for _, k := range network.Maps["rtail"] {
		stats.ClubSizesAverages["rtail"] = stats.ClubSizesAverages["rtail"] + len(k)
		allrTailClubsElements = append(allrTailClubsElements, k...)
	}

	for _, v := range allrTailClubsElements {
		if isNodeUnique(uniquerTailClubsElements, v) {
			uniquerTailClubsElements = append(uniquerTailClubsElements, v)
		}
	}

	stats.ClubSizesAverages["head"] = stats.ClubSizesAverages["head"] / stats.ClubCounts["head"]
	stats.ClubSizesAverages["tail"] = stats.ClubSizesAverages["tail"] / stats.ClubCounts["tail"]

	stats.ClubsCoverage["head"] = float64(100 * (cap(network.Nodes) / len(uniqueHeadClubsElements)))
	stats.ClubsCoverage["tail"] = float64(100 * (cap(network.Nodes) / len(uniqueTailClubsElements)))

	stats.UniqueClubsItems["head"] = uniqueHeadClubsElements
	stats.UniqueClubsItems["tail"] = uniqueTailClubsElements

	stats.ClubSizesAverages["rhead"] = stats.ClubSizesAverages["rhead"] / stats.ClubCounts["rhead"]
	stats.ClubSizesAverages["rtail"] = stats.ClubSizesAverages["rtail"] / stats.ClubCounts["rtail"]

	stats.ClubsCoverage["rhead"] = float64(100 * (cap(network.Nodes) / len(uniquerHeadClubsElements)))
	stats.ClubsCoverage["rtail"] = float64(100 * (cap(network.Nodes) / len(uniquerTailClubsElements)))

	stats.UniqueClubsItems["rhead"] = uniquerHeadClubsElements
	stats.UniqueClubsItems["rtail"] = uniquerTailClubsElements
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
	if nodeA.Cases["head"] == nodeB.Cases["head"] {
		return "Head"
	}

	if nodeA.Cases["tail"] == nodeB.Cases["head"] {
		return "SecondHead"
	}

	for _, e := range network.Maps["tail"][nodeA.Cases["tail"]] {
		if e.Cases["head"] == nodeB.Cases["head"] {
			*found = e
			return "HeadInTails"
		}
	}

	for _, e := range network.Maps["rhead"][nodeA.Cases["head"]] {
		if e.Cases["head"] == nodeB.Cases["head"] {
			*found = e
			return "HeadInSecondTails"
		}
	}

	f, tries := false, 0
	for !f && tries < len(network.Maps["head"][nodeA.Cases["head"]]) {
		random := pickRandomNode(network.Maps["head"][nodeA.Cases["head"]])
		if random.Cases["tail"] != nodeB.Cases["tail"] {
			f = true
			*found = random
			return "TailInHeads"
		}
		tries++
	}

	f, tries = false, 0
	for !f && tries < len(network.Maps["rtail"][nodeA.Cases["tail"]]) {
		random := pickRandomNode(network.Maps["rtail"][nodeA.Cases["tail"]])
		if random.Cases["tail"] != nodeB.Cases["tail"] {
			f = true
			*found = random
			return "TailInSecondHeads"
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

func Router(target Node, destination Node, network *Network, h, t int) (Node, string) {
	found := Node{
		ID:     nil,
		MyRing: nil,
		Cases: map[string]string{
			"head": "",
			"tail": "",
		},
		Clubs: map[string][]Node{
			"head":  make([]Node, 0),
			"tail":  make([]Node, 0),
			"rtail": make([]Node, 0),
			"rhead": make([]Node, 0),
		},
	}
	haveSame := commonClub(target, destination, network, &found)

	fmt.Println("destination and target have the same", haveSame)

	switch haveSame {
	case "Head":
		fmt.Println("We are in head case, looking for the numerically closest node")

		distances := append(
			make([]Node, 0, len(network.Maps["head"][target.Cases["head"]])),
			network.Maps["head"][target.Cases["head"]]...,
		)

		sortByClosestDistance(destination, distances)

		if len(distances) > 0 {
			numericallyClosest := distances[0]
			fmt.Println("Success, found numericaly closest id in the heads club")
			return numericallyClosest, "Head"
		}
		fmt.Println("For some reason, no numerically closest node was picked")
		break

	case "SecondHead":
		fmt.Println("We are in head case (2), looking for the numerically closest node")
		distances := append(
			make([]Node, 0, len(network.Maps["rtail"][target.Cases["tail"]])),
			network.Maps["rtail"][target.Cases["tail"]]...,
		)

		sortByClosestDistance(destination, distances)

		if len(distances) > 0 {
			numericallyClosest := distances[0]
			fmt.Println("Success, found numericaly closest id in the heads (2) club")
			return numericallyClosest, "SecondHead"
		}

		fmt.Println("For some reason, no numerically closest node was picked (2)")
		break

	case "HeadInTails":
		fmt.Println("We are in tail case")

		if found.ID != nil {
			fmt.Println("Success, found a tails club item with same head case")
			return found, "HeadInTails"
		}

		fmt.Println("For some reason, no id from tails club with same head case  as destination was found")

		break

	case "HeadInSecondTails":
		fmt.Println("We are in tail case (2)")
		if found.ID != nil {
			fmt.Println("Success, found a tails club item with same head case (2)")
			return found, "HeadInSecondTails"
		}
		fmt.Println("For some reason, no id from tails club with same head case  as destination was found (2)")
		break

	case "TailInHeads":
		fmt.Println("We are in head case looking for a differnt tail")
		if found.ID != nil {
			fmt.Println("Success, found a tails club item with different head case")
			return found, "TailInHeads"
		}
		fmt.Println("For some reason, no id from heads club with different tail case as destination was found")
		break

	case "TailInSecondHeads":
		fmt.Println("We are in head case looking for a differnt tail (2)")
		if found.ID != nil {
			fmt.Println("Success, found a tails club item with different head case (2)")
			return found, "TailInSecondHeads"
		}
		fmt.Println("For some reason, no id from heads club with different tail case as destination was found (2)")
		break

	default:
		return Node{
			ID:     nil,
			MyRing: nil,
			Cases: map[string]string{
				"head": "",
				"tail": "",
			},
			Clubs: map[string][]Node{
				"head":  make([]Node, 0),
				"tail":  make([]Node, 0),
				"rhead": make([]Node, 0),
				"rtail": make([]Node, 0),
			},
		}, "Undefined"
	}

	return Node{
		ID:     nil,
		MyRing: nil,
		Cases: map[string]string{
			"head": "",
			"tail": "",
		},
		Clubs: map[string][]Node{
			"head":  make([]Node, 0),
			"tail":  make([]Node, 0),
			"rhead": make([]Node, 0),
			"rtail": make([]Node, 0),
		},
	}, "Undefined"
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

			for !route.Routed && route.Hops < 10000 {
				nextHop, status := Router(
					currentTarget,
					destinationNode,
					network,
					network.CaseLengths["head"],
					network.CaseLengths["tail"],
				)

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

func printStats(stats *Stats) {
	fmt.Println("Stats:")
	for k, v := range stats.ClubCounts {
		fmt.Println("*)", k, "clubs")
		fmt.Println("*---> Count:", v)
		fmt.Println("*---> Average (Actual) Clubs Size:", stats.ClubSizesAverages[k])
		fmt.Println("*---> Hat Clubs Network Coverage (%):", stats.ClubsCoverage[k])
		fmt.Println("*---> Number of Nodes Covered by Hat Clubs:", len(stats.UniqueClubsItems[k]))
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

func main() {
	if os.Args[1] == "" {
		panic("You have to supply a network size argument")
	}

	networkSize, _ := StrConversion.Atoi(os.Args[1])
	h, _ := StrConversion.Atoi(os.Args[2])
	t, _ := StrConversion.Atoi(os.Args[3])

	idLength := 128

	IDPool := make([]ID, 0, networkSize)
	network := NewNetwork(h, t, idLength, networkSize)
	statistics := NewStats()

	fmt.Println("Populating the network...")
	populateNetwork(IDPool, network, networkSize, h, t, idLength)

	fmt.Println("Seeding the network...")
	seedNetwork(network)

	fmt.Println("Surveying the network...")
	surveyNetwork(statistics, network)

	fmt.Println("Simulating the routing...")
	simulateRouting(statistics, network)

	printStats(statistics)
}
