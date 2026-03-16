package main

type Level struct {
	Title       string
	Description string
	Grid        [][]Tile
	StartX      int
	StartY      int
}

var level1 = Level{
	Title:       "Level 1",
	Description: "use h, j, k, l to move, u to undo, q to quit",
	Grid: [][]Tile{
		{wall, wall, doorToLevel2, wall, wall},
		{wall, empty, empty, empty, wall},
		{wall, empty, empty, empty, wall},
		{wall, empty, empty, empty, wall},
		{wall, wall, wall, wall, wall},
	},
	StartX: 2,
	StartY: 2,
}
var doorToLevel2 = newDoor(level2)


var level2 = Level{
	Title:       "Level 2",
	Description: "It's a bit empty here.",
	Grid: [][]Tile{
		{empty},

	},
	StartX: 0,
	StartY: 0,
}

var winLevel = Level{
	Title:       "Victory",
	Description: "You've reached the end! Well done.",
	Grid: [][]Tile{
		{wall, wall, wall, wall, wall},
		{wall, empty, empty, empty, wall},
		{wall, empty, empty, empty, wall},
		{wall, empty, empty, empty, wall},
		{wall, wall, wall, wall, wall},
	},
	StartX: 2,
	StartY: 2,
}

var doorToWin = newDoor(winLevel)
