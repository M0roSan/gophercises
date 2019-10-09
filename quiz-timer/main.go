package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return ret
}
func main() {
	csvFileName, timeLimit := readFlags()
	lines := readFile(csvFileName)
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	correct := 0
	for i, problem := range problems {
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		fmt.Printf("Problem #%d: %s = ", i+1, problem.question)

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == problem.answer {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}

func readFlags() (string, int) {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	timeLimit := flag.Int("limit", 30, "time limit for the quiz in seconds")
	flag.Parse()
	return *csvFileName, *timeLimit
}

func readFile(fileName string) [][]string {
	file, err := os.Open(fileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", fileName))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file")
	}
	return lines
}
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
