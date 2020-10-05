package main

import "os"

func main() {
	lines := MapFromFileAsLines(os.Stdin)

	startPos := GetLocation(lines, BeginningCharacter)
	endPos := GetLocation(lines, EndingCharacter)

	path := FollowPath(lines, startPos, endPos)
	println("Letters", PathAsLetters(lines, path))
	println("Path as characters", PathAsCharacters(lines, path))
}
