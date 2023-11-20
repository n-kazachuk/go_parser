package main

import (
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/services"
	"time"
)

func main() {
	defer scriptRunningTime("main")()

	fromCity := "Минск"
	toCity := "Петриков"
	date := "2023-11-24"

	service := services.NewService()
	orders, err := service.ParseOrders(fromCity, toCity, date)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, value := range orders {
		fmt.Println(value)
	}
}

func scriptRunningTime(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
