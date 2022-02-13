package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

const problemsFile string = "/Users/david/code/go/src/github.com/davidfregoli/gophercises/quiz-game/mvp/problems.json"

var wg = sync.WaitGroup{}
var timesUpCh = make(chan struct{})
var answerCh = make(chan string)
var score int = 0

func main() {
	jsonFile, err := os.Open(problemsFile)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	var problems Problems
	json.Unmarshal(byteValue, &problems)
	problemsCount := len(problems.Problems)

	wg.Add(1)
	go func(seconds int32) {
		duration := time.Duration(seconds)
		timer := time.NewTimer(duration * time.Second)
		<-timer.C
		timesUpCh <- struct{}{}
		wg.Done()
	}(30)

GameLoop:
	for i, problem := range problems.Problems {
		wg.Add(1)
		go AskQuestion(problem[0])
		select {
		case answer := <-answerCh:
			wg.Done()
			if answer == problem[1] {
				score++
			}
		case <-timesUpCh:
			wg.Done()
			fmt.Println("\nTIME'S UP!")
			break GameLoop
		}
		if i == problemsCount-1 {
			wg.Done()
		}
	}

	fmt.Println("You guessed right", score, "/", problemsCount, "answers.")
	close(timesUpCh)
	close(answerCh)
}

func AskQuestion(label string) {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+"? ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	s = strings.TrimSpace(s)
	answerCh <- s
}

type Problems struct {
	Problems []Problem `json:"problems"`
}

type Problem [2]string
