package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

type Order struct {
	timeFrom string
	timeTo   string
	cost     int
	isFree   bool
}

func NewOrder(timeFrom, timeTo string, isFree bool) *Order {
	return &Order{
		timeFrom: timeFrom,
		timeTo:   timeTo,
		isFree:   isFree,
	}
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("atlasbus.by"),
	)

	var orders []*Order

	c.OnHTML(".MuiContainer-root .MuiGrid-root.MuiGrid-item.MuiGrid-grid-md-8.MuiGrid-grid-lg-9", func(e *colly.HTMLElement) {
		e.ForEach("div.MuiGrid-root.MuiGrid-container:nth-child(1)", func(i int, el *colly.HTMLElement) {
			if el.DOM.HasClass("MuiGrid-align-items-xs-center") {
				return
			}

			timeFromEl := el.DOM.Find("div.MuiGrid-grid-md-3:nth-child(1) > div:first-child > div:first-child > div:first-child")
			timeToEl := el.DOM.Find("div.MuiGrid-grid-md-3:nth-child(2) > div:first-child > div:first-child > div:first-child")
			isFreeEl := el.DOM.Find("button.MuiButton-containedPrimary")

			orders = append(orders, NewOrder(
				timeFromEl.Text(),
				timeToEl.Text(),
				!isFreeEl.HasClass("Mui-disabled"),
			))
		})
	})

	err := c.Visit("https://atlasbus.by/Маршруты/Минск/Житковичи?date=2023-11-17")
	if err != nil {
		fmt.Println("Error while visiting page: " + err.Error())
	}

	for _, value := range orders {
		fmt.Println(value)
	}
}
