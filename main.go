package main

import (
	"fmt"
	"time"
)

// Функция для запуска таймера
func startTimer(duration int, message string) {
	fmt.Printf("Начинается %s (%d минут)...\n", message, duration)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	totalSeconds := duration * 60
	for i := totalSeconds; i > 0; i-- {
		minutes := i / 60
		seconds := i % 60
		fmt.Printf("\r%02d:%02d осталось", minutes, seconds)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("\nВремя вышло!")
}

func main() {
	// Длительность работы и отдыха в минутах
	workDuration := 25 // 25 минут работы
	breakDuration := 5 // 5 минут отдыха

	fmt.Println("Добро пожаловать в Pomodoro Timer!")

	for {
		// Запуск таймера для работы
		startTimer(workDuration, "таймер работы")

		// Запрос на продолжение
		fmt.Print("Хотите сделать перерыв? (y/n): ")
		var input string
		fmt.Scanln(&input)
		if input != "y" {
			fmt.Println("Работа завершена. До свидания!")
			break
		}

		// Запуск таймера для отдыха
		startTimer(breakDuration, "таймер отдыха")

		// Запрос на продолжение
		fmt.Print("Хотите начать новый цикл? (y/n): ")
		fmt.Scanln(&input)
		if input != "y" {
			fmt.Println("Работа завершена. До свидания!")
			break
		}
	}
}
