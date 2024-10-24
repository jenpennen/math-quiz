package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)
func main() {
	var csvFile *string
	fmt.Printf("Welcome to my math quiz! Firstly, what is your name: ")
	var name string
	fmt.Scan(&name)  // use '&' to store input correctly!!
	fmt.Printf("Hi, %v! How old are you: ", name)
	var age uint
	fmt.Scan(&age)  

	// Check age before starting the quiz
	if age < 0 {
		fmt.Println("Your age must be greater than 0!")	
		return	
	} else if age <= 10 {
		csvFile = flag.String("csv", "problems.csv", "csv file with questions and answers")
		fmt.Println("Ready... Start!")
	} else {
		csvFile = flag.String("csv", "hard-problems.csv", "csv file with questions and answers")
		fmt.Println("Ready...Start!")
	}
	// Set csv file
	// csvFile = flag.String("csv", "hard-problems.csv", "csv file with questions and answers")

	// Set up flags for time limit and randomization
	duration := flag.Int("time", 30, "time limit")
	randomize := flag.Bool("random", false, "randomize questions")
	flag.Parse()

	// Load questions from csv
	reader := bufio.NewReader(os.Stdin)
	questions := loadRecordsFromCsv(*csvFile)
	total := len(questions)
	correct := 0
	mistake := 0

	// Shuffle questions 
	if *randomize {
		questions = shuffle(questions)
	}

	// Quiz logic with a time limit
	done := make(chan bool, 1)
	go func() {
		for i := 0; i < total; i++ {
			if mistake == 3 {
				fmt.Println("You reached the maximum number of wrong answers. Please try again later.")
				done <- true
				return
			}
			fmt.Printf("Question #%d %s = ", i+1, questions[i][0])
			answer, _:= reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			//check if answer is correct
			if strings.Compare(questions[i][1], answer) == 0 {
				correct++
				fmt.Println("Correct!")
			} else {
				mistake++
				fmt.Printf("You got %v wrong. \n", mistake)}
			
		}
		done <- true
	}()

	select {
	case <-done:
		fmt.Printf("You got %v out of %v correct. Good Job!\n", correct, total)
	case <-time.After(time.Duration(*duration) * time.Second):
		fmt.Println("\nYou reached time limit. Better luck next time!")
	}
}


func shuffle(questions [][]string) [][]string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := range questions {
		np := r.Intn(len(questions) - 1)
		questions[i], questions[np] = questions[np], questions[i]
	}
	return questions
}
func loadRecordsFromCsv(csvFile string) [][]string {
	content, err := os.ReadFile(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	
	r := csv.NewReader(bytes.NewReader(content))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return records[1:len(records)]
}