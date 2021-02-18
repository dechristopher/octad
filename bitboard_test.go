package octad

import "testing"

type bitboardTestPair struct {
	initial  uint16
	reversed uint16
}

var (
	tests = []bitboardTestPair{
		{
			uint16(1),
			uint16(32768),
		},
		{
			uint16(24576),
			uint16(6),
		},
		{
			uint16(15),
			uint16(61440),
		},
	}
)

func TestBitboardReverse(t *testing.T) {
	for _, p := range tests {
		r := uint16(bitboard(p.initial).Reverse())
		if r != p.reversed {
			t.Fatalf("bitboard reverse of %s expected %s "+
				"but got %s", intStr(p.initial), intStr(p.reversed), intStr(r))
		}
	}
}

func TestBitboardOccupied(t *testing.T) {
	m := map[Square]bool{
		B3: true,
	}
	bb := newBitboard(m)
	if bb.Occupied(B3) != true {
		t.Fatalf("bitboard occupied of %s expected %t "+
			"but got %t", bb, true, false)
	}
}

func BenchmarkBitboardReverse(b *testing.B) {
	var i uint16
	for i = 0; i < uint16(b.N); i++ {
		u := uint16(42345)
		bitboard(u).Reverse()
	}
}

func intStr(i uint16) string {
	return bitboard(i).String()
}
