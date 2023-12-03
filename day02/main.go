package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// cubeSet represents a set of red, green, and/or blue cubes.
// It may be either be the bag containing all cubes,
// or a collection of cubes pulled from a bag.
type cubeSet struct {
	red   int
	green int
	blue  int
}

// game represents a single game, consisting of an ID and a slice of cubeSets
type game struct {
	id     int
	rounds []cubeSet
}

// newCubeSet returns a pointer to a new cubeSet object
func newCubeSet(red, green, blue int) *cubeSet {
	cubeBag := cubeSet{red: red, green: green, blue: blue}
	return &cubeBag
}

// newGame parses a line into a game object
func newGame(line string) *game {
	lineElements := strings.Split(line, ":")
	if len(lineElements) < 2 {
		log.Fatalf("Invalid game found: %v", line)
	}

	gameNumString := lineElements[0]
	gameResultsString := lineElements[1]

	digitPattern := regexp.MustCompile(`\d+\b`)
	colorPattern := regexp.MustCompile("red|green|blue")

	id, err := strconv.Atoi(digitPattern.FindString(gameNumString))
	if err != nil {
		log.Fatal(err)
	}

	gameElements := strings.Split(gameResultsString, ";")
	thisGame := game{id: id}

	for _, v := range gameElements {
		cubeColors := strings.Split(v, ",")
		thisCubeSet := cubeSet{}

		for _, cubeColor := range cubeColors {
			quantity, err := strconv.Atoi(digitPattern.FindString(cubeColor))
			if err != nil {
				log.Fatal(err)
			}

			color := colorPattern.FindString(cubeColor)

			switch color {
			case "red":
				thisCubeSet.red = quantity
			case "green":
				thisCubeSet.green = quantity
			case "blue":
				thisCubeSet.blue = quantity
			}
		}

		thisGame.rounds = append(thisGame.rounds, thisCubeSet)
	}

	return &thisGame
}

// canContain checks whether a given cubeSet can contain all the cubes of another given cubeSet
func (cubeBag *cubeSet) canContain(containedCubeSet cubeSet) bool {
	return cubeBag.red >= containedCubeSet.red &&
		cubeBag.green >= containedCubeSet.green &&
		cubeBag.blue >= containedCubeSet.blue
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please include the filepath as an argument.")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	idTotal := 0
	scanner := bufio.NewScanner(file)
	mainCubeSet := newCubeSet(12, 13, 14)

	for scanner.Scan() {
		currentGame := newGame(scanner.Text())
		gameIsValid := true

		for _, v := range currentGame.rounds {
			if !mainCubeSet.canContain(v) {
				gameIsValid = false
			}
		}

		if gameIsValid {
			idTotal += currentGame.id
		}
	}

	fmt.Printf("Total of valid game IDs: %v\n", idTotal)
}
