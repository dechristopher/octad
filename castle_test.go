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
		if m.HasTag(NearCastle) || m.HasTag(CenterCastle) || m.HasTag(FarCastle) {
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
		// the far castle is blocked because the 'center' piece sits in its gap
		{"standard white", "ppkn/4/4/NKPP w NCFncf - 0 1", []string{"b1a1", "b1c1"}},
		{"standard black", "ppkn/4/4/NKPP b NCFncf - 0 1", []string{"c4b4", "c4d4"}},

		// near castle (adjacent swap) from each viable king square
		{"knight king b1", "3k/4/4/NK2 w N - 0 1", []string{"b1a1"}},
		{"knight king c1 left", "3k/4/4/1NK1 w N - 0 1", []string{"c1b1"}},
		{"knight king d1 corner", "3k/4/4/2NK w N - 0 1", []string{"d1c1"}},
		// a knight one gap away cannot leap — no castle
		{"knight one-gap blocked", "3k/4/4/1K1N w N - 0 1", nil},

		// center castle (adjacent swap)
		{"center king b1", "3k/4/4/1KP1 w C - 0 1", []string{"b1c1"}},

		// far castle (one-gap leap) from b1 and the a1 corner
		{"far king b1", "3k/4/4/1K1P w NCF - 0 1", []string{"b1d1"}},
		{"far king a1", "3k/4/4/K1P1 w NCF - 0 1", []string{"a1c1"}},
		// far blocked when the gap (the 'center' piece) is occupied; center remains
		{"far blocked by close", "3k/4/4/1KPP w CF - 0 1", []string{"b1c1"}},

		// rights are required: same geometry, no rights -> nothing
		{"no rights", "3k/4/4/NK2 w - - 0 1", nil},
		// the king must be in its home rank
		{"king off home rank", "3k/4/1K2/N3 w N - 0 1", nil},
		// no castling while in check (black rook down the b-file)
		{"in check", "1r1k/4/4/NK2 w N - 0 1", nil},
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
		{"knight swap b1", "3k/4/4/NK2 w N - 0 1", "b1a1", "3k/4/4/KN2 b - - 0 1"},
		{"close swap b1", "3k/4/4/1KP1 w C - 0 1", "b1c1", "3k/4/4/1PK1 b - - 0 1"},
		{"far leap b1->d1", "3k/4/4/1K1P w NCF - 0 1", "b1d1", "3k/4/4/1PK1 b - - 0 1"},
		{"far leap a1->c1", "3k/4/4/K1P1 w NCF - 0 1", "a1c1", "3k/4/4/PK2 b - - 0 1"},
		{"black knight swap", "2kn/4/4/K3 b n - 0 1", "c4d4", "2nk/4/4/K3 w - - 0 2"},
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
		// moving the knight drops only the knight right
		{"knight move drops N", "ppkn/4/4/NKPP w NCFncf - 0 1", "a1c2", "CFncf"},
		// moving the center pawn drops only it's right
		{"center pawn move drops C", "3k/4/4/1KPP w CF - 0 1", "c1c2", "F"},
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
