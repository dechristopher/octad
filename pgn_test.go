package octad

import (
	"strings"
	"testing"
)

func TestPGNGameOption(t *testing.T) {
	pgn := "[Event \"Test\"]\n\n1. d2 a3 *\n"
	pgnFunc, err := PGN(strings.NewReader(pgn))
	if err != nil {
		t.Fatalf("pgn: %v", err)
	}
	g, err := NewGame(pgnFunc)
	if err != nil {
		t.Fatalf("pgn: new game: %v", err)
	}
	moves := g.Moves()
	if len(moves) != 2 {
		t.Fatalf("pgn: expected 2 moves, got %d", len(moves))
	}
	if moves[0].String() != "d1d2" || moves[1].String() != "a4a3" {
		t.Fatalf("pgn: moves = %s, %s; want d1d2, a4a3", moves[0], moves[1])
	}
	if tp := g.GetTagPair("Event"); tp == nil || tp.Value != "Test" {
		t.Fatalf("pgn: expected Event tag pair Test, got %v", tp)
	}
}

func TestPGNRoundtrip(t *testing.T) {
	g, err := NewGame()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range []string{"Nc2", "b3", "d2"} {
		if err := g.MoveStr(m); err != nil {
			t.Fatal(err)
		}
	}
	pgn := g.String()
	parsed, err := PGN(strings.NewReader(pgn))
	if err != nil {
		t.Fatalf("pgn: %v", err)
	}
	cp, err := NewGame(parsed)
	if err != nil {
		t.Fatal(err)
	}
	if cp.String() != g.String() {
		t.Fatalf("pgn: round-trip mismatch\n got: %q\nwant: %q", cp.String(), g.String())
	}
}

func TestPGNFromFENTag(t *testing.T) {
	// FEN tag seeds the start position; the near castle "O" is then applied
	pgn := "[FEN \"3k/4/4/NK2 w N - 0 1\"]\n\n1. O *\n"
	pgnFunc, err := PGN(strings.NewReader(pgn))
	if err != nil {
		t.Fatalf("pgn: %v", err)
	}
	g, err := NewGame(pgnFunc)
	if err != nil {
		t.Fatal(err)
	}
	if got := g.Position().String(); got != "3k/4/4/KN2 b - - 0 1" {
		t.Fatalf("pgn: FEN-seeded near castle produced %q, want 3k/4/4/KN2 b - - 0 1", got)
	}
}

func TestPGNDecodeError(t *testing.T) {
	pgn := "[Event \"x\"]\n\n1. Zz9 *\n"
	if _, err := PGN(strings.NewReader(pgn)); err == nil {
		t.Error("pgn: expected error decoding invalid move")
	}
}

func multiGamePGN() string {
	return "[Event \"G1\"]\n\n1. d2 a3 *\n\n" +
		"[Event \"G2\"]\n\n1. c2 b3 *\n\n"
}

func TestGamesFromPGN(t *testing.T) {
	games, err := GamesFromPGN(strings.NewReader(multiGamePGN()))
	if err != nil {
		t.Fatalf("pgn: %v", err)
	}
	if len(games) != 2 {
		t.Fatalf("pgn: expected 2 games, got %d", len(games))
	}
	if len(games[0].Moves()) != 2 || len(games[1].Moves()) != 2 {
		t.Fatalf("pgn: expected 2 moves per game, got %d and %d",
			len(games[0].Moves()), len(games[1].Moves()))
	}
}

func TestScanner(t *testing.T) {
	s := NewScanner(strings.NewReader(multiGamePGN()))
	var games []*Game
	for s.Scan() {
		games = append(games, s.Next())
	}
	if err := s.Err(); err != nil {
		t.Fatalf("pgn: scanner error: %v", err)
	}
	if len(games) != 2 {
		t.Fatalf("pgn: scanner found %d games, want 2", len(games))
	}
	if games[0].GetTagPair("Event").Value != "G1" || games[1].GetTagPair("Event").Value != "G2" {
		t.Fatalf("pgn: scanner game tags wrong: %v %v",
			games[0].GetTagPair("Event"), games[1].GetTagPair("Event"))
	}
}

func TestScannerError(t *testing.T) {
	bad := "[Event \"x\"]\n\n1. Zz9 *\n\n"
	s := NewScanner(strings.NewReader(bad))
	if s.Scan() {
		t.Error("pgn: scanner should fail on invalid move")
	}
	if s.Err() == nil {
		t.Error("pgn: scanner should report an error on invalid move")
	}
}
