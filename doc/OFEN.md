# Octad Forsyth-Edwards Notation (OFEN)

Forsyth–Edwards Notation (FEN) is a standard notation for describing a
particular board position of a chess game. The purpose of FEN is to provide all
the necessary information to restart a game from a particular position.

OFEN is a derivation of FEN to support the features of Octad.

## OFEN Structure

An OFEN record contains six fields. The separator between fields is a space.
The fields are:

1. Piece placement (from White's perspective). Each rank is described, starting
   with rank 4 and ending with rank 1; within each rank, the contents of each
   square are described from file "a" through file "d". Following the Standard
   Algebraic Notation (SAN), each piece is identified by a single letter taken
   from the standard English names (pawn = "P", knight = "N", bishop = "B",
   rook = "R", queen = "Q" and king = "K"). White pieces are designated using
   upper-case letters ("PNBRQK") while black pieces use lowercase ("pnbrqk").
   Empty squares are noted using digits 1 through 4 (the number of empty
   squares), and "/" separates ranks.

2. Active color. "w" means White moves next, "b" means Black moves next.

3. Castling availability. If neither color can castle, this is "-". Otherwise,
   this has one or more letters: "N" (White can castle knightside), "C" (White
   can castle with close pawn), "F" (White can castle with far pawn),  "n"
   (Black can castle knightside), "c" (Black can castle with close pawn),
   and/or "f" (Black can castle with far pawn).

4. En passant target square in algebraic notation. If there's no en passant
   target square, this is "-". If a pawn has just made a two-square move, this
   is the position "behind" the pawn. This is recorded regardless of whether
   there is a pawn in position to make an en passant capture.

5. Halfmove clock: This is the number of halfmoves since the last capture or
   pawn advance. This is used to determine if a draw can be claimed under the
   twenty-move rule.

6. Fullmove number: The number of the full move. It starts at 1, and is
   incremented after Black's move.

Here is the OFEN for the starting position:
```ppkn/4/4/NKPP w NCFncf - 0 1```

Here is the OFEN after the move 1. c2:
```ppkn/4/2P1/NK1P b NCFncf - 0 1```

*Adapted from https://en.wikipedia.org/wiki/Forsyth-Edwards_Notation*