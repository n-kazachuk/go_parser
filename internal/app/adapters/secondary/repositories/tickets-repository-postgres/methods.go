package tickets_repository_postgres

import (
	"fmt"
	"github.com/n-kazachuk/go_parser/internal/app/domain/ticket"
	"strings"
)

func (s *TicketsRepositoryPostgres) Save(tickets []*ticket.Ticket) error {
	const op = "pgsql.Save"

	query := `
		INSERT INTO ticket (from_city, to_city, date, from_time, to_time, price, is_free)
		VALUES %s
		ON CONFLICT (from_city, to_city, date, from_time, to_time)
		DO UPDATE SET price = EXCLUDED.price
	`

	valueStrings := make([]string, 0, len(tickets))
	valueArgs := make([]interface{}, 0, len(tickets)*6)

	for i, ticket := range tickets {
		valueStrings = append(
			valueStrings,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7),
		)

		valueArgs = append(
			valueArgs,
			ticket.FromCity, ticket.ToCity, ticket.Date, ticket.FromTime, ticket.ToTime, ticket.Price, ticket.IsFree,
		)
	}

	finalQuery := fmt.Sprintf(query, strings.Join(valueStrings, ", "))

	_, err := s.db.Exec(finalQuery, valueArgs...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
