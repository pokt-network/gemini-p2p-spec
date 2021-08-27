package gemini

import (
	bytes "bytes"
	"testing"
)

func TestNewGeminus(t *testing.T) {
	addr := "10.10.210.21"
	g := NewGeminus(addr, 6000, 160, 3)

	err := g.Init()

	if err != nil {
		t.Log("Faulty Gemini Address Length Param or Hash Function", err)
		t.Fail()
	}

	if g.Params.AddrLength != 160 {
		t.Log("Faulty Geminus Address")
		t.Fail()
	}

	if bytes.Compare([]byte(g.Addr.GetRaw()), []byte(addr)) != 0 {
		t.Log("Faulty Geminus Address")
		t.Fail()
	}

	if cap(g.HatClub) != g.Params.ClubSize.Boot {
		t.Log("Faulty Geminus HatClub Instantiation")
		t.Fail()
	}

	if cap(g.BootClub) != g.Params.ClubSize.Boot {
		t.Log("Faulty Geminus BootClub Instantiation")
		t.Fail()
	}
}

func TestSetState(t *testing.T) {
	addr := "10.10.210.21"
	g := NewGeminus(addr, 6000, 160, 3)

	g.Init()

	addressMap := map[string]Case{
		"10.10.230.331":   "",
		"109.20.212.121":  "",
		"100.130.322.222": "",
		"100.130.422.242": "",
		"41.210.412.312":  "",
	}

	for k, _ := range addressMap {
		addressMap[k] = g.SetState(k)
	}

	for k, v := range addressMap {
		foundAddr, status := g.Route(k)
		if foundAddr.GetRaw() != k && status != Forward {
			t.Log("Found address is not we are trying to route to")
			t.Fail()
		}
		if v == Hat && foundAddr.GetRaw() == k && status != HatFind {
			t.Log("Address belongs to Hat Club but RoutingSatus is not HatFind")
			t.Fail()
		} else if v == Boot && foundAddr.GetRaw() == k && status != BootFind {
			t.Log("Address belongs to Boot Club but RoutingStatus is not BootFind")
			t.Fail()
		} else if v == Foreign && status != Forward {
			t.Log("Address belongs to no club but RoutingStatus it not Forward")
			t.Fail()
		}
	}
}

func TestGetState(t *testing.T) {
	t.Log("No test case for Gemini.GetState")
}

func TestGetHatClub(t *testing.T) {
	t.Log("No test case for Gemini.GetHatClub")
}

func TestGetBootClub(t *testing.T) {
	t.Log("No test case for Gemini.GetBootClub")
}

func TestIsInHatClub(t *testing.T) {
	t.Log("No test case for Gemini.IsInHatClub")
}

func TestIsInBootClub(t *testing.T) {
	t.Log("No test case for Gemini.IsInBootClub")
}

func TestGetAddrDistance(t *testing.T) {
	t.Log("No test case for Gemini.GetAddrDistance")
}

func TestRoute(t *testing.T) {
	t.Log("No test case for Gemini.Route")
}
