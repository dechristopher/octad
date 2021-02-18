package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dechristopher/octad"
)

func main() {
	fmt.Print("\033[H\033[2J")

	for {
		var err error
		var game *octad.Game

		var position, tm, rights, enp, hm, m string
		fmt.Println("\nLi(bre oct)ad Test Harness v0.0.1")
		fmt.Printf("Enter OFEN or press enter for new game:")
		_, err = fmt.Scanf("%s %s %s %s %s %s",
			&position, &tm, &rights, &enp, &hm, &m)
		if err != nil {
			game, err = octad.NewGame()
		} else {
			ofen := fmt.Sprintf("%s %s %s %s %s %s",
				position, tm, rights, enp, hm, m)
			ofen = strings.TrimSpace(ofen)
			parsedOFEN, err := octad.OFEN(ofen)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			game, err = octad.NewGame(parsedOFEN)
		}

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Print("\033[H\033[2J")

		for game.Outcome() == octad.NoOutcome {
			fmt.Print("\033[H\033[2J")

			fmt.Println(game.OFEN())
			if len(game.Moves()) > 0 {
				fmt.Printf("last move: %s", game.Moves()[len(game.Moves())-1].String())
				if game.Position().EnPassantSquare() != octad.NoSquare {
					fmt.Printf(", enp: %s", game.Position().EnPassantSquare().String())
				}
			}
			fmt.Println(game.Position().Board().Draw())

			// Print all legal moves
			for _, move := range game.ValidMoves() {
				fmt.Print(move.String() + " ")
			}

			var move string

			fmt.Printf("\nEnter move for %s: ", game.Position().Turn().String())
			_, err := fmt.Scanln(&move)
			if err != nil {
				continue
			}

			if move == "-" {
				break
			} else {
				err = game.MoveStr(move)
			}

			if err != nil {
				fmt.Println(err.Error())
			}
		}

		fmt.Println(game.Position().Board().Draw())
		fmt.Printf("Game completed. %s by %s.\n", game.Outcome(), game.Method())
		fmt.Println(game.String())

		//// create file
		//f, err := os.Create("test.svg")
		//if err != nil {
		//	log.Fatal(err)
		//}
		//defer f.Close()
		//
		//// write board SVG to file
		//blue := color.RGBA{R: 0, G: 64, B: 255, A: 1}
		//mark := image.MarkSquares(blue, octad.C1, octad.C2)
		//if err := image.SVG(f, game.Position().Board(), mark); err != nil {
		//	log.Fatal(err)
		//}
	}
}
