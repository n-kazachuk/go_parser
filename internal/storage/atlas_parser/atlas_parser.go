package atlas_parser

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/n-kazachuk/go_parser/internal/domain/models"
)

const (
	DOMAIN = "atlasbus.by"
)

type AtlasStorage struct {
}

func NewAtlasStorage() *AtlasStorage {
	return &AtlasStorage{}
}

func (s *AtlasStorage) GetOrders(fromCity, toCity, date string) ([]*models.Order, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(DOMAIN),
	)

	var orders []*models.Order

	c.OnHTML(".MuiContainer-root .MuiGrid-root.MuiGrid-item.MuiGrid-grid-md-8.MuiGrid-grid-lg-9", func(e *colly.HTMLElement) {
		e.ForEach("div.MuiGrid-root.MuiGrid-container:nth-child(1)", func(i int, el *colly.HTMLElement) {
			if el.DOM.HasClass("MuiGrid-align-items-xs-center") {
				return
			}

			timeFromEl := el.DOM.Find("div.MuiGrid-grid-md-3:nth-child(1) > div:first-child > div:first-child > div:first-child")
			timeToEl := el.DOM.Find("div.MuiGrid-grid-md-3:nth-child(2) > div:first-child > div:first-child > div:first-child")
			isFreeEl := el.DOM.Find("button.MuiButton-contained")

			orders = append(orders, models.NewOrder(
				timeFromEl.Text(),
				timeToEl.Text(),
				!isFreeEl.HasClass("Mui-disabled"),
			))
		})
	})

	err := c.Visit(fmt.Sprintf("https://%s/Маршруты/%s/%s?date=%s", DOMAIN, fromCity, toCity, date))
	if err != nil {
		return nil, errors.New("Error while visiting page: " + err.Error())
	}

	return orders, nil
}
