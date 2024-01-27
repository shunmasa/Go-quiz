package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type OpenTriviaResponse struct {
	Results []Question `json:"results"`
}

type Question struct {
	Question        string   `json:"question"`
	CorrectAnswer   string   `json:"correct_answer"`
	IncorrectAnswer []string `json:"incorrect_answers"`
}

type QuizOptions struct {
	NumQuestions int
	Difficulty   string
	Category     int
}

func NewQuizOptions() QuizOptions {
	var numQuestions int
	fmt.Print("Enter the number of questions: ")
	fmt.Scan(&numQuestions)

	var difficulty string
	for {
		fmt.Print("Enter the difficulty (easy, medium, hard): ")
		fmt.Scan(&difficulty)
		if difficulty == "easy" || difficulty == "medium" || difficulty == "hard" {
			break
		} else {
			fmt.Println("Error: Invalid difficulty. Choose from easy, medium, or hard.")
		}
	}

	var category int
	for {
		fmt.Println("Choose a category:")
		fmt.Println("9. General Knowledge")
		fmt.Println("14. TV")
		fmt.Println("10. Books")
		fmt.Println("12. Music")
		fmt.Println("11. Film")
		fmt.Print("Enter the category number: ")
		fmt.Scan(&category)

		validCategories := []int{9, 14, 10, 12, 11}
		validCategory := false
		for _, valid := range validCategories {
			if category == valid {
				validCategory = true
				break
			}
		}

		if validCategory {
			break
		} else {
			fmt.Println("Error: Invalid category. Choose a valid category number.")
		}
	}

	return QuizOptions{
		NumQuestions: numQuestions,
		Difficulty:   difficulty,
		Category:     category,
	}
}

func main() {
	options := NewQuizOptions()

	quizAPIURL := fmt.Sprintf("https://opentdb.com/api.php?amount=%d&category=%d&difficulty=%s&type=multiple", options.NumQuestions, options.Category, options.Difficulty)
	quizQuestions := fetchQuizQuestions(quizAPIURL)

	runQuiz(quizQuestions)
}

func fetchQuizQuestions(apiURL string) []Question {
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Failed to fetch quiz questions")
		os.Exit(1)
	}
	defer response.Body.Close()

	var openTriviaResponse OpenTriviaResponse
	err = json.NewDecoder(response.Body).Decode(&openTriviaResponse)
	if err != nil {
		fmt.Printf("Failed to parse JSON response: %s\n", err)
		return nil
	}

	return openTriviaResponse.Results
}

func runQuiz(questions []Question) {
	var score, correctAnswers, incorrectAnswers int

	for i, question := range questions {
		fmt.Printf("Question %d: %s\n", i+1, question.Question)

		choices := append(question.IncorrectAnswer, question.CorrectAnswer)
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(choices), func(i, j int) { choices[i], choices[j] = choices[j], choices[i] })

		for j, choice := range choices {
			fmt.Printf("%d. %s\n", j+1, choice)
		}

		fmt.Print("Your Answer: ")
		var userAnswer string
		fmt.Scan(&userAnswer)

		userAnswer = strings.TrimSpace(userAnswer)

		index := 0
		_, err := fmt.Sscanf(userAnswer, "%d", &index)
		if err == nil && index >= 1 && index <= len(choices) {
			selectedChoice := choices[index-1]
			if selectedChoice == question.CorrectAnswer {
				fmt.Println("Correct!\n")
				score++
				correctAnswers++
			} else {
				fmt.Printf("Incorrect. The correct answer is: %s\n\n", question.CorrectAnswer)
				incorrectAnswers++
			}
		} else {
			fmt.Println("Invalid input. Skipping question.\n")
			incorrectAnswers++
		}
	}

	totalQuestions := float64(len(questions))
	correctPercentage := (float64(correctAnswers) / totalQuestions) * 100.0
	incorrectPercentage := (float64(incorrectAnswers) / totalQuestions) * 100.0

	fmt.Printf("Quiz completed. Your score: %d/%d (%.2f%% correct, %.2f%% incorrect)\n", score, int(totalQuestions), correctPercentage, incorrectPercentage)
}
