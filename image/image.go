// Package image is a go library that creates images from octad board positions
package image

import (
	"fmt"
	"image/color"
	"io"
	"strings"

	svg "github.com/ajstarks/svgo"

	"github.com/dechristopher/octad"
	"github.com/dechristopher/octad/image/internal"
)

// SVG writes the board SVG representation into the writer.
// An error is returned if there is there is an error writing data.
// SVG also takes options which can customize the image output.
func SVG(w io.Writer, b *octad.Board, opts ...func(*encoder)) error {
	e := newEncoder(w, opts)
	return e.EncodeSVG(b)
}

// SquareColors is designed to be used as an optional argument
// to the SVG function.It changes the default light and
// dark square colors to the colors given.
func SquareColors(light, dark color.Color) func(*encoder) {
	return func(e *encoder) {
		e.light = light
		e.dark = dark
	}
}

// MarkSquares is designed to be used as an optional argument
// to the SVG function.It marks the given squares with the
// color.A possible usage includes marking squares of the
// previous move.
func MarkSquares(c color.Color, sqs ...octad.Square) func(*encoder) {
	return func(e *encoder) {
		for _, sq := range sqs {
			e.marks[sq] = c
		}
	}
}

// A Encoder encodes octad boards into images.
type encoder struct {
	w     io.Writer
	light color.Color
	dark  color.Color
	marks map[octad.Square]color.Color
}

// New returns an encoder that writes to the given writer.
// New also takes options which can customize the image
// output.
func newEncoder(w io.Writer, options []func(*encoder)) *encoder {
	e := &encoder{
		w:     w,
		light: color.RGBA{R: 255, G: 255, B: 223, A: 1},
		dark:  color.RGBA{R: 147, G: 175, B: 108, A: 1},
		marks: map[octad.Square]color.Color{},
	}
	for _, op := range options {
		op(e)
	}
	return e
}

const (
	sqWidth     = 45
	sqHeight    = 45
	boardWidth  = 4 * sqWidth
	boardHeight = 4 * sqHeight
)

var (
	orderOfRanks = []octad.Rank{octad.Rank4, octad.Rank3, octad.Rank2, octad.Rank1}
	orderOfFiles = []octad.File{octad.FileA, octad.FileB, octad.FileC, octad.FileD}
)

// EncodeSVG writes the board SVG representation into
// the Encoder's writer.An error is returned if there
// is there is an error writing data.
func (e *encoder) EncodeSVG(b *octad.Board) error {
	boardMap := b.SquareMap()
	canvas := svg.New(e.w)
	canvas.Start(boardWidth, boardHeight)
	canvas.Rect(0, 0, boardWidth, boardHeight)

	for i := 0; i < 16; i++ {
		sq := octad.Square(i)
		x, y := xyForSquare(sq)
		// draw square
		c := e.colorForSquare(sq)
		canvas.Rect(x, y, sqWidth, sqHeight, "fill: "+colorToHex(c))
		markColor, ok := e.marks[sq]
		if ok {
			canvas.Rect(x, y, sqWidth, sqHeight, "fill-opacity:0.2;fill: "+colorToHex(markColor))
		}
		// draw piece
		p := boardMap[sq]
		if p != octad.NoPiece {
			xml := pieceXML(x, y, p)
			if _, err := io.WriteString(canvas.Writer, xml); err != nil {
				return err
			}
		}
		// draw rank text on file A
		txtColor := e.colorForText(sq)
		if sq.File() == octad.FileA {
			style := "font-size:9px;fill: " + colorToHex(txtColor)
			canvas.Text(x+(sqWidth*1/45), y+(sqHeight*5/24), sq.Rank().String(), style)
		}
		// draw file text on rank 1
		if sq.Rank() == octad.Rank1 {
			style := "text-anchor:end;font-size:9px;fill: " + colorToHex(txtColor)
			canvas.Text(x+(sqWidth*49/50), y+sqHeight-(sqHeight*1/40), sq.File().String(), style)
		}
	}
	canvas.End()
	return nil
}

func (e *encoder) colorForSquare(sq octad.Square) color.Color {
	sqSum := int(sq.File()) + int(sq.Rank())
	if sqSum%2 == 0 {
		return e.dark
	}
	return e.light
}

func (e *encoder) colorForText(sq octad.Square) color.Color {
	sqSum := int(sq.File()) + int(sq.Rank())
	if sqSum%2 == 0 {
		return e.light
	}
	return e.dark
}

func xyForSquare(sq octad.Square) (x, y int) {
	fileIndex := int(sq.File())
	rankIndex := 3 - int(sq.Rank())
	return fileIndex * sqWidth, rankIndex * sqHeight
}

func colorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", uint8(float64(r)+0.5), uint8(float64(g)*1.0+0.5), uint8(float64(b)*1.0+0.5))
}

func pieceXML(x, y int, p octad.Piece) string {
	fileName := fmt.Sprintf("pieces/%s%s.svg", p.Color().String(), pieceTypeMap[p.Type()])
	svgStr := string(internal.MustAsset(fileName))
	old := `<svg xmlns="http://www.w3.org/2000/svg" version="1.1" width="45" height="45">`
	newSvgStr := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" version="1.1" width="360" height="360" viewBox="%d %d 360 360">`, -1*x, -1*y)
	return strings.Replace(svgStr, old, newSvgStr, 1)
}

var (
	pieceTypeMap = map[octad.PieceType]string{
		octad.King:   "K",
		octad.Queen:  "Q",
		octad.Rook:   "R",
		octad.Bishop: "B",
		octad.Knight: "N",
		octad.Pawn:   "P",
	}
)
