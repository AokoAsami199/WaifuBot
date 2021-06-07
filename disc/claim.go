package disc

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/Karitham/WaifuBot/anilist"
	"github.com/Karitham/WaifuBot/db"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

// Dropper is used to handle the dropping mechanism
type Dropper struct {
	Waifu   map[discord.ChannelID]anilist.CharStruct
	ChanInc map[discord.ChannelID]uint64
	Mutex   *sync.Mutex
}

func (bot *Bot) drop(m *gateway.MessageCreateEvent) {
	bot.dropper.Mutex.Lock()
	defer bot.dropper.Mutex.Unlock()

	var err error
	bot.dropper.Waifu[m.ChannelID], err = anilist.CharSearchByPopularity(bot.seed.Uint64()%bot.conf.MaxCharacterRoll, []int64{})
	if err != nil {
		log.Err(err).
			Str("Type", "DROP").
			Msg("Error getting char from anilist")

		return
	}

	// Sanitize the name so it's claimable through discord (some characters have double spaces in their name)
	bot.dropper.Waifu[m.ChannelID].Page.Characters[0].Name.Full =
		strings.Join(strings.Fields(bot.dropper.Waifu[m.ChannelID].Page.Characters[0].Name.Full), " ")

	_, err = bot.Ctx.SendMessage(m.ChannelID, "", &discord.Embed{
		Title:       "CHARACTER DROP !",
		Description: "Can you guess who it is ?\nUse w.claim to get this character for yourself",
		Thumbnail: &discord.EmbedThumbnail{
			URL: bot.dropper.Waifu[m.ChannelID].Page.Characters[0].Image.Large,
		},
		Footer: &discord.EmbedFooter{
			Text: "This character's initials are " +
				func(name string) (initials string) {
					for _, v := range strings.Fields(name) {
						initials = initials + strings.ToUpper(string(v[0])) + "."
					}
					return
				}(bot.dropper.Waifu[m.ChannelID].Page.Characters[0].Name.Full),
		},
	})
	if err != nil {
		log.Err(err).Str("Type", "DROP").Msg("Error sending drop message")
	}
}

// Claim a waifu and adds it to the user's database
func (bot *Bot) Claim(m *gateway.MessageCreateEvent, name ...Name) (*discord.Embed, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("if you want to claim a character, use `claim <name>`")
	}

	// Lock because we are reading from the map
	bot.dropper.Mutex.Lock()
	defer bot.dropper.Mutex.Unlock()
	char, ok := bot.dropper.Waifu[m.ChannelID]

	if !ok {
		return nil, fmt.Errorf("there is no character to claim")
	}

	if !strings.EqualFold(
		strings.Join(name, " "),
		char.Page.Characters[0].Name.Full,
	) {
		return nil, fmt.Errorf("wrong name entered")
	}

	// Add to db
	err := bot.DB.InsertChar(context.Background(), db.InsertCharParams{
		ID:     char.Page.Characters[0].ID,
		UserID: int64(m.Author.ID),
		Image:  char.Page.Characters[0].Image.Large,
		Name:   char.Page.Characters[0].Name.Full,
	})
	if err != nil {
		log.Err(err).
			Str("Type", "CLAIM").
			Int64("ID", char.Page.Characters[0].ID).
			Int("UserID", int(m.Author.ID)).
			Msg("Error inserting the char")
		return nil, err
	}

	delete(bot.dropper.Waifu, m.ChannelID)

	return &discord.Embed{
		Title: "Claim successful",
		URL:   char.Page.Characters[0].SiteURL,
		Description: fmt.Sprintf(
			"Well done %s you claimed %d\nIt appears in :\n- %s",
			m.Author.Username, char.Page.Characters[0].ID, char.Page.Characters[0].Media.Nodes[0].Title.Romaji,
		),
		Thumbnail: &discord.EmbedThumbnail{
			URL: char.Page.Characters[0].Image.Large,
		},
	}, nil
}
