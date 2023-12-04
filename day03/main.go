package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type partNum struct {
	id string
	x  int
	y  int
}

type coordinate struct {
	x     int
	y     int
	value string
}

func (num partNum) isAdjacentToSymbol(schematic []string) bool {
	var allAdjacentRunes []coordinate

	for i := range num.id {
		currentCharXCoord := num.x + i
		currentCharYCoord := num.y
		var currentAdjacentRunes []coordinate

		currentAdjacentRunes = append(currentAdjacentRunes, coordinate{x: currentCharXCoord + 1, y: currentCharYCoord})
		currentAdjacentRunes = append(currentAdjacentRunes, coordinate{x: currentCharXCoord - 1, y: currentCharYCoord})
		currentAdjacentRunes = append(currentAdjacentRunes, coordinate{x: currentCharXCoord, y: currentCharYCoord + 1})
		currentAdjacentRunes = append(currentAdjacentRunes, coordinate{x: currentCharXCoord, y: currentCharYCoord - 1})
		currentAdjacentRunes = append(currentAdjacentRunes, coordinate{x: currentCharXCoord + 1, y: currentCharYCoord + 1})
		currentAdjacentRunes = append(currentAdjacentRunes, coordinate{x: currentCharXCoord + 1, y: currentCharYCoord - 1})
		currentAdjacentRunes = append(currentAdjacentRunes, coordinate{x: currentCharXCoord - 1, y: currentCharYCoord + 1})
		currentAdjacentRunes = append(currentAdjacentRunes, coordinate{x: currentCharXCoord - 1, y: currentCharYCoord - 1})

		for _, adjacentElement := range currentAdjacentRunes {
			if adjacentElement.x >= 0 && adjacentElement.x < len(schematic[num.y]) && adjacentElement.y >= 0 && adjacentElement.y < len(schematic) {
				adjacentElement.value = string(schematic[adjacentElement.y][adjacentElement.x])
				allAdjacentRunes = append(allAdjacentRunes, adjacentElement)
			}
		}
	}

	hasAdjacentSymbol := false
	symbolPattern := regexp.MustCompile(`[^a-zA-Z\d\.]`)
	for _, v := range allAdjacentRunes {
		if symbolPattern.MatchString(v.value) {
			hasAdjacentSymbol = true
			break
		}
	}

	return hasAdjacentSymbol
}

func main() {
	filePath := ""
	if len(os.Args) < 2 {
		log.Fatal("Please include the input file as an argment.")
	}

	filePath = os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var schematic []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		schematic = append(schematic, line)
	}

	total := 0
	partNumPattern := regexp.MustCompile(`\d+\b`)

	for lineIndex, line := range schematic {

		currentLinePartNums := partNumPattern.FindAllString(line, -1)
		currentLineIndices := partNumPattern.FindAllStringIndex(line, -1)

		for partNumIndex, v := range currentLinePartNums {
			startCol := currentLineIndices[partNumIndex][0]
			currentNumber := partNum{id: v, x: startCol, y: lineIndex}

			if currentNumber.isAdjacentToSymbol(schematic) {
				idNum, err := strconv.Atoi(currentNumber.id)
				if err != nil {
					log.Fatal(err)
				}

				total += idNum
			}
		}
	}

	fmt.Printf("The total is: %v\n", total)
}
