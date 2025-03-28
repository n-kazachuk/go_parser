package tickets_repository_postgres

import (
	"fmt"
	"strings"

	"github.com/n-kazachuk/go_parser/internal/app/domain/ticket"
)

func (s *TicketsRepositoryPostgres) Save(tickets []*ticket.Ticket) error {
	const op = "pgsql.Save"

	if len(tickets) == 0 {
		return nil
	}

	deleteValueStrings := make([]string, 0, len(tickets))
	deleteValueArgs := make([]interface{}, 0, len(tickets)*5)

	insertValueStrings := make([]string, 0, len(tickets))
	insertValueArgs := make([]interface{}, 0, len(tickets)*7)

	for i, ticket := range tickets {
		deleteValueStrings = append(
			deleteValueStrings,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5),
		)

		deleteValueArgs = append(deleteValueArgs, ticket.FromCity, ticket.ToCity, ticket.Date, ticket.FromTime, ticket.ToTime)

		insertValueStrings = append(
			insertValueStrings,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7),
		)

		insertValueArgs = append(insertValueArgs, ticket.FromCity, ticket.ToCity, ticket.Date, ticket.FromTime, ticket.ToTime, ticket.Price, ticket.IsFree)
	}

	deleteQuery := fmt.Sprintf(`
		DELETE FROM ticket
		WHERE (from_city, to_city, date, from_time, to_time) NOT IN (
			SELECT from_city, to_city, date::date, from_time::time, to_time::time 
			FROM (VALUES %s) AS vals(from_city, to_city, date, from_time, to_time)
		)
	`, strings.Join(deleteValueStrings, ", "))

	_, err := s.db.Exec(deleteQuery, deleteValueArgs...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	insertQuery := fmt.Sprintf(`
		INSERT INTO ticket (from_city, to_city, date, from_time, to_time, price, is_free)
		VALUES %s
		ON CONFLICT (from_city, to_city, date, from_time, to_time)
		DO UPDATE SET price = EXCLUDED.price, is_free = EXCLUDED.is_free
	`, strings.Join(insertValueStrings, ", "))

	_, err = s.db.Exec(insertQuery, insertValueArgs...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
