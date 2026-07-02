package octad

import (
	"log"
	"testing"
)

func unsafeOFEN(s string) *Position {
	pos, err := decodeOFEN(s)
	if err != nil {
		log.Fatal(err)
	}
	return pos
}

type perfTest struct {
	pos           *Position
	nodesPerDepth []int
}

// https://www.chessprogramming.org/Perft_Results
//
// Counts were lowered when castling became position-relative. The legacy
// castleMoves emitted a close/far castle from the king's home square based only
// on the castle-rights bit and the gap being empty, without verifying that the
// partner pawn was actually present. That produced phantom castles — e.g., with
// the king on b1 and c1 empty, it generated "b1c1" a second time (already a
// legal king move) and would have materialized a pawn if played. Requiring the
// real partner piece removes those duplicates, so the corrected node counts are
// strictly smaller.
//
// Counts rose again when the castle slots became square-relative. Rights now
// stay attached to their home-rank squares, so (for example) after 1. c2 the
// d1 pawn still castles under the F right — the piece-identity model reshuffled
// it into the forfeited C slot, wrongly killing the far castle for the rest of
// the game.
var perfResults = []perfTest{
	{pos: unsafeOFEN("ppkn/4/4/NKPP w NCFncf - 0 1"), nodesPerDepth: []int{
		10, 84, 642, 4375, 29309, 183931,
	}},
	{pos: unsafeOFEN("ppkn/4/2P1/NK1P w NCFncf - 0 1"), nodesPerDepth: []int{
		9, 66, 461, 2924,
	}},
}

func TestPerfResults(t *testing.T) {
	for _, perf := range perfResults {
		countMoves(t, perf.pos, []*Position{perf.pos}, perf.nodesPerDepth, len(perf.nodesPerDepth))
	}
}

func countMoves(t *testing.T, originalPosition *Position, positions []*Position, nodesPerDepth []int, maxDepth int) {
	if len(nodesPerDepth) == 0 {
		return
	}
	depth := maxDepth - len(nodesPerDepth) + 1
	expNodes := nodesPerDepth[0]
	newPositions := make([]*Position, 0)
	for _, pos := range positions {
		for _, move := range pos.ValidMoves() {
			newPos := pos.Update(move)
			newPositions = append(newPositions, newPos)
		}
	}
	gotNodes := len(newPositions)
	if expNodes != gotNodes {
		t.Errorf("Depth: %d Expected: %d Got: %d", depth, expNodes, gotNodes)
		t.Log("##############################")
		t.Log("# Original position info")
		t.Log("###")
		t.Log(originalPosition.String())
		t.Log(originalPosition.board.Draw())
		t.Log("##############################")
		t.Log("# Details in JSONL (https://jsonlines.org)")
		t.Log("###")
		//for _, pos := range positions {
		//	//t.Logf(`{"position": "%s", "moves": %d}`, pos.String(), len(pos.ValidMoves()))
		//}
	}
	countMoves(t, originalPosition, newPositions, nodesPerDepth[1:], maxDepth)
}

func BenchmarkValidMoves(b *testing.B) {
	pos := unsafeOFEN("ppkn/4/4/NKPP w NCFncf - 0 1")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pos.ValidMoves()
		pos.validMoves = nil
	}
}
