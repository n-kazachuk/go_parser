package tickets_atlas_gateway

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/n-kazachuk/go_parser/internal/app/domain/ticket"
	"github.com/n-kazachuk/go_parser/internal/libs/helpers"

	ticketsRequest "github.com/n-kazachuk/go_parser/internal/app/domain/tickets-request"
)

func (s *TicketsAtlasGateway) GetTickets(ticketRequest *ticketsRequest.TicketRequest) ([]*ticket.Ticket, error) {
	op := helpers.GetFunctionName()

	c := colly.NewCollector(
		colly.AllowedDomains(DOMAIN),
	)

	c.SetRequestTimeout(s.cfg.Timeout)

	proxy := s.cfg.Proxy
	if proxy != "" {
		if err := c.SetProxy(proxy); err != nil {
			return nil, fmt.Errorf("%s: error while set proxy: %s", op, err.Error())
		}
	}

	var tickets []*ticket.Ticket

	c.OnHTML(".MuiContainer-root .MuiGrid-root.MuiGrid-item.MuiGrid-grid-md-8.MuiGrid-grid-lg-9", func(e *colly.HTMLElement) {
		e.ForEach("div.MuiGrid-root.MuiGrid-container:nth-child(1)", func(i int, el *colly.HTMLElement) {
			if el.DOM.HasClass("MuiGrid-align-items-xs-center") {
				return
			}

			timeFromEl := el.DOM.Find("div.MuiGrid-grid-md-3:nth-child(1) > div:first-child > div:first-child > div:first-child")
			timeToEl := el.DOM.Find("div.MuiGrid-grid-md-3:nth-child(2) > div:first-child > div:first-child > div:first-child")
			isFreeEl := el.DOM.Find("button.MuiButton-contained")
			priceEl := el.DOM.Find("h6")

			ticket, err := s.getTicketFromParsedData(timeFromEl, timeToEl, isFreeEl, priceEl)
			if err != nil {
				return
			}

			ticket.FromCity = ticketRequest.FromCity
			ticket.ToCity = ticketRequest.ToCity
			ticket.Date = ticketRequest.Date

			tickets = append(tickets, ticket)
		})
	})

	err := c.Visit(fmt.Sprintf("https://%s/Маршруты/%s/%s?date=%s", DOMAIN, ticketRequest.FromCity, ticketRequest.ToCity, ticketRequest.Date.Format("2006-01-02")))
	if err != nil {
		return nil, fmt.Errorf("%s: error while visiting page: %s", op, err.Error())
	}

	return tickets, nil
}

func (s *TicketsAtlasGateway) getTicketFromParsedData(
	timeFromEl,
	timeToEl,
	isFreeEl,
	priceEl *goquery.Selection,
) (*ticket.Ticket, error) {
	ticket := ticket.New()

	var err error

	timeFormat := "15:04"

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
