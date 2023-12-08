package main

import (
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/services"
	"time"
)

func main() {
	defer scriptRunningTime("main")()

	service := services.NewService()

	fromCity := "Минск"
	toCity := "Житковичи"
	date := "2023-12-01"

	orders, err := service.Parser.GetOrders(fromCity, toCity, date)
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
