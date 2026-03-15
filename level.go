package main

type Level struct {
	Grid   [][]Tile
	StartX int
	StartY int
}

var level1 = Level{
	Grid: [][]Tile{
		{empty, newBox(), empty, empty, empty, newBox(), empty, empty},
	},
	StartX: 3,
	StartY: 0,
}
