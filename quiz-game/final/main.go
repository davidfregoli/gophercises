package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var timesUpCh = make(chan struct{})
var score int = 0

func main() {
	filePath := flag.String("problems", "problems.json", "The path to the problems file.")
	timeLimit := flag.Int("time", 10, "The time limit for the quiz in seconds.")
	flag.Parse()

	jsonFile, err := os.Open(*filePath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var config Config
	json.Unmarshal(byteValue, &config)
	problems := config.Parse()
	problemsCount := len(problems)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

GameLoop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case answer := <-answerCh:
			if answer == problem.a {
				score++
			}
		case <-timer.C:
			fmt.Println("\nTIME'S UP!")
			break GameLoop
		}
		close(answerCh)
	}

	fmt.Printf("You guessed right %d/%d answers.\n", score, problemsCount)
	close(timesUpCh)
}

type Config struct {
	Problems [][]string `json:"problems"`
}

type Problem struct {
	q string
	a string
}

func (c Config) Parse() []Problem {
	problems := make([]Problem, len(c.Problems))
	for i, problem := range c.Problems {
		problems[i] = Problem{
			q: problem[0],
			a: problem[1],
		}
	}
	return problems
}
