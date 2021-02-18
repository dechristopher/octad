package octad

import (
	"strings"
	"testing"
)

var defaultOFEN = strings.Split(startOFEN, " ")[0]

func TestBoardTextSerialization(t *testing.T) {
	fen := defaultOFEN
	b := &Board{}
	if err := b.UnmarshalText([]byte(defaultOFEN)); err != nil {
		t.Fatal("board: received unexpected error", err)
	}
	txt, err := b.MarshalText()
	if err != nil {
		t.Fatal("board: received unexpected error", err)
	}
	if fen != string(txt) {
		t.Fatalf("board: expected board string %s but got %s", fen, string(txt))
	}
}

func TestBoardBinarySerialization(t *testing.T) {
	g, err := NewGame()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	board := g.Position().Board()
	b, err := board.MarshalBinary()
	if err != nil {
		t.Fatal("board: received unexpected error", err)
	}
	cpBoard := &Board{}
	err = cpBoard.UnmarshalBinary(b)
	if err != nil {
		t.Fatal("board: received unexpected error", err)
	}
	s := defaultOFEN
	if s != cpBoard.String() {
		t.Fatalf("board: expected board string %s but got %s", s, cpBoard.String())
	}
}

func TestBoardRotation(t *testing.T) {
	fens := []string{
		"N2p/K2p/P2k/P2n",
		"PPKN/4/4/nkpp",
		"n2P/k2P/p2K/p2N",
		"ppkn/4/4/NKPP",
	}
	g, err := NewGame()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	board := g.Position().Board()
	for i := 0; i < 4; i++ {
		board = board.Rotate()
		if fens[i] != board.String() {
			t.Fatalf("board: expected board string %s but got %s", fens[i], board.String())
		}
	}
}

func TestBoardFlip(t *testing.T) {
	g, err := NewGame()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	board := g.Position().Board()
	board = board.Flip(UpDown)
	b := "NKPP/4/4/ppkn"
	if b != board.String() {
		t.Fatalf("board: flip1 expected board string %s but got %s", b, board.String())
	}
	board = board.Flip(UpDown)
	b = "ppkn/4/4/NKPP"
	if b != board.String() {
		t.Fatalf("board: flip2 expected board string %s but got %s", b, board.String())
	}
	board = board.Flip(LeftRight)
	b = "nkpp/4/4/PPKN"
	if b != board.String() {
		t.Fatalf("board: flip3 expected board string %s but got %s", b, board.String())
	}
	board = board.Flip(LeftRight)
	b = "ppkn/4/4/NKPP"
	if b != board.String() {
		t.Fatalf("board: flip4 expected board string %s but got %s", b, board.String())
	}
}

func TestBoardTranspose(t *testing.T) {
	g, err := NewGame()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	board := g.Position().Board()
	board = board.Transpose()
	b := "p2N/p2K/k2P/n2P"
	if b != board.String() {
		t.Fatalf("board: expected board string %s but got %s", b, board.String())
	}
}
