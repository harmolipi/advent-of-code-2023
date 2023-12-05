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
	valid bool
}

// isAdjacentToSymbol checks whether a given partNum is adjacent to a symbol
func (num partNum) isAdjacentToSymbol(schematic []string) bool {
	var allAdjacentValues []coordinate

	for i := range num.id {
		currentAdjacentValues := getAdjacentCoordinates(coordinate{x: (num.x + i), y: num.y}, schematic)

		for _, v := range currentAdjacentValues {
			allAdjacentValues = append(allAdjacentValues, v)
		}

	}

	hasAdjacentSymbol := false
	symbolPattern := regexp.MustCompile(`[^a-zA-Z\d\.]`)
	for _, v := range allAdjacentValues {
		if !v.valid {
			continue
		}

		if symbolPattern.MatchString(v.value) {
			hasAdjacentSymbol = true
			break
		}
	}

	return hasAdjacentSymbol
}

// isAdjacentToGearRatio checks whether a given gear part * is adjacent
// to exactly 2 partnums, and if so, returns true and both of them
// in a slice.
func (num partNum) isAdjacentToGearRatio(schematic []string) (bool, []string) {
	singleDigitPattern := regexp.MustCompile(`\d`)

	numAdjacentGearRatios := 0
	topRowHasRatios := false
	middleRowHasRatios := false
	bottomRowHasRatios := false

	allAdjacentValues := getAdjacentCoordinates(coordinate{x: num.x, y: num.y}, schematic)
	numAdjacentDigits := 0
	hasAdjacentGearRatio := false

	for _, v := range allAdjacentValues {
		if !v.valid {
			continue
		}

		if singleDigitPattern.MatchString(v.value) {
			numAdjacentDigits++
		}
	}

	if allAdjacentValues["topLeftCoordinate"].valid && singleDigitPattern.MatchString(allAdjacentValues["topLeftCoordinate"].value) {
		numAdjacentGearRatios++
		topRowHasRatios = true

		if allAdjacentValues["topMiddleCoordinate"].valid && !singleDigitPattern.MatchString(allAdjacentValues["topMiddleCoordinate"].value) &&
			(allAdjacentValues["topRightCoordinate"].valid && singleDigitPattern.MatchString(allAdjacentValues["topRightCoordinate"].value)) {
			numAdjacentGearRatios++
		}

	} else if allAdjacentValues["topRightCoordinate"].valid && singleDigitPattern.MatchString(allAdjacentValues["topRightCoordinate"].value) {
		numAdjacentGearRatios++
		topRowHasRatios = true
	} else if allAdjacentValues["topMiddleCoordinate"].valid && singleDigitPattern.MatchString(allAdjacentValues["topMiddleCoordinate"].value) {
		numAdjacentGearRatios++
		topRowHasRatios = true
	}

	if allAdjacentValues["leftCoordinate"].valid && singleDigitPattern.MatchString(allAdjacentValues["leftCoordinate"].value) {
		numAdjacentGearRatios++
		middleRowHasRatios = true
	}

	if allAdjacentValues["rightCoordinate"].valid && singleDigitPattern.MatchString(allAdjacentValues["rightCoordinate"].value) {
		numAdjacentGearRatios++
		middleRowHasRatios = true
	}

	if allAdjacentValues["bottomLeftCoordinate"].valid && singleDigitPattern.MatchString(allAdjacentValues["bottomLeftCoordinate"].value) {
		numAdjacentGearRatios++
		bottomRowHasRatios = true

		if allAdjacentValues["bottomMiddleCoordinate"].valid && !singleDigitPattern.MatchString(allAdjacentValues["bottomMiddleCoordinate"].value) &&
			(allAdjacentValues["bottomRightCoordinate"].valid && singleDigitPattern.MatchString(allAdjacentValues["bottomRightCoordinate"].value)) {
			numAdjacentGearRatios++
		}

	} else if allAdjacentValues["bottomRightCoordinate"].valid && singleDigitPattern.MatchString(allAdjacentValues["bottomRightCoordinate"].value) {
		numAdjacentGearRatios++
		bottomRowHasRatios = true
	} else if allAdjacentValues["bottomMiddleCoordinate"].valid && singleDigitPattern.MatchString(allAdjacentValues["bottomMiddleCoordinate"].value) {
		numAdjacentGearRatios++
		bottomRowHasRatios = true
	}

	hasAdjacentGearRatio = numAdjacentGearRatios == 2
	var partNumPatterns []string

	if hasAdjacentGearRatio {
		partNumPattern := regexp.MustCompile(`\d+`)

		if topRowHasRatios {
			partNumLocations := partNumPattern.FindAllStringIndex(schematic[num.y-1], -1)

			for _, v := range partNumLocations {
				if (v[0] >= (num.x-1) && v[0] <= (num.x+1)) || (v[1] > (num.x-1) && v[1] <= (num.x+2) || (v[0] < (num.x-1) && v[1] > (num.x+2))) {
					partNumPatterns = append(partNumPatterns, schematic[num.y-1][v[0]:v[1]])
				}
			}
		}

		if middleRowHasRatios {
			partNumLocations := partNumPattern.FindAllStringIndex(schematic[num.y], -1)

			for _, v := range partNumLocations {
				if (v[0] >= (num.x-1) && v[0] <= (num.x+1)) || (v[1] > (num.x-1) && v[1] <= (num.x+2) || (v[0] < (num.x-1) && v[1] > (num.x+2))) {
					partNumPatterns = append(partNumPatterns, schematic[num.y][v[0]:v[1]])
				}
			}
		}

		if bottomRowHasRatios {
			partNumLocations := partNumPattern.FindAllStringIndex(schematic[num.y+1], -1)

			for _, v := range partNumLocations {
				if (v[0] >= (num.x-1) && v[0] <= (num.x+1)) || (v[1] > (num.x-1) && v[1] <= (num.x+2) || (v[0] < (num.x-1) && v[1] > (num.x+2))) {
					partNumPatterns = append(partNumPatterns, schematic[num.y+1][v[0]:v[1]])
				}
			}
		}
	}

	return hasAdjacentGearRatio, partNumPatterns
}

// getAdjacentCoordinates assembles adjacent coordinates into a map, and indicates whether each
// one is valid, and what its value is.
func getAdjacentCoordinates(element coordinate, schematic []string) map[string]coordinate {
	elementXCoord := element.x
	elementYCoord := element.y
	adjacentCoordinates := make(map[string]coordinate)

	adjacentCoordinates["topLeftCoordinate"] = coordinate{x: elementXCoord - 1, y: elementYCoord - 1}
	adjacentCoordinates["topMiddleCoordinate"] = coordinate{x: elementXCoord, y: elementYCoord - 1}
	adjacentCoordinates["topRightCoordinate"] = coordinate{x: elementXCoord + 1, y: elementYCoord - 1}
	adjacentCoordinates["leftCoordinate"] = coordinate{x: elementXCoord - 1, y: elementYCoord}
	adjacentCoordinates["rightCoordinate"] = coordinate{x: elementXCoord + 1, y: elementYCoord}
	adjacentCoordinates["bottomLeftCoordinate"] = coordinate{x: elementXCoord - 1, y: elementYCoord + 1}
	adjacentCoordinates["bottomMiddleCoordinate"] = coordinate{x: elementXCoord, y: elementYCoord + 1}
	adjacentCoordinates["bottomRightCoordinate"] = coordinate{x: elementXCoord + 1, y: elementYCoord + 1}

	for elementLocation, adjacentElement := range adjacentCoordinates {
		if adjacentElement.x >= 0 && adjacentElement.x < len(schematic[elementYCoord]) && adjacentElement.y >= 0 && adjacentElement.y < len(schematic) {
			adjacentElement.value = string(schematic[adjacentElement.y][adjacentElement.x])
			adjacentElement.valid = true
			adjacentCoordinates[elementLocation] = adjacentElement
		}
	}

	return adjacentCoordinates
}

// getPartNumTotal finds the total value of all partNums
// in the schematic.
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

// getGearRatioTotal returns the total of all gear ratios
// in the schematic.
func getGearRatioTotal(schematic []string) int {
	total := 0
	gearRatioPattern := regexp.MustCompile(`\*`)

	for lineIndex, line := range schematic {
		currentLineGearRatios := gearRatioPattern.FindAllStringIndex(line, -1)

		for partNumIndex := range currentLineGearRatios {
			xCoord := currentLineGearRatios[partNumIndex][0]
			currentGear := partNum{id: "*", x: xCoord, y: lineIndex}

			isAdjacent, partNums := currentGear.isAdjacentToGearRatio(schematic)

			if isAdjacent {
				partNum1, err := strconv.Atoi(partNums[0])
				if err != nil {
					log.Fatal(err)
				}

				partNum2, err := strconv.Atoi(partNums[1])
				if err != nil {
					log.Fatal(err)
				}

				product := partNum1 * partNum2

				total += product
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

	partNumTotal := getPartNumTotal(schematic)
	gearRatioTotal := getGearRatioTotal(schematic)

	fmt.Printf("The part num total is: %v\n", partNumTotal)
	fmt.Printf("The gear ratio product total is: %v\n", gearRatioTotal)
}
