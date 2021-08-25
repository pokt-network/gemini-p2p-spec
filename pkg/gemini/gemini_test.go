package gemini

import (
	bytes "bytes"
	"testing"
)

func TestNewGeminus(t *testing.T) {
	addr := make([]byte, 64)
	g := NewGeminus(addr)

	if g.Params.AddrLength != 64 {
		t.Log("Faulty Geminus Address")
		t.Fail()
	}

	if bytes.Compare(g.Addr, addr) != 0 {
		t.Log("Faulty Geminus Address")
		t.Fail()
	}

	if g.HatClub.Start != 0 {
		t.Log("Faulty Geminus HatClub Start Value")
		t.Fail()
	}

	if g.HatClub.End != g.Params.ClubSize.Hat-1 {
		t.Log("Faulty Geminus HatClub End Value")
		t.Fail()
	}

	if g.BootClub.Start != g.Params.ClubSize.Hat {
		t.Log("Faulty Geminus BootClub Start Value")
		t.Fail()
	}

	if g.BootClub.End != g.Params.ClubSize.Hat+g.Params.ClubSize.Boot-1 {
		t.Log("Faulty Geminus BootClub End Value")
		t.Fail()
	}
}

func TestGetState(t *testing.T) {
	t.Log("No test case for Gemini.GetState")
	t.Fail()
}

func TestSetState(t *testing.T) {
	t.Log("No test case for Gemini.SetState")
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
