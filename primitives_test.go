package octad

import (
	"strings"
	"testing"
)

func TestColorOther(t *testing.T) {
	if White.Other() != Black {
		t.Errorf("color: White.Other() = %v, want Black", White.Other())
	}
	if Black.Other() != White {
		t.Errorf("color: Black.Other() = %v, want White", Black.Other())
	}
	if NoColor.Other() != NoColor {
		t.Errorf("color: NoColor.Other() = %v, want NoColor", NoColor.Other())
	}
}

func TestColorStringAndName(t *testing.T) {
	cases := []struct {
		c    Color
		str  string
		name string
	}{
		{White, "w", "White"},
		{Black, "b", "Black"},
		{NoColor, "-", "No Color"},
	}
	for _, tc := range cases {
		if tc.c.String() != tc.str {
			t.Errorf("color: %v String() = %q, want %q", tc.c, tc.c.String(), tc.str)
		}
		if tc.c.Name() != tc.name {
			t.Errorf("color: %v Name() = %q, want %q", tc.c, tc.c.Name(), tc.name)
		}
	}
}

func TestPieceTypesAndPromotable(t *testing.T) {
	got := PieceTypes()
	want := [6]PieceType{King, Queen, Rook, Bishop, Knight, Pawn}
	if got != want {
		t.Fatalf("PieceTypes() = %v, want %v", got, want)
	}
	promotable := map[PieceType]bool{
		King: false, Queen: true, Rook: true, Bishop: true, Knight: true, Pawn: false, NoPieceType: false,
	}
	for pt, exp := range promotable {
		if pt.promotableTo() != exp {
			t.Errorf("piece: %v promotableTo() = %v, want %v", pt, pt.promotableTo(), exp)
		}
	}
}

func TestPieceTypeColorAndChars(t *testing.T) {
	cases := []struct {
		p    Piece
		typ  PieceType
		col  Color
		ofen string
	}{
		{WhiteKing, King, White, "K"},
		{WhiteQueen, Queen, White, "Q"},
		{WhiteRook, Rook, White, "R"},
		{WhiteBishop, Bishop, White, "B"},
		{WhiteKnight, Knight, White, "N"},
		{WhitePawn, Pawn, White, "P"},
		{BlackKing, King, Black, "k"},
		{BlackPawn, Pawn, Black, "p"},
		{NoPiece, NoPieceType, NoColor, ""},
	}
	for _, tc := range cases {
		if tc.p.Type() != tc.typ {
			t.Errorf("piece: %v Type() = %v, want %v", tc.p, tc.p.Type(), tc.typ)
		}
		if tc.p.Color() != tc.col {
			t.Errorf("piece: %v Color() = %v, want %v", tc.p, tc.p.Color(), tc.col)
		}
		if tc.p.getOFENChar() != tc.ofen {
			t.Errorf("piece: %v getOFENChar() = %q, want %q", tc.p, tc.p.getOFENChar(), tc.ofen)
		}
	}
}

func TestGetPiece(t *testing.T) {
	if getPiece(Knight, White) != WhiteKnight {
		t.Errorf("getPiece(Knight, White) = %v, want WhiteKnight", getPiece(Knight, White))
	}
	if getPiece(Pawn, Black) != BlackPawn {
		t.Errorf("getPiece(Pawn, Black) = %v, want BlackPawn", getPiece(Pawn, Black))
	}
	if getPiece(King, NoColor) != NoPiece {
		t.Errorf("getPiece(King, NoColor) = %v, want NoPiece", getPiece(King, NoColor))
	}
}

func TestPieceUnicodeString(t *testing.T) {
	if WhiteKing.String() != "♔" {
		t.Errorf("piece: WhiteKing.String() = %q, want ♔", WhiteKing.String())
	}
	if BlackPawn.String() != "♟" {
		t.Errorf("piece: BlackPawn.String() = %q, want ♟", BlackPawn.String())
	}
	if NoPiece.String() != " " {
		t.Errorf("piece: NoPiece.String() = %q, want space", NoPiece.String())
	}
}

func TestSquareColor(t *testing.T) {
	// rank 1 alternates dark/light starting dark on a1
	cases := map[Square]Color{
		A1: Black, B1: White, C1: Black, D1: White,
		A4: White, D4: Black,
	}
	for sq, want := range cases {
		if sq.Color() != want {
			t.Errorf("square: %s Color() = %v, want %v", sq, sq.Color(), want)
		}
	}
}

func TestSquareFileRankString(t *testing.T) {
	if A1.String() != "a1" || D4.String() != "d4" || B3.String() != "b3" {
		t.Errorf("square: unexpected String() for a1/d4/b3: %s %s %s", A1, D4, B3)
	}
	if A1.File() != FileA || A1.Rank() != Rank1 {
		t.Errorf("square: A1 file/rank = %v/%v, want FileA/Rank1", A1.File(), A1.Rank())
	}
	if D4.File() != FileD || D4.Rank() != Rank4 {
		t.Errorf("square: D4 file/rank = %v/%v, want FileD/Rank4", D4.File(), D4.Rank())
	}
	if getSquare(FileC, Rank2) != C2 {
		t.Errorf("square: getSquare(FileC, Rank2) = %v, want C2", getSquare(FileC, Rank2))
	}
	if FileB.String() != "b" || Rank3.String() != "3" {
		t.Errorf("square: file/rank strings = %q/%q, want b/3", FileB.String(), Rank3.String())
	}
}

func TestBitboardMappingAndString(t *testing.T) {
	m := map[Square]bool{A1: true, D4: true, B3: true}
	bb := newBitboard(m)
	if got := bb.Mapping(); len(got) != 3 || !got[A1] || !got[D4] || !got[B3] {
		t.Fatalf("bitboard: Mapping() = %v, want %v", got, m)
	}
	// A1 is the most significant of 16 bits
	only := newBitboard(map[Square]bool{A1: true})
	if only.String() != "1000000000000000" {
		t.Errorf("bitboard: A1 String() = %q, want 1000000000000000", only.String())
	}
	if len(bb.String()) != squaresOnBoard {
		t.Errorf("bitboard: String() length = %d, want %d", len(bb.String()), squaresOnBoard)
	}
}

func TestBitboardDraw(t *testing.T) {
	bb := newBitboard(map[Square]bool{A1: true})
	draw := bb.Draw()
	if !strings.Contains(draw, "A B C D") {
		t.Errorf("bitboard: Draw() missing header:\n%s", draw)
	}
	if !strings.Contains(draw, "1") {
		t.Errorf("bitboard: Draw() missing set square:\n%s", draw)
	}
}

func TestMethodString(t *testing.T) {
	cases := map[Method]string{
		NoMethod:             "NoMethod",
		Checkmate:            "Checkmate",
		Resignation:          "Resignation",
		DrawOffer:            "DrawOffer",
		Stalemate:            "Stalemate",
		ThreefoldRepetition:  "ThreefoldRepetition",
		TwentyFiveMoveRule:   "TwentyFiveMoveRule",
		InsufficientMaterial: "InsufficientMaterial",
	}
	for m, want := range cases {
		if m.String() != want {
			t.Errorf("method: %d String() = %q, want %q", m, m.String(), want)
		}
	}
	if got := Method(99).String(); got != "Method(99)" {
		t.Errorf("method: out-of-range String() = %q, want Method(99)", got)
	}
}

func TestOutcomeString(t *testing.T) {
	cases := map[Outcome]string{
		NoOutcome: "*",
		WhiteWon:  "1-0",
		BlackWon:  "0-1",
		Draw:      "1/2-1/2",
	}
	for o, want := range cases {
		if o.String() != want {
			t.Errorf("outcome: String() = %q, want %q", o.String(), want)
		}
	}
}
