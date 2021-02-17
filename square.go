package octad

const (
	squaresOnBoard = 16
	squaresInRow   = 4
)

// A Square is one of the 16 rank and file combinations that make up a board.
type Square int8

// File returns the square's file.
func (sq Square) File() File {
	return File(int(sq) % squaresInRow)
}

// Rank returns the square's rank.
func (sq Square) Rank() Rank {
	return Rank(int(sq) / squaresInRow)
}

func (sq Square) String() string {
	return sq.File().String() + sq.Rank().String()
}

// Color returns the color of a given square
func (sq Square) Color() Color {
	if ((sq / 4) % 2) == (sq % 2) {
		return Black
	}
	return White
}

func getSquare(f File, r Rank) Square {
	return Square((int(r) * 4) + int(f))
}

const (
	NoSquare Square = iota - 1
	A1
	B1
	C1
	D1
	A2
	B2
	C2
	D2
	A3
	B3
	C3
	D3
	A4
	B4
	C4
	D4
)

const (
	fileChars = "abcd"
	rankChars = "1234"
)

// A Rank is the rank of a square.
type Rank int8

const (
	Rank1 Rank = iota
	Rank2
	Rank3
	Rank4
)

func (r Rank) String() string {
	return rankChars[r : r+1]
}

// A File is the file of a square.
type File int8

const (
	FileA File = iota
	FileB
	FileC
	FileD
)

func (f File) String() string {
	return fileChars[f : f+1]
}

var (
	strToSquareMap = map[string]Square{
		"a1": A1, "a2": A2, "a3": A3, "a4": A4,
		"b1": B1, "b2": B2, "b3": B3, "b4": B4,
		"c1": C1, "c2": C2, "c3": C3, "c4": C4,
		"d1": D1, "d2": D2, "d3": D3, "d4": D4,
	}
)
