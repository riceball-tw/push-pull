[2026-03-15 11:17] - Updated by Junie
{
    "TYPE": "positive",
    "CATEGORY": "praise and enhancement",
    "EXPECTATION": "User is happy with single-wall collision but now wants walls represented as a 2D array to support multiple walls.",
    "NEW INSTRUCTION": "WHEN modeling walls in this game THEN use a 2D array and update collision/rendering"
}

[2026-03-15 11:54] - Updated by Junie
{
    "TYPE": "correction",
    "CATEGORY": "grid data model",
    "EXPECTATION": "Grid should support multiple element types; do not use bool for walls.",
    "NEW INSTRUCTION": "WHEN designing the grid data structure THEN use a typed cell enum, not bool"
}

[2026-03-15 12:40] - Updated by Junie
{
    "TYPE": "correction",
    "CATEGORY": "runtime panic bug",
    "EXPECTATION": "After entering a door, movement must use the new map's dimensions and never index outside the grid.",
    "NEW INSTRUCTION": "WHEN accessing grid cells after map change THEN bounds-check using current map width/height"
}

[2026-03-15 12:58] - Updated by Junie
{
    "TYPE": "correction",
    "CATEGORY": "insufficient detail",
    "EXPECTATION": "User wants a more detailed, step-by-step explanation and examples for the requested change.",
    "NEW INSTRUCTION": "WHEN explaining how to modify tiles or assets THEN provide step-by-step instructions with code snippets and examples"
}

[2026-03-15 13:29] - Updated by Junie
{
    "TYPE": "correction",
    "CATEGORY": "sound playback",
    "EXPECTATION": "User wants actual audio to play on tile entry instead of showing a message.",
    "NEW INSTRUCTION": "WHEN entering a tile that has sound THEN play the associated audio file, not a message"
}

[2026-03-15 13:31] - Updated by Junie
{
    "TYPE": "correction",
    "CATEGORY": "sound playback",
    "EXPECTATION": "User wants actual audio to play on tile entry instead of showing a message.",
    "NEW INSTRUCTION": "WHEN entering a tile with a sound THEN trigger audio playback, avoid TUI message"
}

[2026-03-15 13:44] - Updated by Junie
{
    "TYPE": "correction",
    "CATEGORY": "build/run instructions",
    "EXPECTATION": "After refactoring into multiple files, the program should run without undefined symbol errors; instructions should indicate how to run all files.",
    "NEW INSTRUCTION": "WHEN refactoring into multiple files THEN instruct to use `go run .` or include all files"
}

