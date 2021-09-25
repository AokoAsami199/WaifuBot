package discord

import (
	"fmt"

	"github.com/Karitham/WaifuBot/service/anilist"
	"github.com/Karitham/WaifuBot/service/store"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

//go:generate go-mockgen -f github.com/Karitham/WaifuBot/discord -i randomer -o randomerMock_test.go
type randomer interface {
	Random(notIn []int) (anilist.Character, error)
}

//go:generate go-mockgen -f github.com/Karitham/WaifuBot/discord -i storer -o storerMock_test.go
type storer interface {
	Put(userID int, c store.Character) error
	Get(userID int) ([]store.Character, error)
}

func (b *bot) Roller(r randomer, s storer) slashCmd {
	return func(e *gateway.InteractionCreateEvent, d discord.CommandInteractionData) api.InteractionResponseData {
		chars, err := s.Get(int(e.Member.User.ID))
		if err != nil {
			return api.InteractionResponseData{
				Content: option.NewNullableString("An error occurred dialing the database, please try again later"),
				Flags:   api.EphemeralResponse,
			}
		}

		c, err := r.Random(IDs(chars))
		if err != nil {
			return api.InteractionResponseData{
				Content: option.NewNullableString("An error getting a random character occurred, please try again later"),
				Flags:   api.EphemeralResponse,
			}
		}

		return api.InteractionResponseData{
			Embeds: &[]discord.Embed{{
				Title:       c.Name.Full,
				Description: fmt.Sprintf("You rolled %s.\nCongratulations!", c.Name.Full),
				URL:         c.Siteurl,
				Color:       anilist.Color,
				Footer: &discord.EmbedFooter{
					Icon: anilist.IconURL,
					Text: "View them on anilist",
				},
				Thumbnail: &discord.EmbedThumbnail{
					URL: c.Image.Large,
				},
			}},
		}
	}
}

func IDs(c []store.Character) []int {
	ids := make([]int, len(c))
	for i, v := range c {
		ids[i] = int(v.ID)
	}

	return ids
}
