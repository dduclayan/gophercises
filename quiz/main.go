package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	questionsPath = flag.String("path", "questions.csv", "path to questions csv file")
	timer         = flag.Int("timer", 60, "how many seconds to answer all questions in the quiz")
	randomize     = flag.Bool("randomize", false, "randomize question order")
)

// Quiz represents a timed quiz game
type Quiz struct {
	questions      []questionAnswer
	correctAnswers int
}

type questionAnswer struct {
	question string
	answer   string
}

func startTimer(timeLimit int, ch chan bool) {
	limit := time.Second * time.Duration(timeLimit)
	start := time.Now()
	for {
		if time.Since(start) >= limit {
			ch <- true
		}
	}
}

func extractQuestions(path string, shuffle bool) ([]questionAnswer, error) {
	infoLogger := log.New(os.Stderr, "INFO: ", log.Lshortfile)
	var questions []questionAnswer
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	csv := csv.NewReader(f)
	for {
		row, err := csv.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return questions, nil
			}
			return nil, fmt.Errorf("could not read row: %w", err)
		}
		question := row[0]
		answer := row[1]

		// skip header
		if question == "question" {
			continue
		}

		questions = append(questions, questionAnswer{question: question, answer: answer})

		if shuffle {
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
		}

		infoLogger.Printf("added question %q and answer %q to question bank", question, answer)
	}
}

func (q *Quiz) play(timeLimit int) {
	fmt.Printf("Press enter when you are ready to play. You will have %d seconds for the whole quiz\n", *timer)
	input := bufio.NewReader(os.Stdin)
	input.ReadString('\n')

	timesUp := make(chan bool)
	go startTimer(timeLimit, timesUp)

	go func() {
		for _, val := range q.questions {
			fmt.Printf("\n%q\n", val.question)
			input := bufio.NewReader(os.Stdin)
			userAnswer, _ := input.ReadString('\n')
			userAnswer = strings.Trim(userAnswer, "\n")
			if strings.EqualFold(strings.ToLower(userAnswer), strings.ToLower(val.answer)) {
				q.correctAnswers++
			}
		}
		fmt.Printf("\nScore: %d/%d\n", q.correctAnswers, len(q.questions))
	}()

	if <-timesUp {
		fmt.Println("times up!")
		fmt.Printf("\nScore: %d/%d\n", q.correctAnswers, len(q.questions))
		return
	}
}

func new(path string) Quiz {
	q := Quiz{}
	questionBank, err := extractQuestions(path, *randomize)
	if err != nil {
		fmt.Printf("could not extract questions: %v\n", err)
		os.Exit(1)
	}
	q.questions = questionBank
	return q
}

func main() {
	flag.Parse()
	quiz := new(*questionsPath)
	quiz.play(*timer)
}
