package octad

import (
	"testing"
)

func TestPositionBinary(t *testing.T) {
	for _, ofen := range validOFENs {
		pos, err := decodeOFEN(ofen)
		if err != nil {
			t.Fatal(err)
		}
		b, err := pos.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		cp := &Position{}
		if err := cp.UnmarshalBinary(b); err != nil {
			t.Fatal(err)
		}
		if pos.String() != cp.String() {
			t.Fatalf("position: expected %s but got %s", pos.String(), cp.String())
		}
	}
}
