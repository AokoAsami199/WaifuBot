package search

import "github.com/machinebox/graphql"

type Anilist struct {
	c   *graphql.Client
	URL string
}

func New() *Anilist {
	const graphURL = "https://graphql.anilist.co"

	return &Anilist{URL: graphURL, c: graphql.NewClient(graphURL)}
}
