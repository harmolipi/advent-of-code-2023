package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type scratchCard struct {
	id                int
	winningNumbers    []string
	chosenNumbers     []string
	score             int
	numWinningNumbers int
}

func newCard(line string) scratchCard {
	cardContents := strings.Split(line, ": ")
	numberPattern := regexp.MustCompile(`\d+`)
	cardId := numberPattern.FindString(cardContents[0])

	cardIdNum, err := strconv.Atoi(cardId)
	if err != nil {
		log.Fatal(err)
	}

	splitNumbers := strings.Split(cardContents[1], " | ")

	winningNumberSet := splitNumbers[0]
	winningNumbers := numberPattern.FindAllString(winningNumberSet, -1)

	chosenNumberSet := splitNumbers[1]
	chosenNumbers := numberPattern.FindAllString(chosenNumberSet, -1)

	currentScore := 0
	numWinningNumbers := 0

	for _, v := range chosenNumbers {
		if slices.Contains(winningNumbers, v) {
			numWinningNumbers++
			if currentScore == 0 {
				currentScore++
			} else {
				currentScore *= 2
			}
		}
	}

	return scratchCard{
		id:                cardIdNum,
		winningNumbers:    winningNumbers,
		chosenNumbers:     chosenNumbers,
		score:             currentScore,
		numWinningNumbers: numWinningNumbers,
	}
}

func getTotalScore(scratchCards []scratchCard) int {
	score := 0

	for _, v := range scratchCards {
		score += v.score
	}

	return score
}

func getTotalWinningCards(cardSet []scratchCard) int {
	numWinningCards := 0

	for i := range cardSet {
		numWinningCards += 1 + getNumWinningCards(cardSet, i)
	}

	return numWinningCards
}

func getNumWinningCards(cardSet []scratchCard, currentNum int) int {
	card := cardSet[currentNum]
	numWinners := card.numWinningNumbers

	if numWinners == 0 {
		return 0
	} else {
		total := 0
		for i := 1; i <= numWinners; i++ {
			total += getNumWinningCards(cardSet, (currentNum + i))
		}
		return numWinners + total
	}
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
	numWinningCards := getTotalWinningCards(allScratchCards)

	fmt.Printf("The total score of this pile of cards is: %v\n", totalScore)
	fmt.Printf("The number of copied scratch cards is: %v\n", numWinningCards)
}
