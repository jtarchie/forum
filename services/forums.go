package services

import (
	_ "embed"
	"fmt"

	"github.com/jtarchie/forum/db"
)

type Forum struct {
	ID          int64
	Name        string
	Description string
	ParentID    int64
}

type Forums []Forum

func ListForums(client db.Client) (Forums, error) {
	rows, err := client.Query("SELECT * FROM ordered_forums;")
	if err != nil {
		return nil, fmt.Errorf("could not read forums: %w", err)
	}

	var forums Forums

	for rows.Next() {
		var forum Forum

		err := rows.Scan(&forum.ID, &forum.Name, &forum.Description, &forum.ParentID)
		if err != nil {
			return nil, fmt.Errorf("could not read properties from forums: %w", err)
		}

		forums = append(forums, forum)
	}

	return forums, nil
}
