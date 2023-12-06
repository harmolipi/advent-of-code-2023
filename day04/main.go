package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
)

type scratchCard struct {
	winningNumbers []string
	chosenNumbers  []string
	score          int
}

func newCard(line string) scratchCard {
	cardContents := strings.Split(line, ": ")
	splitNumbers := strings.Split(cardContents[1], " | ")
	numberPattern := regexp.MustCompile(`\d+`)

	winningNumberSet := splitNumbers[0]
	winningNumbers := numberPattern.FindAllString(winningNumberSet, -1)

	chosenNumberSet := splitNumbers[1]
	chosenNumbers := numberPattern.FindAllString(chosenNumberSet, -1)

	currentScore := 0

	for _, v := range chosenNumbers {
		if slices.Contains(winningNumbers, v) {
			if currentScore == 0 {
				currentScore++
			} else {
				currentScore *= 2
			}
		}
	}

	return scratchCard{
		winningNumbers: winningNumbers,
		chosenNumbers:  chosenNumbers,
		score:          currentScore,
	}
}

func getTotalScore(scratchCards []scratchCard) int {
	score := 0

	for _, v := range scratchCards {
		score += v.score
	}

	return score
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please include an input file as an argument.")
	}

	filepath := os.Args[1]

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var allScratchCards []scratchCard

	for scanner.Scan() {
		thisCard := newCard(scanner.Text())
		allScratchCards = append(allScratchCards, thisCard)
	}

	totalScore := getTotalScore(allScratchCards)

	fmt.Printf("The total score of this pile of cards is: %v\n", totalScore)
}
