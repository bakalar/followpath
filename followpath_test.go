package main

import (
	"os"
	"testing"
)

func TestMap1(t *testing.T) {
	testMapFile(t, "map1.txt", "ACB", "@---A---+|C|+---+|+-B-x")
}

func TestMap2(t *testing.T) {
	testMapFile(t, "map2.txt", "ABCD", "@|A+---B--+|+----C|-||+---D--+|x")
}

func TestMap3(t *testing.T) {
	testMapFile(t, "map3.txt", "BEEFCAKE", "@---+B||E--+|E|+--F--+|C|||A--|-----K|||+--E--Ex")
}

func TestMapInfiniteLoop(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	testMapFile(t, "map_infinite_loop.txt", "", "")
}

func TestMapNoEnd(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	testMapFile(t, "map_no_end.txt", "", "")
}

func TestMapDoubleCrossing(t *testing.T) {
	testMapFile(t, "map_double_crossing.txt", "BCDEFGH", "@---+B||C--+||D|++|E+----+|||F--|-----G||||+--DE---H-x")
}

func testMapFile(t *testing.T, filename string, expectedLetters string, expectedCharacters string) {
	lines, path := followPathFromFile(t, filename)
	if PathAsLetters(lines, path) != expectedLetters {
		t.Fail()
	}
	if PathAsCharacters(lines, path) != expectedCharacters {
		t.Fail()
	}
}

func followPathFromFile(t *testing.T, filename string) ([]string, []CharacterLocation) {
	file, err := os.Open(filename)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	lines := MapFromFileAsLines(file)

	startPos := GetLocation(lines, BeginningCharacter)
	endPos := GetLocation(lines, EndingCharacter)

	path := FollowPath(lines, startPos, endPos)
	return lines, path
}
