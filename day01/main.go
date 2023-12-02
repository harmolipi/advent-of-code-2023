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

	scanner := bufio.NewScanner(file)
	pattern := regexp.MustCompile(`\d`)
	calibrationTotal := 0

	for scanner.Scan() {
		lineDigits := pattern.FindAllString(scanner.Text(), -1)
		calibrationString := lineDigits[0] + lineDigits[len(lineDigits)-1]

		calibrationValue, err := strconv.Atoi(calibrationString)
		if err != nil {
			log.Fatal("There was an error parsing the calibration value!")
		}

		calibrationTotal += calibrationValue
	}

	fmt.Println(calibrationTotal)
}
