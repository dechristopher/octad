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
	// NoSquare represents an invalid square
	NoSquare Square = iota - 1
	A1              // The A1 square, index 0
	B1              // The B1 square, index 1
	C1              // The C1 square, index 2
	D1              // The D1 square, index 3
	A2              // The A2 square, index 4
	B2              // The B2 square, index 5
	C2              // The C2 square, index 6
	D2              // The D2 square, index 7
	A3              // The A3 square, index 8
	B3              // The B3 square, index 9
	C3              // The C3 square, index 10
	D3              // The D3 square, index 11
	A4              // The A4 square, index 12
	B4              // The B4 square, index 13
	C4              // The C4 square, index 14
	D4              // The D4 square, index 15
)

const (
	fileChars = "abcd"
	rankChars = "1234"
)

// A Rank is the rank of a square.
type Rank int8

const (
	Rank1 Rank = iota // Rank1 is the first rank, index 0
	Rank2             // Rank2 is the second rank, index 1
	Rank3             // Rank3 is the third rank, index 2
	Rank4             // Rank4 is the fourth rank, index 3
)

func (r Rank) String() string {
	return rankChars[r : r+1]
}

// A File is the file of a square.
type File int8

const (
	FileA File = iota // FileA is the A file, index 0
	FileB             // FileB is the B file, index 1
	FileC             // FileC is the C file, index 2
	FileD             // FileD is the D file, index 3
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
