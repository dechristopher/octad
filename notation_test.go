package octad

import "testing"

// notationRoundtrip encodes every legal move in the position with the given
// notation and checks that decoding the result yields the same move.
func notationRoundtrip(t *testing.T, n Notation, ofen string) {
	t.Helper()
	g := gameFromOFEN(t, ofen)
	pos := g.Position()
	moves := pos.ValidMoves()
	if len(moves) == 0 {
		t.Fatalf("%s: no legal moves in %q", n, ofen)
	}
	for _, m := range moves {
		enc := n.Encode(pos, m)
		dec, err := n.Decode(pos, enc)
		if err != nil {
			t.Fatalf("%s: decode %q (from %s) in %q: %v", n, enc, m, ofen, err)
		}
		if dec.S1() != m.S1() || dec.S2() != m.S2() || dec.Promo() != m.Promo() {
			t.Fatalf("%s: roundtrip of %s in %q encoded %q decoded to %s%s%s",
				n, m, ofen, enc, dec.S1(), dec.S2(), dec.Promo())
		}
	}
}

func TestNotationStrings(t *testing.T) {
	if (UOINotation{}).String() != "UCI Notation" {
		t.Errorf("notation: UOI String() = %q", UOINotation{}.String())
	}
	if (AlgebraicNotation{}).String() != "Algebraic Notation" {
		t.Errorf("notation: Algebraic String() = %q", AlgebraicNotation{}.String())
	}
	if (LongAlgebraicNotation{}).String() != "Long Algebraic Notation" {
		t.Errorf("notation: LongAlgebraic String() = %q", LongAlgebraicNotation{}.String())
	}
}

func TestNotationRoundtrip(t *testing.T) {
	positions := []string{
		startOFEN,
		"3k/4/4/NK2 w N - 0 1",      // near castle available
		"3k/4/4/1KPP w CF - 0 1",    // center castle available, far blocked
		"3k/4/4/1K1P w NCF - 0 1",   // far castle available
		"3k/P3/4/K3 w - - 0 1",      // promotion available (a3 pawn)
		"3k/4/1Pp1/K3 w - c3 0 2",   // en passant capture available (b2xc3)
		"rk1n/4/2PP/QN1K w - - 0 1", // captures and piece variety
	}
	notations := []Notation{UOINotation{}, AlgebraicNotation{}, LongAlgebraicNotation{}}
	for _, n := range notations {
		for _, ofen := range positions {
			notationRoundtrip(t, n, ofen)
		}
	}
}

func TestUOIEncode(t *testing.T) {
	g := gameFromOFEN(t, startOFEN)
	pos := g.Position()
	want := map[string]string{ // UOI string -> itself (encode should reproduce S1S2)
		"a1c2": "a1c2", // knight
		"c1c2": "c1c2", // pawn
		"b1a1": "b1a1", // near castle
		"b1c1": "b1c1", // center castle
	}
	got := map[string]bool{}
	for _, m := range pos.ValidMoves() {
		got[UOINotation{}.Encode(pos, m)] = true
	}
	for uoi := range want {
		if !got[uoi] {
			t.Errorf("uoi: expected move %q to be encodable from start; got set %v", uoi, got)
		}
	}
}

func TestUOIDecodeErrorsAndNilPos(t *testing.T) {
	bad := []string{"a1", "a1c2c3", "z9z9", "a1c2x"}
	for _, s := range bad {
		if _, err := (UOINotation{}).Decode(nil, s); err == nil {
			t.Errorf("uoi: expected error decoding %q", s)
		}
	}
	// with a nil position a well-formed string decodes to a bare move (no tags)
	m, err := UOINotation{}.Decode(nil, "a1c2")
	if err != nil {
		t.Fatalf("uoi: decode with nil pos: %v", err)
	}
	if m.S1() != A1 || m.S2() != C2 {
		t.Fatalf("uoi: nil-pos decode = %s%s, want a1c2", m.S1(), m.S2())
	}
}

func TestAlgebraicCastleStrings(t *testing.T) {
	cases := []struct {
		ofen string
		uoi  string
		want string
	}{
		{"3k/4/4/NK2 w N - 0 1", "b1a1", "O"},        // near castle
		{"3k/4/4/1KP1 w C - 0 1", "b1c1", "O-O"},     // center castle
		{"3k/4/4/1K1P w NCF - 0 1", "b1d1", "O-O-O"}, // far castle
	}
	for _, tc := range cases {
		g := gameFromOFEN(t, tc.ofen)
		pos := g.Position()
		m := findMove(t, pos, tc.uoi)
		if got := (AlgebraicNotation{}).Encode(pos, m); got != tc.want {
			t.Errorf("algebraic: castle %s encoded %q, want %q", tc.uoi, got, tc.want)
		}
		if got := (LongAlgebraicNotation{}).Encode(pos, m); got != tc.want {
			t.Errorf("long algebraic: castle %s encoded %q, want %q", tc.uoi, got, tc.want)
		}
	}
}

func TestAlgebraicEncodeAnnotations(t *testing.T) {
	// promotion (king placed so the new queen gives no check)
	g := gameFromOFEN(t, "4/P2k/4/K3 w - - 0 1")
	promo := findMove(t, g.Position(), "a3a4q")
	if got := (AlgebraicNotation{}).Encode(g.Position(), promo); got != "a4=Q" {
		t.Errorf("algebraic: promotion encoded %q, want a4=Q", got)
	}

	// check (+): rook delivers check but not mate
	g = gameFromOFEN(t, "3k/4/4/RK2 w - - 0 1")
	check := findMove(t, g.Position(), "a1a4")
	if got := (AlgebraicNotation{}).Encode(g.Position(), check); got != "Ra4+" {
		t.Errorf("algebraic: check encoded %q, want Ra4+", got)
	}

	// checkmate (#)
	g = gameFromOFEN(t, "4/1K2/1Q2/3k w - - 7 7")
	mate := findMove(t, g.Position(), "b2c2")
	if got := (AlgebraicNotation{}).Encode(g.Position(), mate); got != "Qc2#" {
		t.Errorf("algebraic: mate encoded %q, want Qc2#", got)
	}
}

func TestAlgebraicDisambiguation(t *testing.T) {
	// two white rooks on a1 and c1 both reach b1 -> file disambiguation
	g := gameFromOFEN(t, "k2K/4/4/R1R1 w - - 0 1")
	m := findMove(t, g.Position(), "a1b1")
	if got := (AlgebraicNotation{}).Encode(g.Position(), m); got != "Rab1" {
		t.Errorf("algebraic: disambiguation encoded %q, want Rab1", got)
	}
}

func TestAlgebraicDecodeError(t *testing.T) {
	g := gameFromOFEN(t, startOFEN)
	if _, err := (AlgebraicNotation{}).Decode(g.Position(), "Zz9"); err == nil {
		t.Error("algebraic: expected error decoding nonsense move")
	}
	if _, err := (LongAlgebraicNotation{}).Decode(g.Position(), "q9q9"); err == nil {
		t.Error("long algebraic: expected error decoding nonsense move")
	}
}

// findMove returns the legal move in pos whose UOI string matches uoi.
func findMove(t *testing.T, pos *Position, uoi string) *Move {
	t.Helper()
	for _, m := range pos.ValidMoves() {
		if m.String() == uoi {
			return m
		}
	}
	t.Fatalf("notation: move %q not legal in %s", uoi, pos)
	return nil
}
