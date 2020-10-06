package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var correct int = 0

//Question ...
type Question struct {
	problem string
	answer  string
}

func pressEnter() {
	fmt.Println("\nREADY TO START THE QUIZ!,,, press ENTER :")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func questionRound(data [][]string) []Question {
	val := make([]Question, len(data))
	for i, eachLine := range data {
		val[i] = Question{
			problem: eachLine[0],
			answer:  strings.TrimSpace(eachLine[1]),
		}
	}
	return val
}

func evalEach(i int, ques Question) int {
	fmt.Printf("\nQuestion no.%d: %s = ", i+1, ques.problem)
	var inp string

	fmt.Scanf("%s\n", &inp)
	if ques.answer == inp {
		correct++
		//fmt.Println("OK,AC \n")
	} /*else {
		fmt.Println("oops wrong answer")
	}*/
	return correct
}

func main() {

	timeout := flag.Int("time", 20, "quiz time limit")
	pressEnter()
	csFile, err := os.Open("input.csv")
	if err != nil {
		log.Fatal("file error %v", err)
	}
	defer csFile.Close() //close CSV file

	parsedFile := csv.NewReader(csFile)

	data, err := parsedFile.ReadAll()
	if err != nil {
		log.Fatal("parsing error %v", err)
	}

	//var timeOut int = 10
	timer := time.NewTimer(time.Duration(*timeout) * time.Second)

	var questions []Question
	var acVal int
	completed := make(chan bool, 1)

	go func() {
		questions = questionRound(data)

		for i, q := range questions {
			acVal = evalEach(i, q)
		}
		completed <- true
	}()

	//when timer goes up
	select {
	case <-timer.C:
		fmt.Println("\n\nOops, time is up!!\n")
	case <-completed:
		fmt.Println("\nHurray!!Quiz done & dusted\n")
	}

	fmt.Printf("You scored %d out of %d \n", acVal, len(questions))

}
