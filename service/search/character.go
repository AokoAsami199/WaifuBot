package search

import (
	"context"
	"errors"

	"github.com/machinebox/graphql"
)

type name struct {
	Full string `json:"full"`
}

type image struct {
	Large string `json:"large"`
}

type Character struct {
	Name        name   `json:"name"`
	Description string `json:"description"`
	Siteurl     string `json:"siteUrl"`
	Image       image  `json:"image"`
	ID          int    `json:"id"`
}

// User queries the anilist user
func (a *Anilist) Character(name string) ([]Character, error) {
	var q struct {
		Page struct {
			Characters []Character `json:"characters"`
		} `json:"page"`
	}
	req := graphql.NewRequest(`
query ($name: String) {
  Page {
    characters(search: $name) {
      id
      name {
        full
      }
      description
      siteUrl
      image {
        large
      }
    }
  }
}

	`)
	req.Var("name", name)

	err := a.c.Run(context.Background(), req, &q)
	if err != nil {
		return nil, err
	}
	if len(q.Page.Characters) < 1 {
		return nil, errors.New("no characters found")
	}

	return q.Page.Characters, nil
}
