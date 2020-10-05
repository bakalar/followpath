package main

import (
	"bufio"
	"os"
	"strings"
)

// BeginningCharacter represents default character used for a beginning of a path
const BeginningCharacter = '@'

// EndingCharacter represents default character used for a path end
const EndingCharacter = 'x'

const spaceCharacter = ' '
const verticalConnectionCharacter = '|'
const horizontalConnectionCharacter = '-'
const edgeConnectionCharacter = '+'

type pathDirection int

const (
	none  pathDirection = 0
	up    pathDirection = 1
	right pathDirection = 2
	down  pathDirection = 3
	left  pathDirection = 4
)

// CharacterLocation represents character's location in a map
type CharacterLocation struct {
	lineIndex   int
	columnIndex int
}

// MapFromFileAsLines loads a map from file and returns it as a list of lines
func MapFromFileAsLines(file *os.File) []string {
	sc := bufio.NewScanner(file)

	// Characters valid in a map (not including letters)
	validCharacters := string([]byte{BeginningCharacter, EndingCharacter, spaceCharacter, verticalConnectionCharacter, horizontalConnectionCharacter, edgeConnectionCharacter})

	var lines []string
	for {
		if !sc.Scan() {
			break
		}
		mutableLine := []rune(sc.Text())
		for index := 0; index < len(mutableLine); index++ {
			character := mutableLine[index]
			if !isLetter(character) && !strings.ContainsRune(validCharacters, character) {
				mutableLine[index] = spaceCharacter
			}
		}
		lines = append(lines, string(mutableLine))
	}
	return lines
}

// GetLocation returns location of a character in a given map
func GetLocation(lines []string, character byte) CharacterLocation {
	for lineIndex, line := range lines {
		columnIndex := strings.IndexByte(line, character)
		if columnIndex != -1 {
			return CharacterLocation{lineIndex, columnIndex}
		}
	}
	panic("character " + string(character) + " not found")
}

// FollowPath follows a path on a map from beggining to end and returns it as a list of locations
func FollowPath(lines []string, start CharacterLocation, end CharacterLocation) []CharacterLocation {
	var direction pathDirection = none
	currentLocation := start
	path := []CharacterLocation{CharacterLocation{currentLocation.lineIndex, currentLocation.columnIndex}}
	for {
		if currentLocation == end {
			break
		}
		path = onePathStep(path, lines, direction, currentLocation)

		// next location is contained in the last element
		nextLocation := path[len(path)-1]

		hDiff := nextLocation.columnIndex - currentLocation.columnIndex
		vDiff := nextLocation.lineIndex - currentLocation.lineIndex
		direction = nextDirection(hDiff, vDiff)
		currentLocation = nextLocation
	}
	return path
}

// PathAsLetters converts a path to a list of letters
func PathAsLetters(lines []string, path []CharacterLocation) string {
	var letters string
	letterPath := []CharacterLocation{}
	for _, location := range path {
		character := lines[location.lineIndex][location.columnIndex]
		if isLetter(rune(character)) && !contains(letterPath, location) {
			letters += string(lines[location.lineIndex][location.columnIndex])
			letterPath = append(letterPath, location)
		}
	}
	return letters
}

// PathAsCharacters converts a path to a list of characters
func PathAsCharacters(lines []string, path []CharacterLocation) string {
	var characters string
	for _, location := range path {
		characters += string(lines[location.lineIndex][location.columnIndex])
	}
	return characters
}

func isLetter(character rune) bool {
	return character >= 'A' && character <= 'Z'
}

// max returns the larger of x or y.
func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// min returns the smaller of x or y.
func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func nextDirection(hDiff int, vDiff int) pathDirection {
	if hDiff > 0 {
		return right
	} else if hDiff < 0 {
		return left
	} else if vDiff > 0 {
		return down
	} else if vDiff < 0 {
		return up
	} else {
		panic("invalid path direction")
	}
}

func diffs(direction pathDirection) (int, int) {
	var hDiff int
	var vDiff int
	switch direction {
	case right:
		hDiff = 1
		vDiff = 0
	case down:
		hDiff = 0
		vDiff = 1
	case left:
		hDiff = -1
		vDiff = 0
	case up:
		hDiff = 0
		vDiff = -1
	case none:
		hDiff = 0
		vDiff = 0
	}
	return hDiff, vDiff
}

func onePathStep(path []CharacterLocation, lines []string, direction pathDirection, baseLocation CharacterLocation) []CharacterLocation {
	baseCharacter := lines[baseLocation.lineIndex][baseLocation.columnIndex]
	baseCharacterIsLetter := isLetter(rune(baseCharacter))
	baseCharacterIsVerticalConnection := baseCharacter == verticalConnectionCharacter
	baseCharacterIsHorizontalConnection := baseCharacter == horizontalConnectionCharacter
	hDiff, vDiff := diffs(direction)
	preferredNextLocation := CharacterLocation{baseLocation.lineIndex + vDiff, baseLocation.columnIndex + hDiff}

	// try all directions
	var validLocations []CharacterLocation
	for lineIndex := max(0, baseLocation.lineIndex-1); lineIndex <= min(baseLocation.lineIndex+1, len(lines)-1); lineIndex++ {
		line := lines[lineIndex]

		var columnIndexStart int
		var columnIndexEnd int
		if lineIndex == baseLocation.lineIndex {
			columnIndexStart = max(0, baseLocation.columnIndex-1)
			columnIndexEnd = min(baseLocation.columnIndex+1, len(line)-1)
		} else {
			// diagonal direction is not allowed
			columnIndexStart = max(0, baseLocation.columnIndex)
			columnIndexEnd = min(baseLocation.columnIndex, len(line)-1)
		}

		for columnIndex := columnIndexStart; columnIndex <= columnIndexEnd; columnIndex++ {
			if lineIndex == baseLocation.lineIndex-vDiff && columnIndex == baseLocation.columnIndex-hDiff {
				// can't go to backwards (to previous location)
				continue
			}
			currentLocation := CharacterLocation{lineIndex, columnIndex}
			if currentLocation == baseLocation {
				// this will not move from baseLocation
				continue
			}
			currentCharacter := line[columnIndex]
			if currentCharacter == spaceCharacter {
				// invalid direction
				continue
			}
			if baseCharacterIsLetter {
				if lineIndex == baseLocation.lineIndex {
					// vertical connection is not allowed here
					if currentCharacter != verticalConnectionCharacter {
						validLocations = append(validLocations, currentLocation)
					}
				}
				if columnIndex == baseLocation.columnIndex {
					// horizontal connection is not allowed here
					if currentCharacter != horizontalConnectionCharacter {
						validLocations = append(validLocations, currentLocation)
					}
				}
			} else if baseCharacterIsVerticalConnection {
				// only up and down directions are allowed
				if columnIndex == baseLocation.columnIndex {
					lineIndex2 := lineIndex
					for currentCharacter == horizontalConnectionCharacter {
						path = append(path, CharacterLocation{lineIndex2, columnIndex})
						if direction == down {
							lineIndex2++
						} else {
							lineIndex2--
						}
						if lineIndex2 < 0 || lineIndex2 >= len(lines) {
							break
						}
						currentCharacter = lines[lineIndex2][columnIndex]
					}
					if lineIndex2 != lineIndex && lineIndex2 < len(lines) {
						if preferredNextLocation == currentLocation {
							preferredNextLocation = CharacterLocation{lineIndex2, columnIndex}
						}
						currentLocation = CharacterLocation{lineIndex2, columnIndex}
					}
					validLocations = append(validLocations, currentLocation)
				}
			} else if baseCharacterIsHorizontalConnection {
				// only left and right directions are allowed
				if lineIndex == baseLocation.lineIndex {
					columnIndex2 := columnIndex
					for currentCharacter == verticalConnectionCharacter {
						path = append(path, CharacterLocation{lineIndex, columnIndex2})
						if direction == right {
							columnIndex2++
						} else {
							columnIndex2--
						}
						if columnIndex2 < 0 || columnIndex2 >= len(line) {
							break
						}
						currentCharacter = lines[lineIndex][columnIndex2]
					}
					if columnIndex2 != columnIndex && columnIndex2 < len(line) {
						if preferredNextLocation == currentLocation {
							preferredNextLocation = CharacterLocation{lineIndex, columnIndex2}
						}
						currentLocation = CharacterLocation{lineIndex, columnIndex2}
					}
					validLocations = append(validLocations, currentLocation)
				}
			} else {
				validLocations = append(validLocations, currentLocation)
			}
		}
	}
	if preferredNextLocation != baseLocation {
		for _, validLocation := range validLocations {
			if validLocation == preferredNextLocation {
				return append(path, preferredNextLocation)
			}
		}
	}

	if len(validLocations) == 0 {
		panic("no direction to follow")
	}
	return append(path, validLocations[0])
}

// Is location containt in this path?
func contains(path []CharacterLocation, location CharacterLocation) bool {
	for _, location1 := range path {
		if location1 == location {
			return true
		}
	}
	return false
}
