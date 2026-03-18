package game

type Level struct {
	Title       string
	Description string
	Grid        [][]Tile
	StartX      int
	StartY      int
}

var level1 = Level{
	Title:       "Level 1 - Walk",
	Description: "use h, j, k, l to move, u to undo, q to quit",
	Grid: [][]Tile{
		{w, w, d2, w, w},
		{w, e, e, e, w},
		{w, e, e, e, w},
		{w, e, e, e, w},
		{w, w, w, w, w},
	},
	StartX: 2,
	StartY: 2,
}
var d2 = newDoor(level2)

var level2 = Level{
	Title:       "Level 2 - Push",
	Description: "Push box +1 count",
	Grid: [][]Tile{
		{w, w, w, w, w, w, w, w},
		{w, e, e, e, e, e, e, w},
		{w, e, e, e, b(), e, d3, w},
		{w, e, e, e, e, e, e, w},
		{w, w, w, w, w, w, w, w},
	},
	StartX: 2,
	StartY: 2,
}
var d3 = newLock(1, newDoor(level3))

var level3 = Level{
	Title:       "Level 3 - Pull",
	Description: "Hold shift to pull box, -1 count, minium 0",
	Grid: [][]Tile{
		{w, w, w, w, w, w, w, w, w, w},
		{w, e, e, e, e, e, e, e, e, w},
		{w, e, e, e, d4, e, e, b(), e, w},
		{w, e, e, e, e, e, e, e, e, w},
		{w, w, w, w, w, w, w, w, w, w},
	},
	StartX: 2,
	StartY: 2,
}
var d4 = newLock(1, newDoor(level4))

var level4 = Level{
	Title:       "Level 4 - Push Push",
	Description: "Combine 2 box",
	Grid: [][]Tile{
		{w, w, w, w, w, w, w, w},
		{w, e, b(), e, b(), e, d5, w},
		{w, e, e, e, e, e, e, w},
		{w, w, w, w, w, w, w, w},
	},
	StartX: 1,
	StartY: 1,
}
var d5 = newLock(4, newDoor(level5))

var level5 = Level{
	Title:       "Level 5 - Pull Push",
	Description: "x",
	Grid: [][]Tile{
		{w, w, w, w, w, w, w, w},
		{w, e, b(), e, e, b(), d6, w},
		{w, e, e, e, e, e, e, w},
		{w, w, w, w, w, w, w, w},
	},
	StartX: 1,
	StartY: 1,
}
var d6 = newLock(2, newDoor(level6))

var level6 = Level{
	Title:       "Level 6 - Charge",
	Description: "x",
	Grid: [][]Tile{
		{w, w, w, w, w, w, w, w, w, w, w, w, w},
		{w, e, e, e, e, e, e, e, e, e, b(), d7, w},
		{w, w, w, w, w, w, w, w, w, w, w, w, w},
	},
	StartX: 1,
	StartY: 1,
}
var d7 = newLock(8, newDoor(level7))

var level7 = Level{
	Title:       "Level 7 - Circle",
	Description: "x",
	Grid: [][]Tile{
		{w, w, w, w, w, w},
		{w, e, e, e, w, w},
		{w, e, d8, e, e, w},
		{w, e, e, b(), e, w, w, e, e},
		{w, w, e, e, w, w},
	},
	StartX: 1,
	StartY: 1,
}
var d8 = newLock(2, newDoor(level8))

var level8 = Level{
	Title:       "Level 8 - x",
	Description: "x",
	Grid: [][]Tile{
		{w, w, w, w, w},
		{w, e, e, w, w},
		{w, d9, b(), e, w},
		{w, b(), e, e, w, w, e, e},
		{w, e, e, w, w},
	},
	StartX: 1,
	StartY: 1,
}
var d9 = newLock(3, newDoor(winLevel))

var winLevel = Level{
	Title:       "Victory",
	Description: "You've reached the end! Well done.",
	Grid: [][]Tile{
		{w, w, w, w, w},
		{w, e, e, e, w},
		{w, e, e, e, w},
		{w, e, e, e, w},
		{w, w, w, w, w},
	},
	StartX: 2,
	StartY: 2,
}

var doorToWin = newDoor(winLevel)
