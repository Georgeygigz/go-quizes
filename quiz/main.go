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
	q string
	a string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 8, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s Error: %s\n", *csvFilename, err))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file")
	}
	problems := parseLines(lines)
	correct := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

problemloop:
	for _, p := range problems {
		fmt.Printf("What is %s\n", p.q)

		answerCh := make(chan string)
		go func() {
			var ans string
			fmt.Scanf("%s", &ans)
			answerCh <- ans
		}()

		select {
		case <-timer.C:
			fmt.Println("\n>>>TIME EXCEED<<<")
			break problemloop
		case ans := <-answerCh:
			if ans == p.a {
				correct++
			}
		}

	}
	fmt.Printf(">>>>>CHALLENGE COMPLETED<<<<<\n Total Questions: %3d\n Correct Answers: %3d\n Wrong Answers: %3d\n", len(lines), correct, len(lines)-correct)

}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
