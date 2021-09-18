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

func NewNode(id *ID, h, b, t, hcSize, bcSize, tcSize int, ring *Ring) *Node {
	clubs := map[string][]Node{
		"head": make([]Node, 0, hcSize),
		"body": make([]Node, 0, bcSize),
		"tail": make([]Node, 0, tcSize),
	}

	cases := map[string]string{
		"head": Case("head", id.IntRep, h),
		"body": Case("body", id.IntRep, b),
		"tail": Case("tail", id.IntRep, t),
	}

	return &Node{
		MyRing: ring,
		ID:     id,
		Clubs:  clubs,
		Cases:  cases,
	}
}

func NewNetwork(headL, tailL, bodyL, order, netSize int) *Network {
	ring := NewRing(order)

	caseLengths := map[string]int{
		"head": headL,
		"body": bodyL,
		"tail": tailL,
	}

	maps := map[string]map[string][]Node{
		"head": make(map[string][]Node),
		"body": make(map[string][]Node),
		"tail": make(map[string][]Node),
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
		"head": 0,
		"tail": 0,
		"body": 0,
	}

	clubSizesAverage := map[string]int{
		"hat":  0,
		"tail": 0,
		"body": 0,
	}

	clubsCoverage := map[string]float64{
		"hat":  0,
		"tail": 0,
		"body": 0,
	}

	uniqueClubsItems := map[string][]Node{
		"head": make([]Node, 0),
		"body": make([]Node, 0),
		"tail": make([]Node, 0),
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

	case "body":
		addressLength := len(bin)
		midAddress := int((addressLength / 2) - 1)
		midCase := int(math.Floor(float64(length / 2)))
		fmt.Println("body", bin[midAddress-midCase:midAddress+midCase])
		return bin[midAddress-midCase : midAddress+midCase]

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

func populateNetwork(idPool []ID, network *Network, netSize, h, b, t, idLength int) {
	hcSize := netSize / int(math.Pow(2, float64(h)))
	bcSize := netSize / int(math.Pow(2, float64(b)))
	tcSize := netSize / int(math.Pow(2, float64(t)))

	for i := 0; i < netSize; i++ {
		id := NewID(network.Ring)
		if isIDUnique(idPool, *id) {
			idPool = append(idPool, *id)
			node := NewNode(id, h, t, b, hcSize, bcSize, tcSize, network.Ring)
			network.Nodes = append(network.Nodes, *node)
		}
	}
}

func seedNetwork(network *Network) {
	for i := 0; i < len(network.Nodes); i++ {
		iHead := network.Nodes[i].Cases["head"]
		iBody := network.Nodes[i].Cases["body"]
		iTail := network.Nodes[i].Cases["tail"]

		_, hExists := network.Maps["head"][iHead]
		if !hExists {
			network.Maps["head"][iHead] = []Node{network.Nodes[i]}
		} else {
			network.Maps["head"][iHead] = append(network.Maps["head"][iHead], network.Nodes[i])
		}

		_, bExists := network.Maps["body"][iBody]
		if !bExists {
			network.Maps["body"][iBody] = []Node{network.Nodes[i]}
		} else {
			network.Maps["body"][iBody] = append(network.Maps["body"][iBody], network.Nodes[i])
		}

		_, tExists := network.Maps["tail"][iTail]
		if !tExists {
			network.Maps["tail"][iTail] = []Node{network.Nodes[i]}
		} else {
			network.Maps["tail"][iTail] = append(network.Maps["tail"][iTail], network.Nodes[i])
		}
	}

	for i := 0; i < len(network.Nodes); i++ {
		iHead := network.Nodes[i].Cases["head"]
		iBody := network.Nodes[i].Cases["body"]
		iTail := network.Nodes[i].Cases["tail"]

		network.Nodes[i].Clubs["head"] = network.Maps["head"][iHead]
		network.Nodes[i].Clubs["body"] = network.Maps["body"][iBody]
		network.Nodes[i].Clubs["tail"] = network.Maps["tail"][iTail]
	}
}

func surveyNetwork(stats *Stats, network *Network) {
	stats.ClubCounts["head"] = len(network.Maps["head"])
	stats.ClubCounts["body"] = len(network.Maps["body"])
	stats.ClubCounts["tail"] = len(network.Maps["tail"])

	allHeadClubsElements := make([]Node, 0, cap(network.Nodes)*2)
	allBodyClubsElements := make([]Node, 0, cap(network.Nodes)*2)
	allTailClubsElements := make([]Node, 0, cap(network.Nodes)*2)

	uniqueHeadClubsElements := make([]Node, 0, cap(network.Nodes))
	uniqueBodyClubsElements := make([]Node, 0, cap(network.Nodes))
	uniqueTailClubsElements := make([]Node, 0, cap(network.Nodes))

	for _, k := range network.Maps["head"] {
		stats.ClubSizesAverages["head"] = stats.ClubSizesAverages["head"] + len(k)
		allHeadClubsElements = append(allHeadClubsElements, k...)
	}

	for _, v := range allHeadClubsElements {
		if isNodeUnique(uniqueHeadClubsElements, v) {
			uniqueHeadClubsElements = append(uniqueHeadClubsElements, v)
		}
	}

	for _, k := range network.Maps["body"] {
		stats.ClubSizesAverages["body"] = stats.ClubSizesAverages["body"] + len(k)
		allBodyClubsElements = append(allBodyClubsElements, k...)
	}

	for _, v := range allBodyClubsElements {
		if isNodeUnique(uniqueBodyClubsElements, v) {
			uniqueBodyClubsElements = append(uniqueBodyClubsElements, v)
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

	stats.ClubSizesAverages["head"] = stats.ClubSizesAverages["head"] / stats.ClubCounts["head"]
	stats.ClubSizesAverages["body"] = stats.ClubSizesAverages["body"] / stats.ClubCounts["body"]
	stats.ClubSizesAverages["tail"] = stats.ClubSizesAverages["tail"] / stats.ClubCounts["tail"]

	stats.ClubsCoverage["head"] = float64(100 * (cap(network.Nodes) / len(uniqueHeadClubsElements)))
	stats.ClubsCoverage["body"] = float64(100 * (cap(network.Nodes) / len(uniqueBodyClubsElements)))
	stats.ClubsCoverage["tail"] = float64(100 * (cap(network.Nodes) / len(uniqueTailClubsElements)))

	stats.UniqueClubsItems["head"] = uniqueHeadClubsElements
	stats.UniqueClubsItems["body"] = uniqueBodyClubsElements
	stats.UniqueClubsItems["tail"] = uniqueTailClubsElements
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

	for _, e := range network.Maps["tail"][nodeA.Cases["tail"]] {
		if e.Cases["head"] == nodeB.Cases["head"] {
			*found = e
			return "HeadInTails"
		}
	}

	for _, e := range network.Maps["body"][nodeA.Cases["body"]] {
		if e.Cases["head"] == nodeB.Cases["head"] {
			*found = e
			return "HeadInBodies"
		}
	}

	f, tries := false, 0
	for !f && tries < len(network.Maps["body"][nodeA.Cases["body"]]) {
		random := pickRandomNode(network.Maps["body"][nodeA.Cases["body"]])
		if random.Cases["tail"] != nodeB.Cases["tail"] {
			f = true
			*found = random
			return "ATailInBodies"
		}
		tries++
	}

	f, tries = false, 0
	for !f && tries < len(network.Maps["head"][nodeA.Cases["head"]]) {
		random := pickRandomNode(network.Maps["head"][nodeA.Cases["head"]])
		if random.Cases["tail"] != nodeB.Cases["tail"] {
			f = true
			*found = random
			return "ATailInHeads"
		}
		tries++
	}

	f, tries = false, 0
	for !f && tries < len(network.Maps["head"][nodeA.Cases["head"]]) {
		random := pickRandomNode(network.Maps["head"][nodeA.Cases["head"]])
		if random.Cases["body"] != nodeB.Cases["body"] {
			f = true
			*found = random
			return "ABodyInHeads"
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

func Router(target Node, destination Node, network *Network, h, t, b int) (Node, string) {
	found := Node{
		ID:     nil,
		MyRing: nil,
		Cases: map[string]string{
			"head": "",
			"body": "",
			"tail": "",
		},
		Clubs: map[string][]Node{
			"head": make([]Node, 0),
			"body": make([]Node, 0),
			"tail": make([]Node, 0),
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

	case "HeadInTails":
		fmt.Println("We are in tail case")
		if found.ID != nil {
			fmt.Println("Success, found a tails club item with same head case")
			return found, "HeadInTails"
		}
		fmt.Println("For some reason, no id from tails club with same head case  as destination was found")
		break

	case "HeadInBodies":
		fmt.Println("We are in body case")
		if found.ID != nil {
			fmt.Println("Success, found a bodies club item with same head case")
			return found, "HeadInBodies"
		}
		fmt.Println("For some reason, no id from bodies club with same head case as destination was found")
		break

	case "ATailInHeads":
		fmt.Println("We are in head case looking for a differnt tail")
		if found.ID != nil {
			fmt.Println("Success, found a tails club item with different head case")
			return found, "TailInHeads"
		}
		fmt.Println("For some reason, no id from heads club with different tail case as destination was found")
		break

	case "ATailInBodies":
		fmt.Println("We are in body case looking for a different tail")
		if found.ID != nil {
			fmt.Println("Success, found a bodies club item with different head case")
			return found, "TailInBodies"
		}
		fmt.Println("For some reason, no id from heads club with different tail case as destination was found")
		break

	case "ABodyInHeads":
		fmt.Println("We are in hat case looking for a differnt body")
		if found.ID != nil {
			fmt.Println("Success, found a heads club item with different body case")
			return found, "BodyInHeads"
		}
		fmt.Println("For some reason, no id from heads club with different body case as destination was found")
		break

	default:
		return Node{
			ID:     nil,
			MyRing: nil,
			Cases: map[string]string{
				"head": "",
				"body": "",
				"tail": "",
			},
			Clubs: map[string][]Node{
				"head": make([]Node, 0),
				"body": make([]Node, 0),
				"tail": make([]Node, 0),
			},
		}, "Undefined"
	}

	return Node{
		ID:     nil,
		MyRing: nil,
		Cases: map[string]string{
			"head": "",
			"body": "",
			"tail": "",
		},
		Clubs: map[string][]Node{
			"head": make([]Node, 0),
			"body": make([]Node, 0),
			"tail": make([]Node, 0),
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
					network.CaseLengths["body"],
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
	b, _ := StrConversion.Atoi(os.Args[3])
	t, _ := StrConversion.Atoi(os.Args[4])

	idLength := 128

	IDPool := make([]ID, 0, networkSize)
	network := NewNetwork(h, b, t, idLength, networkSize)
	statistics := NewStats()

	fmt.Println("Populating the network...")
	populateNetwork(IDPool, network, networkSize, h, b, t, idLength)

	fmt.Println("Seeding the network...")
	seedNetwork(network)

	fmt.Println("Surveying the network...")
	surveyNetwork(statistics, network)

	fmt.Println("Simulating the routing...")
	simulateRouting(statistics, network)

	printStats(statistics)
}
