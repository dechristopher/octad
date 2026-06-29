package octad

import (
	"strings"
	"testing"
)

// ---- Board ----

func TestBoardDrawAndSquareMap(t *testing.T) {
	g, _ := NewGame()
	b := g.Position().Board()

	draw := b.Draw()
	if !strings.Contains(draw, "A B C D") {
		t.Errorf("board: Draw() missing header:\n%s", draw)
	}

	sm := b.SquareMap()
	if len(sm) != 8 {
		t.Fatalf("board: SquareMap() has %d entries, want 8", len(sm))
	}
	if sm[B1] != WhiteKing || sm[C4] != BlackKing {
		t.Errorf("board: SquareMap mislocated kings: b1=%v c4=%v", sm[B1], sm[C4])
	}
	if _, ok := sm[A2]; ok {
		t.Errorf("board: SquareMap should omit empty square a2")
	}
}

func TestBoardPieceLookup(t *testing.T) {
	g, _ := NewGame()
	b := g.Position().Board()
	if b.Piece(A1) != WhiteKnight {
		t.Errorf("board: Piece(a1) = %v, want WhiteKnight", b.Piece(A1))
	}
	if b.Piece(A2) != NoPiece {
		t.Errorf("board: Piece(a2) = %v, want NoPiece", b.Piece(A2))
	}
}

func TestBoardTransposeInvolution(t *testing.T) {
	g, _ := NewGame()
	b := g.Position().Board()
	if got := b.Transpose().Transpose().String(); got != b.String() {
		t.Errorf("board: Transpose twice = %q, want %q", got, b.String())
	}
}

func TestBoardUnmarshalTextError(t *testing.T) {
	b := &Board{}
	if err := b.UnmarshalText([]byte("not-a-board")); err == nil {
		t.Error("board: expected error unmarshalling invalid board text")
	}
}

// ---- Position ----

func TestStartingPosition(t *testing.T) {
	pos, err := StartingPosition()
	if err != nil {
		t.Fatal(err)
	}
	if pos.String() != startOFEN {
		t.Fatalf("position: StartingPosition() = %q, want %q", pos.String(), startOFEN)
	}
	if pos.Turn() != White {
		t.Errorf("position: starting turn = %v, want White", pos.Turn())
	}
}

func TestPositionInCheckAndCheckSquare(t *testing.T) {
	// black to move and in check from the white rook down the d-file. Built via
	// the OFEN game option (not unsafeOFEN) so the inCheck field is computed.
	pos := gameFromOFEN(t, "3k/4/3R/3K b - - 0 1").Position()
	if !pos.InCheck() {
		t.Fatal("position: expected InCheck() true")
	}
	if pos.CheckSquare() != D4 {
		t.Errorf("position: CheckSquare() = %s, want d4", pos.CheckSquare())
	}

	start, _ := StartingPosition()
	if start.InCheck() {
		t.Error("position: start should not be in check")
	}
	if start.CheckSquare() != NoSquare {
		t.Errorf("position: start CheckSquare() = %s, want NoSquare", start.CheckSquare())
	}
}

func TestPositionEnPassantSquare(t *testing.T) {
	pos := unsafeOFEN("3k/4/1Pp1/K3 w - c3 0 2")
	if pos.EnPassantSquare() != C3 {
		t.Errorf("position: EnPassantSquare() = %s, want c3", pos.EnPassantSquare())
	}
	start, _ := StartingPosition()
	if start.EnPassantSquare() != NoSquare {
		t.Errorf("position: start EnPassantSquare() = %s, want NoSquare", start.EnPassantSquare())
	}
}

func TestPositionCastleRights(t *testing.T) {
	start, _ := StartingPosition()
	if start.CastleRights() != "NCFncf" {
		t.Fatalf("position: start CastleRights() = %q, want NCFncf", start.CastleRights())
	}
	for _, side := range []Side{NearSide, CenterSide, FarSide} {
		if !start.CastleRights().CanCastle(White, side) || !start.CastleRights().CanCastle(Black, side) {
			t.Errorf("position: start should allow castling side %v for both colors", side)
		}
	}
	none := CastleRights("-")
	if none.CanCastle(White, NearSide) {
		t.Error("position: '-' rights should not allow castling")
	}
}

func TestPositionTextSerialization(t *testing.T) {
	for _, ofen := range validOFENs {
		pos := unsafeOFEN(ofen)
		txt, err := pos.MarshalText()
		if err != nil {
			t.Fatal(err)
		}
		cp := &Position{}
		if err := cp.UnmarshalText(txt); err != nil {
			t.Fatal(err)
		}
		if cp.String() != pos.String() {
			t.Fatalf("position: text roundtrip got %q want %q", cp.String(), pos.String())
		}
	}
}

func TestPositionHashDiffers(t *testing.T) {
	g, _ := NewGame()
	h1 := g.Position().Hash()
	if err := g.MoveStr("Nc2"); err != nil {
		t.Fatal(err)
	}
	if g.Position().Hash() == h1 {
		t.Error("position: hash should change after a move")
	}
}

func TestPositionUpdateHalfmove(t *testing.T) {
	// a quiet non-pawn move that touches no castling piece increments the clock
	// (rook on a2 -> b2; "R3" is the rank-2 segment of the OFEN)
	pos := unsafeOFEN("3k/4/R3/3K w - - 5 3")
	next := pos.Update(findMove(t, pos, "a2b2"))
	if field := strings.Split(next.String(), " ")[4]; field != "6" {
		t.Errorf("position: halfmove after quiet move = %s, want 6 (%s)", field, next.String())
	}
}

// ---- Game ----

func TestGameResign(t *testing.T) {
	g, _ := NewGame()
	g.Resign(White)
	if g.Outcome() != BlackWon || g.Method() != Resignation {
		t.Fatalf("game: white resign -> %s/%s, want 0-1/Resignation", g.Outcome(), g.Method())
	}
	// resigning a finished game is a no-op
	g.Resign(Black)
	if g.Outcome() != BlackWon {
		t.Errorf("game: resign after completion changed outcome to %s", g.Outcome())
	}

	g2, _ := NewGame()
	g2.Resign(Black)
	if g2.Outcome() != WhiteWon {
		t.Errorf("game: black resign -> %s, want 1-0", g2.Outcome())
	}

	g3, _ := NewGame()
	g3.Resign(NoColor)
	if g3.Outcome() != NoOutcome {
		t.Errorf("game: NoColor resign changed outcome to %s", g3.Outcome())
	}
}

func TestGameClone(t *testing.T) {
	g, _ := NewGame()
	if err := g.MoveStr("Nc2"); err != nil {
		t.Fatal(err)
	}
	cp := g.Clone()
	if cp.String() != g.String() {
		t.Fatalf("game: clone string mismatch\n got %q\nwant %q", cp.String(), g.String())
	}
	// mutating the clone must not affect the original
	if err := cp.MoveStr("b3"); err != nil {
		t.Fatal(err)
	}
	if len(g.Moves()) != 1 {
		t.Errorf("game: original move count changed to %d after clone move", len(g.Moves()))
	}
}

func TestGameEligibleDraws(t *testing.T) {
	g, _ := NewGame()
	if got := g.EligibleDraws(); len(got) != 1 || got[0] != DrawOffer {
		t.Fatalf("game: start EligibleDraws() = %v, want [DrawOffer]", got)
	}
	for _, m := range []string{"Nc2", "Nb3", "Na1", "Nd4", "Nc2", "Nb3", "Na1", "Nd4", "Nc2", "Nb3", "Na1", "Nd4"} {
		if err := g.MoveStr(m); err != nil {
			t.Fatal(err)
		}
	}
	draws := g.EligibleDraws()
	if !containsMethod(draws, DrawOffer) || !containsMethod(draws, ThreefoldRepetition) {
		t.Fatalf("game: EligibleDraws() after threefold = %v, want DrawOffer + ThreefoldRepetition", draws)
	}
}

func TestGameDrawErrors(t *testing.T) {
	g, _ := NewGame()
	if err := g.Draw(Stalemate); err == nil {
		t.Error("game: Draw(Stalemate) should be unsupported")
	}
	if err := g.Draw(ThreefoldRepetition); err == nil {
		t.Error("game: Draw(ThreefoldRepetition) without repetitions should error")
	}
}

func TestGameUseNotation(t *testing.T) {
	g, err := NewGame(UseNotation(UOINotation{}))
	if err != nil {
		t.Fatal(err)
	}
	if err := g.MoveStr("c1c2"); err != nil {
		t.Fatalf("game: UOI MoveStr failed: %v", err)
	}
	if g.Moves()[0].String() != "c1c2" {
		t.Errorf("game: move = %s, want c1c2", g.Moves()[0])
	}
}

func TestGameTextSerialization(t *testing.T) {
	g, _ := NewGame()
	for _, m := range []string{"Nc2", "b3"} {
		if err := g.MoveStr(m); err != nil {
			t.Fatal(err)
		}
	}
	txt, err := g.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	// unmarshal into a NewGame so the notation is initialized
	cp, _ := NewGame()
	if err := cp.UnmarshalText(txt); err != nil {
		t.Fatal(err)
	}
	if cp.String() != g.String() {
		t.Fatalf("game: text roundtrip\n got %q\nwant %q", cp.String(), g.String())
	}
}

func TestGameInvalidInputs(t *testing.T) {
	if _, err := OFEN("totally invalid ofen"); err == nil {
		t.Error("game: OFEN should reject invalid input")
	}
	g, _ := NewGame()
	if err := g.MoveStr("nonsense"); err == nil {
		t.Error("game: MoveStr should reject nonsense")
	}
	if err := g.Move(&Move{s1: A1, s2: A4}); err == nil {
		t.Error("game: Move should reject an illegal move")
	}
}

func TestGameAddTagPairOverwrite(t *testing.T) {
	g, _ := NewGame()
	if g.AddTagPair("Site", "lioctad.org") {
		t.Error("game: first AddTagPair should not report an overwrite")
	}
	if !g.AddTagPair("Site", "example.com") {
		t.Error("game: re-adding a key should report an overwrite")
	}
	if tp := g.GetTagPair("Site"); tp == nil || tp.Value != "example.com" {
		t.Errorf("game: tag value = %v, want example.com", tp)
	}
}

func containsMethod(ms []Method, target Method) bool {
	for _, m := range ms {
		if m == target {
			return true
		}
	}
	return false
}
