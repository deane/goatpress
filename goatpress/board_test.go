package goatpress

import (
	"testing"
)

func TestMakingNewBoards(t *testing.T) {
	bg := defaultBoardGenerator()
	board := bg.newBoard(5)
	if board.Size != 5 {
		t.Errorf("board.Size is not right", board.Size, 5)
	}
	if len(board.Letters) != 5 {
		t.Errorf("board.Letters is not right size", len(board.Letters), 5)
	}
	if len(board.Letters[0]) != 5 {
		t.Errorf("board.Letters[0] is not right size", len(board.Letters[0]), 5)
	}
	if board.Letters[0][0] == "" {
		t.Errorf("board.Letters hasn't been filled in")
	}
}

func TestBoard(t *testing.T) {
	bg := defaultBoardGenerator()
	board := bg.newBoard(5)
	for i := 0; i < 25; i++ {
		board.Letters[i/5][i%5] = "q"
	}

	board.Letters[0][0] = "h"
	board.Letters[0][1] = "e"
	board.Letters[0][2] = "l"
	board.Letters[0][3] = "l"
	board.Letters[0][4] = "o"

	word := board.WordFromTiles([]Tile{newTile(0, 0), newTile(0, 1), newTile(0, 2), newTile(0, 3), newTile(0, 4)})
	if word != "hello" {
		t.Errorf("board.WordFromMove", word, "hello")
	}
	letters := board.HasLetters()

	if !letters["h"] {
		t.Errorf("board.HasLetters does not include 'h'")
	}
	if letters["x"] {
		t.Errorf("board.HasLetters does include 'x'")
	}

	if len(board.RandomMoveFromWord("hell").Tiles) != 4 {
		t.Errorf("board.MoveFromWord('hell') didn't return four tiles")
	}
	if board.RandomMoveFromWord("hell").Tiles[1][1] != 1 {
		t.Errorf("board.MoveFromWord('hell') didn't return correct move")
	}

	tilesA := board.TilesForLetterExcluding("a", []Tile{})
	if len(tilesA) != 0 {
		t.Errorf("board.TilesForLetterExcluding('a') shouldn't have returned anything")
	}

	tilesH := board.TilesForLetterExcluding("h", []Tile{})
	if len(tilesH) != 1 {
		t.Errorf("board.TilesForLetterExcluding('h') should be one tile", len(tilesH), 1)
	}

	tilesL := board.TilesForLetterExcluding("l", []Tile{})
	if len(tilesL) != 2 {
		t.Errorf("board.TilesForLetterExcluding('l') should be two tiles", len(tilesL), 2)
	}
	if tilesL[0][0] != 0 || tilesL[0][1] != 2 {
		t.Errorf("first l tile is not as expected")
	}
	if tilesL[1][0] != 0 || tilesL[1][1] != 3 {
		t.Errorf("second l tile is not as expected")
	}

	tilesL2 := board.TilesForLetterExcluding("l", []Tile{newTile(0, 2)})
	if len(tilesL2) != 1 {
		t.Errorf("board.TilesForLetterExcluding('l', 0-2) should be one tiles", len(tilesL2), 1)
	}
	if tilesL2[0][0] != 0 || tilesL2[0][1] != 3 {
		t.Errorf("first excluded l tile is not as expected", tilesL2[0][1], 3)
	}

}

func TestColorMask(t *testing.T) {
	bg := defaultBoardGenerator()
	board := bg.newBoard(5)
	SetupBoard(board)
	colorMask := newColorMask(board,
		[]Move{
			board.MoveFromTiles([]Tile{newTile(0, 0), newTile(1, 2)}),
			board.MoveFromTiles([]Tile{newTile(3, 2), newTile(1, 2)})})
	expColorMaskString := "10000 00200 00000 00200 00000"
	if colorMask.ToString() != expColorMaskString {
		t.Errorf("color mask is wrong", colorMask.ToString(), expColorMaskString)
	}
	if colorMask.Score(1) != 1 {
		t.Errorf("wrong score for player 1", colorMask.Score(1), 1)
	}
	if colorMask.Score(2) != 2 {
		t.Errorf("wrong score for player 2", colorMask.Score(2), 2)
	}
}

func TestColorMaskDarkTaking(t *testing.T) {
	bg := defaultBoardGenerator()
	board := bg.newBoard(5)
	SetupBoard(board)
	colorMask := newColorMask(board,
		[]Move{
			board.MoveFromTiles([]Tile{newTile(0, 0), newTile(0, 1), newTile(1, 0)})})
	expColorMaskString := "11000 10000 00000 00000 00000"
	if colorMask.ToString() != expColorMaskString {
		t.Errorf("color mask is wrong", colorMask.ToString(), expColorMaskString)
	}
	if colorMask.Score(1) != 3 {
		t.Errorf("wrong score for player 1", colorMask.Score(1), 3)
	}
	if colorMask.Score(2) != 0 {
		t.Errorf("wrong score for player 2", colorMask.Score(2), 0)
	}

	colorMask = newColorMask(board,
		[]Move{
			board.MoveFromTiles([]Tile{newTile(0, 0), newTile(0, 1), newTile(1, 0)}),
			board.MoveFromTiles([]Tile{newTile(0, 0), newTile(0, 1), newTile(1, 0), newTile(1, 1)})})
	expColorMaskString = "12000 22000 00000 00000 00000"
	if colorMask.ToString() != expColorMaskString {
		t.Errorf("color mask is wrong", colorMask.ToString(), expColorMaskString)
	}
	if colorMask.Score(1) != 1 {
		t.Errorf("wrong score for player 1", colorMask.Score(1), 1)
	}
	if colorMask.Score(2) != 3 {
		t.Errorf("wrong score for player 2", colorMask.Score(2), 3)
	}
}

func SetupBoard(board *Board) {
	board.Letters[0][0] = "h"
	board.Letters[0][1] = "e"
	board.Letters[0][2] = "l"
	board.Letters[0][3] = "l"
	board.Letters[0][4] = "o"

	board.Letters[1][0] = "s"
	board.Letters[1][1] = "t"
	board.Letters[1][2] = "a"
	board.Letters[1][3] = "t"
	board.Letters[1][4] = "e"

	board.Letters[2][0] = "j"
	board.Letters[2][1] = "e"
	board.Letters[2][2] = "n"
	board.Letters[2][3] = "g"
	board.Letters[2][4] = "a"

	board.Letters[3][0] = "p"
	board.Letters[3][1] = "a"
	board.Letters[3][2] = "g"
	board.Letters[3][3] = "e"
	board.Letters[3][4] = "s"

	board.Letters[4][0] = "v"
	board.Letters[4][1] = "a"
	board.Letters[4][2] = "l"
	board.Letters[4][3] = "i"
	board.Letters[4][4] = "d"
}
