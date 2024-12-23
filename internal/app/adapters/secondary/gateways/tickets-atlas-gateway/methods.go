package tickets_atlas_gateway

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/n-kazachuk/go_parser/internal/app/domain/model"
	"strconv"
	"strings"
	"time"
)

func (s *TicketsAtlasGateway) GetTickets(fromCity, toCity, date string) ([]*model.Ticket, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(DOMAIN),
	)

	c.SetRequestTimeout(s.cfg.Timeout)

	proxy := s.cfg.Proxy
	if proxy != "" {
		if err := c.SetProxy(proxy); err != nil {
			return nil, errors.New("Error while set proxy: " + err.Error())
		}
	}

	var tickets []*model.Ticket

	c.OnHTML(".MuiContainer-root .MuiGrid-root.MuiGrid-item.MuiGrid-grid-md-8.MuiGrid-grid-lg-9", func(e *colly.HTMLElement) {
		e.ForEach("div.MuiGrid-root.MuiGrid-container:nth-child(1)", func(i int, el *colly.HTMLElement) {
			if el.DOM.HasClass("MuiGrid-align-items-xs-center") {
				fmt.Println("Return ?")
				return
			}

			timeFromEl := el.DOM.Find("div.MuiGrid-grid-md-3:nth-child(1) > div:first-child > div:first-child > div:first-child")
			timeToEl := el.DOM.Find("div.MuiGrid-grid-md-3:nth-child(2) > div:first-child > div:first-child > div:first-child")
			isFreeEl := el.DOM.Find("button.MuiButton-contained")
			priceEl := el.DOM.Find("h6")

			ticket, err := s.getTicketFromParsedData(timeFromEl, timeToEl, isFreeEl, priceEl, fromCity, toCity, date)
			if err != nil {
				return
			}

			tickets = append(tickets, ticket)
		})
	})

	err := c.Visit(fmt.Sprintf("https://%s/Маршруты/%s/%s?date=%s", DOMAIN, fromCity, toCity, date))
	if err != nil {
		return nil, errors.New("Error while visiting page: " + err.Error())
	}

	return tickets, nil
}

func (s *TicketsAtlasGateway) getTicketFromParsedData(
	timeFromEl, timeToEl, isFreeEl, priceEl *goquery.Selection,
	fromCity, toCity, date string,
) (*model.Ticket, error) {
	ticket := model.NewTicket()

	ticket.FromCity = fromCity
	ticket.ToCity = toCity

	var err error

	dateFormat := "2006-01-02"
	timeFormat := "15:04"

	ticket.Date, err = time.Parse(dateFormat, date)
	if err != nil {
		return nil, err
	}

	ticket.FromTime, err = time.Parse(timeFormat, timeFromEl.Text())
	if err != nil {
		return nil, err
	}

	ticket.ToTime, err = time.Parse(timeFormat, timeToEl.Text())
	if err != nil {
		return nil, err
	}

	ticket.IsFree = !isFreeEl.HasClass("Mui-disabled")

	if ticket.IsFree {
		ticket.Price, err = strconv.ParseFloat(strings.Split(priceEl.Text(), " ")[0], 64)

		if err != nil {
			return nil, err
		}
	}

	return ticket, nil
}
