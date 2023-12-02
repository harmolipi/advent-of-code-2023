package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	if !(len(os.Args) >= 2) {
		log.Fatal("Please include the name of the input file. For example:\n\t`go run main.go input.txt`")
	}

	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("There was an error reading the file! Please make sure it exists.")
	}
	defer file.Close()

	calibrationTotal := 0
	counter := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		calibrationTotal += getCalibrationValue(scanner.Text(), counter)
		counter++
	}

	fmt.Println(calibrationTotal)
}

func getCalibrationValue(line string, count int) int {
	numberMap := map[string]string{
		"one":   "1",
		"two":   "2",
		"three": "3",
		"four":  "4",
		"five":  "5",
		"six":   "6",
		"seven": "7",
		"eight": "8",
		"nine":  "9",
	}

	pattern := regexp.MustCompile(`\d|one|two|three|four|five|six|seven|eight|nine`)

	var matches []string

	for i := range line {
		foundMatch := pattern.FindString(line[i:])
		if foundMatch != "" {
			matches = append(matches, foundMatch)
		}
	}

	firstDigit, ok := numberMap[matches[0]]
	if !ok {
		firstDigit = matches[0]
	}

	secondDigit, ok := numberMap[matches[len(matches)-1]]
	if !ok {
		secondDigit = matches[len(matches)-1]
	}

	calibrationValue, err := strconv.Atoi(firstDigit + secondDigit)
	if err != nil {
		log.Fatal(err)
	}

	return calibrationValue
}
