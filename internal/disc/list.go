package disc

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Karitham/WaifuBot/internal/db"
	"github.com/diamondburned/arikawa/v2/bot/extras/arguments"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/dgwidgets"
	"github.com/rs/zerolog/log"
)

// List shows the user's list
func (b *Bot) List(m *gateway.MessageCreateEvent, _ ...*arguments.UserMention) error {
	user := parseUser(m)

	uData, err := b.conn.GetUserList(context.Background(), int64(m.Author.ID))
	if err == sql.ErrNoRows {
		err := b.conn.CreateUser(b.Ctx.Context(), int64(m.Author.ID))
		if err != nil {
			log.Err(err).
				Str("Type", "LIST").
				Int("user", int(m.Author.ID)).
				Msg("Could not create user")

			return err
		}
	} else if err != nil {
		log.Err(err).
			Str("Type", "LIST").
			Msg("Could not get list")

		return err
	}
	// Create widget
	p := dgwidgets.NewPaginator(b.Ctx.State, m.ChannelID)

	// What to do when timeout
	p.SetTimeout(b.conf.ListMaxUpdateTime.Duration)
	p.ColourWhenDone = 0xFFFF00

	// Make pages
	for j := 0; j <= len(uData)/b.conf.ListLen; j++ {
		p.Add(discord.Embed{
			Title: fmt.Sprintf("%s's list", user.Username),
			Description: func(l []db.Character) string {
				var s strings.Builder
				if len(l) >= 0 {
					for i := b.conf.ListLen * j; i < b.conf.ListLen+b.conf.ListLen*j && i < len(l); i++ {
						s.WriteString(fmt.Sprintf("`%d`\f - %s\n", l[i].ID, l[i].Name.String))
					}
					return s.String()
				}
				return ""
			}(uData),

			Color: 3447003,
		})
	}

	log.Trace().
		Str("Type", "LIST").
		Int("User", int(user.ID)).
		Msg("Sent list embed")

	return p.Spawn()
}

// Verify verify if someone has a waifu
func (b *Bot) Verify(m *gateway.MessageCreateEvent, ID int64, _ ...*arguments.UserMention) (string, error) {
	user := parseUser(m)
	log.Trace().
		Str("Type", "VERIFY").
		Int("User", int(user.ID)).
		Int64("Char", ID).
		Msg("verifying ownership")

	if _, err := b.conn.GetChar(context.Background(), db.GetCharParams{
		ID:     ID,
		UserID: int64(user.ID),
	}); err == sql.ErrNoRows {
		return fmt.Sprintf("%s doesn't own the character %d", user.Username, ID), nil
	} else if err != nil {
		log.Err(err).
			Str("Type", "VERIFY").
			Int("ID", int(user.ID)).
			Msg("Could not verify")
	}
	return fmt.Sprintf("%s owns the character %d", user.Username, ID), nil
}
