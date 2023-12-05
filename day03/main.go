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
	var allAdjacentValues []coordinate

	for i := range num.id {
		currentAdjacentValues := getAdjacentCoordinates(coordinate{x: (num.x + i), y: num.y}, schematic)
		allAdjacentValues = append(allAdjacentValues, currentAdjacentValues...)
	}

	hasAdjacentSymbol := false
	symbolPattern := regexp.MustCompile(`[^a-zA-Z\d\.]`)
	for _, v := range allAdjacentValues {
		if symbolPattern.MatchString(v.value) {
			hasAdjacentSymbol = true
			break
		}
	}

	return hasAdjacentSymbol
}

func getAdjacentCoordinates(element coordinate, schematic []string) []coordinate {
	elementXCoord := element.x
	elementYCoord := element.y
	var allAdjacentValues []coordinate
	var currentAdjacentCoordinates []coordinate

	currentAdjacentCoordinates = append(currentAdjacentCoordinates, coordinate{x: elementXCoord + 1, y: elementYCoord})
	currentAdjacentCoordinates = append(currentAdjacentCoordinates, coordinate{x: elementXCoord - 1, y: elementYCoord})
	currentAdjacentCoordinates = append(currentAdjacentCoordinates, coordinate{x: elementXCoord, y: elementYCoord + 1})
	currentAdjacentCoordinates = append(currentAdjacentCoordinates, coordinate{x: elementXCoord, y: elementYCoord - 1})
	currentAdjacentCoordinates = append(currentAdjacentCoordinates, coordinate{x: elementXCoord + 1, y: elementYCoord + 1})
	currentAdjacentCoordinates = append(currentAdjacentCoordinates, coordinate{x: elementXCoord + 1, y: elementYCoord - 1})
	currentAdjacentCoordinates = append(currentAdjacentCoordinates, coordinate{x: elementXCoord - 1, y: elementYCoord + 1})
	currentAdjacentCoordinates = append(currentAdjacentCoordinates, coordinate{x: elementXCoord - 1, y: elementYCoord - 1})

	for _, adjacentElement := range currentAdjacentCoordinates {
		if adjacentElement.x >= 0 && adjacentElement.x < len(schematic[elementYCoord]) && adjacentElement.y >= 0 && adjacentElement.y < len(schematic) {
			adjacentElement.value = string(schematic[adjacentElement.y][adjacentElement.x])
			allAdjacentValues = append(allAdjacentValues, adjacentElement)
		}
	}

	return allAdjacentValues
}

func getPartNumTotal(schematic []string) int {
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

	return total
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

	total := getPartNumTotal(schematic)

	fmt.Printf("The total is: %v\n", total)
}
