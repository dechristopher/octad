package octad

import (
	"sort"
	"strings"
	"testing"
)

// gameFromOFEN builds a game from an OFEN string, failing the test on error.
func gameFromOFEN(t *testing.T, ofen string) *Game {
	t.Helper()
	fromPos, err := OFEN(ofen)
	if err != nil {
		t.Fatalf("bad ofen %q: %v", ofen, err)
	}
	g, err := NewGame(fromPos)
	if err != nil {
		t.Fatalf("new game from %q: %v", ofen, err)
	}
	return g
}

// castleMoveSet returns the UOI strings of every castle move available in the
// position described by ofen.
func castleMoveSet(t *testing.T, ofen string) map[string]bool {
	t.Helper()
	g := gameFromOFEN(t, ofen)
	set := map[string]bool{}
	for _, m := range g.ValidMoves() {
		if m.castles() {
			set[m.String()] = true
		}
	}
	return set
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// applyUOI finds the legal move with the given UOI string and plays it.
func applyUOI(t *testing.T, g *Game, uoi string) {
	t.Helper()
	for _, m := range g.ValidMoves() {
		if m.String() == uoi {
			if err := g.Move(m); err != nil {
				t.Fatalf("apply %s: %v", uoi, err)
			}
			return
		}
	}
	var got []string
	for _, m := range g.ValidMoves() {
		got = append(got, m.String())
	}
	sort.Strings(got)
	t.Fatalf("move %s not legal; valid moves: %v", uoi, got)
}

func TestCastleAvailability(t *testing.T) {
	cases := []struct {
		name string
		ofen string
		want []string
	}{
		// legacy parity: standard start offers the near and center castles;
		// the far castle is blocked because the center piece sits between
		{"standard white", "ppkn/4/4/NKPP w NCFncf - 0 1", []string{"b1a1", "b1c1"}},
		{"standard black", "ppkn/4/4/NKPP b NCFncf - 0 1", []string{"c4b4", "c4d4"}},

		// near castle (adjacent swap) from each viable king square; the near
		// slot is always the closest square, toward the king's nearest edge
		{"near king b1", "3k/4/4/NK2 w N - 0 1", []string{"b1a1"}},
		{"near king c1", "3k/4/4/2KN w N - 0 1", []string{"c1d1"}},
		{"near king d1 corner", "3k/4/4/2NK w N - 0 1", []string{"d1c1"}},
		// a king on c1 has its center slot on b1 (near goes toward the d edge)
		{"center king c1", "3k/4/4/1NK1 w C - 0 1", []string{"c1b1"}},

		// center castle (adjacent swap)
		{"center king b1", "3k/4/4/1KP1 w C - 0 1", []string{"b1c1"}},
		// a corner king's center slot is two files out: the pair crosses
		{"center king a1 crossing", "3k/4/4/K1P1 w C - 0 1", []string{"a1c1"}},

		// far castle crossings, pawn or knight alike
		{"far pawn king b1", "3k/4/4/1K1P w NCF - 0 1", []string{"b1d1"}},
		{"far knight king b1", "3k/4/4/1K1N w F - 0 1", []string{"b1d1"}},
		// corner-to-corner: the king hops two squares, the partner crosses
		{"far knight opposite corner", "3k/4/4/K2N w F - 0 1", []string{"a1d1"}},
		{"far pawn opposite corner", "3k/4/4/K2P w F - 0 1", []string{"a1d1"}},

		// every square between the pair must be empty
		{"far blocked by center piece", "3k/4/4/1KPP w CF - 0 1", []string{"b1c1"}},
		{"corner far blocked at b1", "3k/4/4/KN1P w NF - 0 1", []string{"a1b1"}},
		{"corner far blocked at c1", "3k/4/4/K1PN w CF - 0 1", []string{"a1c1"}},

		// two adjacent partners: distinct near and center castles
		{"king between partners", "3k/4/4/PKP1 w NC - 0 1", []string{"b1a1", "b1c1"}},

		// rights are slot-based: a knight on the far slot needs the F right
		{"wrong slot right", "3k/4/4/1K1N w N - 0 1", nil},
		// rights are required: same geometry, no rights -> nothing
		{"no rights", "3k/4/4/NK2 w - - 0 1", nil},
		// the king must be on its home rank
		{"king off home rank", "3k/4/1K2/N3 w N - 0 1", nil},
		// no castling while in check (black rook down the b-file)
		{"in check", "1r1k/4/4/NK2 w N - 0 1", nil},
		// no castling through check: the rook covers b1, which the king crosses
		{"through check", "1r1k/4/4/K2P w F - 0 1", nil},
		// no castling into check: the rook covers c1, where the king lands
		{"into check", "2rk/4/4/K2P w F - 0 1", nil},

		// black mirrors, home rank 4 with slots running toward the a-file
		{"black far corner", "n2k/4/4/K3 b f - 0 1", []string{"d4a4"}},
		{"black center knight", "1n1k/4/4/K3 b c - 0 1", []string{"d4b4"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := sortedKeys(castleMoveSet(t, tc.ofen))
			want := tc.want
			sort.Strings(want)
			if strings.Join(got, ",") != strings.Join(want, ",") {
				t.Fatalf("castles = %v, want %v", got, want)
			}
		})
	}
}

func TestCastleApplyResult(t *testing.T) {
	cases := []struct {
		name string
		ofen string
		uoi  string
		want string // resulting OFEN
	}{
		// adjacent swaps
		{"near swap b1", "3k/4/4/NK2 w N - 0 1", "b1a1", "3k/4/4/KN2 b - - 0 1"},
		{"center swap b1", "3k/4/4/1KP1 w C - 0 1", "b1c1", "3k/4/4/1PK1 b - - 0 1"},
		// one-gap crossings: the king lands on the gap, the partner behind it
		{"far pawn crossing b1d1", "3k/4/4/1K1P w NCF - 0 1", "b1d1", "3k/4/4/1PK1 b - - 0 1"},
		{"far knight crossing b1d1", "3k/4/4/1K1N w F - 0 1", "b1d1", "3k/4/4/1NK1 b - - 0 1"},
		{"center pawn crossing a1c1", "3k/4/4/K1P1 w C - 0 1", "a1c1", "3k/4/4/PK2 b - - 0 1"},
		// corner-to-corner: the king hops two squares to c1, the partner to b1
		{"corner knight crossing a1d1", "3k/4/4/K2N w F - 0 1", "a1d1", "3k/4/4/1NK1 b - - 0 1"},
		{"corner pawn crossing a1d1", "3k/4/4/K2P w F - 0 1", "a1d1", "3k/4/4/1PK1 b - - 0 1"},
		// black mirrors
		{"black near swap", "2kn/4/4/K3 b n - 0 1", "c4d4", "2nk/4/4/K3 w - - 0 2"},
		{"black corner crossing", "n2k/4/4/K3 b f - 0 1", "d4a4", "1kn1/4/4/K3 w - - 0 2"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := gameFromOFEN(t, tc.ofen)
			applyUOI(t, g, tc.uoi)
			if got := g.OFEN(); got != tc.want {
				t.Fatalf("after %s: ofen = %q, want %q", tc.uoi, got, tc.want)
			}
		})
	}
}

func TestCastleRightsLoss(t *testing.T) {
	cases := []struct {
		name string
		ofen string
		uoi  string
		want string // resulting castle-rights field
	}{
		// any king move forfeits all of that side's rights
		{"king move drops all", "ppkn/4/4/NKPP w NCFncf - 0 1", "b1b2", "ncf"},
		// moving the near-slot knight drops only the near right
		{"knight move drops N", "ppkn/4/4/NKPP w NCFncf - 0 1", "a1c2", "CFncf"},
		// moving the center-slot pawn drops only its right
		{"center pawn move drops C", "3k/4/4/1KPP w CF - 0 1", "c1c2", "F"},
		// slots are square-relative: a pawn on the far slot of a c1 king
		// carries the F right, so moving it drops F and leaves N
		{"far slot pawn drops F", "3k/4/4/P1KN w NF - 0 1", "a1a2", "N"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := gameFromOFEN(t, tc.ofen)
			applyUOI(t, g, tc.uoi)
			fields := strings.Split(g.OFEN(), " ")
			if fields[2] != tc.want {
				t.Fatalf("after %s: rights = %q, want %q (ofen %q)", tc.uoi, fields[2], tc.want, g.OFEN())
			}
		})
	}
}

// TestCastleNotation verifies that O / O-O / O-O-O track the near, center,
// and far slots regardless of piece type or distance, and that each decodes
// back to the same move.
func TestCastleNotation(t *testing.T) {
	cases := []struct {
		name string
		ofen string
		uoi  string
		san  string
	}{
		{"near is O", "ppkn/4/4/NKPP w NCFncf - 0 1", "b1a1", "O"},
		{"center is O-O", "ppkn/4/4/NKPP w NCFncf - 0 1", "b1c1", "O-O"},
		{"crossing center is O-O", "3k/4/4/K1P1 w C - 0 1", "a1c1", "O-O"},
		{"far pawn is O-O-O", "3k/4/4/1K1P w NCF - 0 1", "b1d1", "O-O-O"},
		{"corner knight is O-O-O", "3k/4/4/K2N w F - 0 1", "a1d1", "O-O-O"},
		{"adjacent partners near", "3k/4/4/PKP1 w NC - 0 1", "b1a1", "O"},
		{"adjacent partners center", "3k/4/4/PKP1 w NC - 0 1", "b1c1", "O-O"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := gameFromOFEN(t, tc.ofen)
			pos := g.Position()

			var move *Move
			for _, m := range g.ValidMoves() {
				if m.String() == tc.uoi {
					move = m
					break
				}
			}
			if move == nil {
				t.Fatalf("move %s not legal in %q", tc.uoi, tc.ofen)
			}

			if got := (AlgebraicNotation{}).Encode(pos, move); got != tc.san {
				t.Fatalf("encode %s: san = %q, want %q", tc.uoi, got, tc.san)
			}

			decoded, err := AlgebraicNotation{}.Decode(pos, tc.san)
			if err != nil {
				t.Fatalf("decode %q: %v", tc.san, err)
			}
			if decoded.String() != tc.uoi {
				t.Fatalf("decode %q: move = %s, want %s", tc.san, decoded, tc.uoi)
			}

			// UOI decode must reproduce the slot tag so the move matches the
			// generated one exactly
			uoiMove, err := UOINotation{}.Decode(pos, tc.uoi)
			if err != nil {
				t.Fatalf("uoi decode %q: %v", tc.uoi, err)
			}
			if !uoiMove.Equals(move) {
				t.Fatalf("uoi decode %q: tags = %b, want %b", tc.uoi, uoiMove.tags, move.tags)
			}
		})
	}
}
