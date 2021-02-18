package octad

import "testing"

var (
	validOFENs = []string{
		"ppkn/4/4/NKPP w NCFncf - 0 1",
		"ppkn/2P1/4/NK1P b NFncf c2 0 1",
		"p1kn/2P1/1p2/NK1P w NFnf b3 0 2",
		"p1kn/2P1/1pN1/1K1P b Fnf - 0 2",
		"3k/4/4/K2P w - - 0 1",
		"3k/4/4/NK1P w NF - 0 1",
		"3k/4/4/NKPP w NCF - 0 1",
		"4/2k1/4/1KP1 w C - 0 1",
		"4/3k/PK2/4 w - - 0 1",
		"rk1n/4/2PP/QN1K w - - 0 1",
		"4/1K2/1Q2/3k w - - 7 7",
		"1pkn/p3/K1N1/2PP b nc - 0 2",
		"1p1n/p1k1/K1N1/2PP w - - 0 3",
	}

	invalidOFENs = []string{
		"ppkn/4/4/NKP w NCFncf - 0 1",
		"ppkn/4/2P1/NK1P w NCFncf c5 0 1",
		"ppkn/4/2P1/NK1P w NCFncf - 0 -1",
		"ppkn/4/2P1/NK1P w NCFncf - -1 1",
		"ppkn/4/2P1/NK1P w NCFncf - 0 0",
		"ppkn/4/2P1/NK1P w NCFncf - - 0 1",
		"ppkn/4/2P1/NK1P w - 0 1",
		"ppkn/4/2P1/NK1P w e4 - 0 1",
	}
)

func TestValidOFENs(t *testing.T) {
	for _, f := range validOFENs {
		state, err := decodeOFEN(f)
		if err != nil {
			t.Fatal("ofen: received unexpected error", err)
		}
		if f != state.String() {
			t.Fatalf("ofen: expected board string %s but got %s", f, state.String())
		}
	}
}

func TestInvalidOFENs(t *testing.T) {
	for _, f := range invalidOFENs {
		if _, err := decodeOFEN(f); err == nil {
			t.Fatal("ofen: expected error from ", f)
		}
	}
}
