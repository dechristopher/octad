package octad

import (
	"log"
	"strings"
	"testing"
)

func TestUndoMove(t *testing.T) {
	g, err := NewGame()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	if err = g.MoveStr("c2"); err != nil {
		t.Fatal(err)
	}

	g.UndoMove()

	if g.Position().String() != startOFEN {
		t.Fatalf("game: expected default OFEN after UndoMove but got %s", g.Position().String())
	}

	if len(g.Moves()) != 0 {
		t.Fatalf("game: expected no moves played after UndoMove but got %d", len(g.Moves()))
	}
}

func TestUndoMoveState(t *testing.T) {
	startPos := "3R/4/2K1/k3 w - - 7 7"
	ofen, err := OFEN(startPos)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	g, err := NewGame(ofen)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	legalMoves := g.ValidMoves()

	// white checkmates black with rook to a4
	if err = g.MoveStr("Ra4"); err != nil {
		t.Fatal(err)
	}

	if g.Outcome() != WhiteWon || g.Method() != Checkmate {
		t.Fatalf("game: expected checkmate after d4a4 but got %s - %s",
			g.Outcome().String(), g.Method().String())
	}

	g.UndoMove()

	if g.Position().String() != startPos {
		t.Fatalf("game: expected test start position after UndoMove but got %s", g.Position().String())
	}

	if g.Outcome() != NoOutcome || g.Method() != NoMethod {
		t.Fatalf("game: expected no outcome after UndoMove but got %s - %s",
			g.Outcome().String(), g.Method().String())
	}

	if len(g.ValidMoves()) != len(legalMoves) {
		t.Fatalf("game: expected %d valid moes after UndoMove but got %d",
			len(legalMoves), len(g.ValidMoves()))
	}
}

func TestCheckmate(t *testing.T) {
	ofenStr := "4/1K2/1Q2/3k w - - 7 7"
	ofen, err := OFEN(ofenStr)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(ofen)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	if err = g.MoveStr("Qc2#"); err != nil {
		t.Fatal(err)
	}
	if g.Method() != Checkmate {
		t.Fatalf("game: expected method %s but got %s", Checkmate, g.Method())
	}
	if g.Outcome() != WhiteWon {
		t.Fatalf("game: expected outcome %s but got %s", WhiteWon, g.Outcome())
	}

	// TODO is checkmating by castling possible?
}

func TestCheckmateFromOFEN(t *testing.T) {
	ofenStr := "4/1K2/2Q1/3k b - - 8 7"
	ofen, err := OFEN(ofenStr)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(ofen)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	if g.Method() != Checkmate {
		t.Error(g.Position().Board().Draw())
		t.Fatalf("game: expected method %s but got %s", Checkmate, g.Method())
	}
	if g.Outcome() != WhiteWon {
		t.Fatalf("game: expected outcome %s but got %s", WhiteWon, g.Outcome())
	}
}

func TestStalemate(t *testing.T) {
	ofenStr := "4/1K2/Q3/3k w - - 7 7"
	ofen, err := OFEN(ofenStr)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(ofen)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	if err = g.MoveStr("Qb2"); err != nil {
		t.Fatal(err)
	}
	if g.Method() != Stalemate {
		t.Fatalf("game: expected method %s but got %s", Stalemate, g.Method())
	}
	if g.Outcome() != Draw {
		t.Fatalf("game: expected outcome %s but got %s", Draw, g.Outcome())
	}
}

// position shouldn't result in stalemate because pawn can move http://en.lichess.org/Pc6mJDZN#138
func TestInvalidStalemate(t *testing.T) {
	ofenStr := "4/k1P1/p1K1/4 w - - 7 7"
	ofen, err := OFEN(ofenStr)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(ofen)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	if err = g.MoveStr("c4=Q"); err != nil {
		t.Fatal(err)
	}
	if g.Outcome() != NoOutcome {
		t.Fatalf("game: expected outcome %s but got %s", NoOutcome, g.Outcome())
	}
}

func TestThreeFoldRepetition(t *testing.T) {
	g, err := NewGame()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	moves := []string{
		"Nc2", "Nb3", "Na1", "Nd4",
		"Nc2", "Nb3", "Na1", "Nd4",
		"Nc2", "Nb3", "Na1", "Nd4",
	}
	for _, m := range moves {
		if err = g.MoveStr(m); err != nil {
			t.Fatal(err)
		}
	}
	if err = g.Draw(ThreefoldRepetition); err != nil {
		for _, pos := range g.Positions() {
			log.Println(pos.String())
		}
		t.Fatalf("game: %s - %d reps", err.Error(), g.numRepetitions())
	}
}

func TestInvalidThreeFoldRepetition(t *testing.T) {
	g, err := NewGame()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	moves := []string{
		"Nc2", "Nb3", "Na1", "Nd4",
		"Nc2", "Nb3", "Na1", "Nd4",
	}
	for _, m := range moves {
		if err := g.MoveStr(m); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.Draw(ThreefoldRepetition); err == nil {
		t.Fatal("game: should require three repeated board states")
	}
}

func TestTwentyFiveMoveRule(t *testing.T) {
	ofen, _ := OFEN("n2k/4/3K/N3 b - - 49 24")
	g, err := NewGame(ofen)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	if err = g.MoveStr("Kc4"); err != nil {
		t.Fatal(err)
	}
	if g.Outcome() != Draw || g.Method() != TwentyFiveMoveRule {
		t.Fatal("game: should automatically draw after twenty five moves w/ no pawn move or capture")
	}
}

func TestInsufficientMaterial(t *testing.T) {
	ofens := []string{
		"4/1n1k/4/1K2 w - - 1 5",
		"3k/4/4/K3 w - - 1 1",
		"k2n/4/4/1K2 w - - 1 1",
		"bk2/4/4/K3 w - - 1 1",
	}
	for _, o := range ofens {
		ofen, err := OFEN(o)
		if err != nil {
			t.Fatal(err)
		}
		g, err := NewGame(ofen)
		if err != nil {
			t.Fatalf(err.Error())
			return
		}
		if g.Outcome() != Draw || g.Method() != InsufficientMaterial {
			log.Println(g.Position().Board().Draw())
			t.Fatalf("game: %s should automatically draw by insufficient material", o)
		}
	}
}

func TestSufficientMaterial(t *testing.T) {
	ofens := []string{
		"k3/3P/2K1/4 w - - 1 1",
		"kbn1/4/4/K3 w - - 1 1",
	}
	for _, o := range ofens {
		ofen, err := OFEN(o)
		if err != nil {
			t.Fatal(err)
		}
		g, err := NewGame(ofen)
		if err != nil {
			t.Fatalf(err.Error())
			return
		}
		if g.Outcome() != NoOutcome {
			log.Println(g.Position().Board().Draw())
			t.Fatalf("game: %s should not find insufficient material", o)
		}
	}
}

func TestSerializationCycle(t *testing.T) {
	g, err := NewGame()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	pgn, err := PGN(strings.NewReader(g.String()))
	if err != nil {
		t.Fatal(err)
	}
	cp, err := NewGame(pgn)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	if cp.String() != g.String() {
		t.Fatalf("game: expected %s but got %s", g.String(), cp.String())
	}
}

func TestInitialNumOfValidMoves(t *testing.T) {
	g, err := NewGame()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	if len(g.ValidMoves()) != 10 {
		t.Fatal("game: should find 10 valid moves from the initial position")
	}
}

func TestTagPairs(t *testing.T) {
	g, err := NewGame()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	g.AddTagPair("Draw Offer", "White")
	tagPair := g.GetTagPair("Draw Offer")
	if tagPair == nil {
		t.Fatalf("game: expected %s but got %s", "White", "nil")
	}
	if tagPair.Value != "White" {
		t.Fatalf("game: expected %s but got %s", "White", tagPair.Value)
	}
	g.RemoveTagPair("Draw Offer")
	tagPair = g.GetTagPair("Draw Offer")
	if tagPair != nil {
		t.Fatalf("game: expected %s but got %s", "nil", "not nil")
	}
}

func TestPositionHash(t *testing.T) {
	g1, err := NewGame()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	for _, s := range []string{"Nc2", "b3", "d2"} {
		err = g1.MoveStr(s)
		if err != nil {
			t.Fatalf(err.Error())
			return
		}
	}
	g2, err := NewGame()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	for _, s := range []string{"d2", "b3", "Nc2"} {
		err = g2.MoveStr(s)
		if err != nil {
			t.Fatalf(err.Error())
			return
		}
	}
	if g1.Position().Hash() != g2.Position().Hash() {
		t.Fatalf("game: expected position hashes to be equal but got %s and %s", g1.Position().Hash(), g2.Position().Hash())
	}
}

func BenchmarkStalemateStatus(b *testing.B) {
	ofenStr := "3k/4/3p/3K b - - 2 10"
	ofen, err := OFEN(ofenStr)
	if err != nil {
		b.Fatal(err)
	}
	g, err := NewGame(ofen)
	if err != nil {
		b.Fatalf(err.Error())
		return
	}
	if err = g.MoveStr("Kd3"); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.Position().Status()
	}
}

func BenchmarkInvalidStalemateStatus(b *testing.B) {
	ofenStr := "4/P2k/3p/3K w - - 2 10"
	ofen, err := OFEN(ofenStr)
	if err != nil {
		b.Fatal(err)
	}
	g, err := NewGame(ofen)
	if err != nil {
		b.Fatalf(err.Error())
		return
	}
	if err = g.MoveStr("a4=Q"); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.Position().Status()
	}
}

func BenchmarkPositionHash(b *testing.B) {
	ofenStr := "4/P2k/3p/3K w - - 2 10"
	ofen, err := OFEN(ofenStr)
	if err != nil {
		b.Fatal(err)
	}
	g, err := NewGame(ofen)
	if err != nil {
		b.Fatalf(err.Error())
		return
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.Position().Hash()
	}
}
