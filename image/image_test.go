package image

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image/color"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/dechristopher/octad"
)

const expectedMD5 = "44481de63ea7b587f37e791eb55a99aa"

func TestSVG(t *testing.T) {
	// create buffer of actual svg
	buf := bytes.NewBuffer([]byte{})
	fenStr := "ppkn/4/2P1/NK1P b NCFncf - 0 1"
	pos := &octad.Position{}
	if err := pos.UnmarshalText([]byte(fenStr)); err != nil {
		t.Error(err)
	}
	mark := MarkSquares(color.RGBA{R: 255, G: 255, A: 1}, octad.C1, octad.C2)
	if err := SVG(buf, pos.Board(), mark); err != nil {
		t.Error(err)
	}

	// compare to expected svg
	actualSVG := strings.TrimSpace(buf.String())
	actualMD5 := fmt.Sprintf("%x", md5.Sum([]byte(actualSVG)))
	if actualMD5 != expectedMD5 {
		t.Errorf("expected actual md5 hash to be %s but got %s", expectedMD5, actualMD5)
	}

	// create actual svg file for visualization
	f, err := os.Create("example.svg")
	defer f.Close()
	if err != nil {
		t.Error(err)
	}
	if _, err := io.Copy(f, bytes.NewBufferString(actualSVG)); err != nil {
		t.Error(err)
	}
}
