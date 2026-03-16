package main

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
		{w, e, e, e, newBox(), e, d3, w},
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
		{w, e, e, e, d4, e, e, newBox(), e, w},
		{w, e, e, e, e, e, e, e, e, w},
		{w, w, w, w, w, w, w, w, w, w},

	},
	StartX: 2,
	StartY: 2,
}
var d4 = newLock(1, newDoor(level4))


var level4 = Level{
	Title:       "Level 4 - Combine",
	Description: "Combine 2 box",
	Grid: [][]Tile{
		{w, w, w, w, w, w, w, w},
		{w, e, newBox(), e, newBox(), e, d5, w},
		{w, e, e, e, e, e, e, w},
		{w, w, w, w, w, w, w, w},
	},
	StartX: 1,
	StartY: 1,
}
var d5 = newLock(4, newDoor(winLevel))

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
