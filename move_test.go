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
var perfResults = []perfTest{
	{pos: unsafeOFEN("ppkn/4/4/NKPP w NCFncf - 0 1"), nodesPerDepth: []int{
		10, 84, 642, 4348, 29382, 185195,
	}},
	{pos: unsafeOFEN("ppkn/4/2P1/NK1P w NCFncf - 0 1"), nodesPerDepth: []int{
		10, 75, 548, 3483,
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
		t.Log("# Details in JSONL (http://jsonlines.org)")
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
