package discord

import (
	"fmt"

	"github.com/Karitham/WaifuBot/service/anilist"
	"github.com/Karitham/WaifuBot/service/store"
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

type Randomer interface {
	Random(notIn []int) (anilist.Character, error)
}

type Storager interface {
	Put(userID int, c store.Character) error
	Get(userID int) ([]store.Character, error)
}

func (b *bot) Roller(r Randomer, s Storager) Command {
	cmd := api.CreateCommandData{
		Name:        "roll",
		Description: "Roll a character, obtain it for yourself",
	}

	fn := func(e *gateway.InteractionCreateEvent) {
		chars, err := s.Get(int(e.Member.User.ID))
		if err != nil {
			b.RespondWithError(e, "An error occurred dialing the database, please try again later", err).Msg("database error")
			return
		}

		c, err := r.Random(IDs(chars))
		if err != nil {
			b.RespondWithError(e, "An error getting a random character occurred, please try again later", err).Msg("anilist error")
		}

		b.RespondWithEmbed(e, discord.Embed{
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
		})
	}

	return Command{
		callback:   fn,
		apiCommand: cmd,
	}
}

func IDs(c []store.Character) []int {
	ids := make([]int, len(c))
	for i, v := range c {
		ids[i] = int(v.ID)
	}

	return ids
}
