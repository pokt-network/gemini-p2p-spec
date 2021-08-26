package gemini

import (
	bytes "bytes"
	"testing"
)

func TestNewGeminus(t *testing.T) {
	addr := "10.10.210.21"
	g := NewGeminus(addr)

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

	if len(g.HatClub) != g.Params.ClubSize.Boot {
		t.Log("Faulty Geminus HatClub Instantiation")
		t.Fail()
	}

	if len(g.BootClub) != g.Params.ClubSize.Boot {
		t.Log("Faulty Geminus BootClub Instantiation")
		t.Fail()
	}
}

func TestSetState(t *testing.T) {
	addr := "10.10.210.21"
	g := NewGeminus(addr)

	g.Init()

	g.SetState("10.10.210.21")
	g.SetState("10.10.230.331")
	g.SetState("109.20.212.121")
	g.SetState("100.130.322.222")
	g.SetState("100.130.322.222")
	g.SetState("41.210.412.312")

	//for _, v := range g.GetState() {
	//	t.Logf("State %s", v)
	//}

	t.Log("No test case for Gemini.SetState")
	t.Fail()
}

func TestGetState(t *testing.T) {
	t.Log("No test case for Gemini.GetState")
	t.Fail()
}

func TestGetHatClub(t *testing.T) {
	t.Log("No test case for Gemini.GetHatClub")
	t.Fail()
}

func TestGetBootClub(t *testing.T) {
	t.Log("No test case for Gemini.GetBootClub")
	t.Fail()
}

func TestIsInHatClub(t *testing.T) {
	t.Log("No test case for Gemini.IsInHatClub")
	t.Fail()
}

func TestIsInBootClub(t *testing.T) {
	t.Log("No test case for Gemini.IsInBootClub")
	t.Fail()
}

func TestGetAddrDistance(t *testing.T) {
	t.Log("No test case for Gemini.GetAddrDistance")
	t.Fail()
}

func TestRoute(t *testing.T) {
	t.Log("No test case for Gemini.Route")
	t.Fail()
}
