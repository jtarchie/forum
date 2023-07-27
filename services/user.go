package services

import (
	"fmt"
	"math/rand"

	"github.com/dustinkirkland/golang-petname"
	"github.com/jtarchie/forum/db"
)

type User struct {
	Email    string
	Provider string
}

func UpsertUser(
	client db.Client,
	user User,
) error {
	email := user.Email
	provider := user.Provider

	rows, err := client.Query("SELECT COUNT(*) FROM users WHERE email = ? AND provider = ?;", email, provider)
	if err != nil {
		return fmt.Errorf("could not find user (%q): %w", email, err)
	}
	defer rows.Close()

	var userCount int

	if !rows.Next() {
		return fmt.Errorf("could not query user (%q): %w", email, err)
	}

	err = rows.Scan(&userCount)
	if err != nil {
		return fmt.Errorf("could not determine user (%q): %w", email, err)
	}

	err = rows.Close()
	if err != nil {
		return fmt.Errorf("could close user cursor: %w", err)
	}

	if userCount == 0 {
		username := fmt.Sprintf("%s-%d", petname.Generate(2, "-"), rand.Intn(9999))

		err := client.Execute("INSERT INTO users (email, provider, username) VALUES (?, ?, ?);", email, provider, username)
		if err != nil {
			return fmt.Errorf("could not insert user (%q): %w", email, err)
		}
	}

	return nil
}
