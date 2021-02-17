# octad
Standard library for the Octad chess variant written in Go.

Octad is a creation of Andrew DeChristopher. Rules and information about the
game can be found below. Octad is thought to be a deterministic game but needs
formal verification to prove that.

## Octad Game
### Board Layout
Each player begins with four pieces, a knight, their king, and two pawns placed
in that order. An example of this can be seen in the board diagram below:

![Octad board](doc/octad.png "Octad board")

### Rules
All standard chess rules apply:

* En Passant is possible
* Pawn promotion to any piece
* Stalemates are a draw

The only catch, however, is that castling is possible between the king and any
of its pieces on the starting rank before movement. The king will simply switch
spaces with the castling piece in all cases except the far pawn, in which case
the king will travel one space to the right and the pawn will lie where the king
was before. An example of white castling with their d pawn can be expressed as
`[ 1. c2 b3 2. O-O-O ... ]` with the resulting structure leaving the knight on
a1, a pawn on b1, the king on c1, and the other pawn on c2.

#### Castling notation
* Knight-color castle: **O**
* Close pawn castle: **O-O**
* Far pawn castle: **O-O-O**