package main

type Level struct {
	Grid   [][]Tile
	StartX int
	StartY int
}

var level1 = Level{
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
var doorToLevel2 = newDoor(level2.Grid, level2.StartX, level2.StartY)


var level2 = Level{
	Grid: [][]Tile{
		{empty},

	},
	StartX: 0,
	StartY: 0,
}

var winLevel = Level{
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

var doorToWin = newDoor(winLevel.Grid, winLevel.StartX, winLevel.StartY)
