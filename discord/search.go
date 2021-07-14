package discord

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Karitham/WaifuBot/service/search"
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/rs/zerolog/log"
)

type SearchProvider interface {
	Anime(string) ([]search.Media, error)
	Manga(string) ([]search.Media, error)
	Character(string) ([]search.Character, error)
	User(string) ([]search.User, error)
}

type Interaction struct {
	Token   string
	Options []gateway.InteractionOption
	ID      discord.InteractionID
}

func (b *bot) Search(s SearchProvider) Command {
	subcommands := map[string]func(e Interaction){}
	a := b.regAnimeInteraction(subcommands, s)
	m := b.regMangaInteraction(subcommands, s)
	c := b.regCharInteraction(subcommands, s)
	u := b.regUserInteraction(subcommands, s)

	callback := func(e *gateway.InteractionCreateEvent) {
		// Call the right function based on the subcommand
		for _, opt := range e.Data.Options {
			subcommands[opt.Name](Interaction{
				Token:   e.Token,
				Options: opt.Options,
				ID:      e.ID,
			})
		}
	}

	return Command{
		callback: callback,
		apiCommand: api.CreateCommandData{
			Name:        "search",
			Description: "search for anything on anilist poggers",
			Options:     []discord.CommandOption{a, m, c, u},
		},
	}
}

func (b *bot) regAnimeInteraction(
	subcom map[string]func(e Interaction),
	s SearchProvider,
) discord.CommandOption {
	opt := discord.CommandOption{
		Name:        "anime",
		Description: "search for an anime",
		Type:        discord.SubcommandOption,
		Required:    false,
		Options: []discord.CommandOption{{
			Type:        discord.StringOption,
			Name:        "title",
			Description: "title of the anime",
			Required:    true,
		}},
	}

	fn := func(i Interaction) {
		log.Trace().Interface("options", i.Options).Msg("Searching for anime")

		anime, err := s.Anime(i.Options[0].Value)
		if err != nil {
			_ = b.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: "Error searching for this anime, either it doesn't exist or something went wrong",
				},
			})
			log.Err(err).Str("title", i.Options[0].Value).Msg("Error searching for the anime")
		}

		if len(anime) < 1 {
			_ = b.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: "Error searching for this anime, either it doesn't exist or something went wrong",
				},
			})
			return
		}

		err = b.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
			Type: api.MessageInteractionWithSource,
			Data: &api.InteractionResponseData{
				Embeds: []discord.Embed{buildMediaEmbed(anime[0])},
			},
		})
		if err != nil {
			log.Err(err).Msg("failed to send interaction callback")
		}
	}

	subcom[opt.Name] = fn
	return opt
}

func (b *bot) regMangaInteraction(
	subcom map[string]func(e Interaction),
	s SearchProvider,
) discord.CommandOption {
	opt := discord.CommandOption{
		Name:        "manga",
		Description: "search for a manga",
		Type:        discord.SubcommandOption,
		Required:    false,
		Options: []discord.CommandOption{{
			Type:        discord.StringOption,
			Name:        "title",
			Description: "title of the manga",
			Required:    true,
		}},
	}

	fn := func(i Interaction) {
		log.Trace().Interface("options", i.Options).Msg("Searching for a manga")

		manga, err := s.Manga(i.Options[0].Value)
		if err != nil {
			_ = b.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: "Error searching for this manga, either it doesn't exist or something went wrong",
				},
			})
			log.Err(err).Str("title", i.Options[0].Value).Msg("Error searching for the manga")
		}

		if len(manga) < 1 {
			_ = b.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: "Error searching for this manga, either it doesn't exist or something went wrong",
				},
			})
			return
		}

		err = b.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
			Type: api.MessageInteractionWithSource,
			Data: &api.InteractionResponseData{
				Embeds: []discord.Embed{buildMediaEmbed(manga[0])},
			},
		})
		if err != nil {
			log.Err(err).Msg("failed to send interaction callback")
		}
	}

	subcom[opt.Name] = fn
	return opt
}

func (b *bot) regUserInteraction(
	subcom map[string]func(e Interaction),
	s SearchProvider,
) discord.CommandOption {
	opt := discord.CommandOption{
		Name:        "user",
		Description: "search for a user",
		Type:        discord.SubcommandOption,
		Required:    false,
		Options: []discord.CommandOption{{
			Type:        discord.StringOption,
			Name:        "name",
			Description: "name of the user",
			Required:    true,
		}},
	}

	fn := func(i Interaction) {
		log.Trace().Interface("options", i.Options).Msg("Searching for user")

		user, err := s.User(i.Options[0].Value)
		if err != nil {
			_ = b.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: "Error searching for this user, either it doesn't exist or something went wrong",
				},
			})
			log.Err(err).Str("name", i.Options[0].Value).Msg("Error searching for the user")
		}

		if len(user) < 1 {
			_ = b.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: "Error searching for this user, either it doesn't exist or something went wrong",
				},
			})
			return
		}

		err = b.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
			Type: api.MessageInteractionWithSource,
			Data: &api.InteractionResponseData{
				Embeds: []discord.Embed{buildUserEmbed(user[0])},
			},
		})
		if err != nil {
			log.Err(err).Msg("failed to send interaction callback")
		}
	}

	subcom[opt.Name] = fn
	return opt
}

func (b *bot) regCharInteraction(
	subcom map[string]func(e Interaction),
	s SearchProvider,
) discord.CommandOption {
	opt := discord.CommandOption{
		Name:        "char",
		Description: "search for a character",
		Type:        discord.SubcommandOption,
		Required:    false,
		Options: []discord.CommandOption{{
			Type:        discord.StringOption,
			Name:        "name",
			Description: "name of the character",
			Required:    true,
		}},
	}

	fn := func(i Interaction) {
		log.Trace().Interface("options", i.Options).Msg("Searching for char")

		char, err := s.Character(i.Options[0].Value)
		if err != nil {
			_ = b.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: "Error searching for this char, either it doesn't exist or something went wrong",
				},
			})
			log.Err(err).Str("name", i.Options[0].Value).Msg("Error searching for the char")
		}

		if len(char) < 1 {
			_ = b.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: "Error searching for this char, either it doesn't exist or something went wrong",
				},
			})
			return
		}

		err = b.s.RespondInteraction(i.ID, i.Token, api.InteractionResponse{
			Type: api.MessageInteractionWithSource,
			Data: &api.InteractionResponseData{
				Embeds: []discord.Embed{buildCharEmbed(char[0])},
			},
		})
		if err != nil {
			log.Err(err).Msg("failed to send interaction callback")
		}
	}

	subcom[opt.Name] = fn
	return opt
}

func buildMediaEmbed(m search.Media) discord.Embed {
	return discord.Embed{
		Title:       FixString(m.Title.Romaji),
		Description: Sanitize(m.Description),
		URL:         m.Siteurl,
		Thumbnail:   &discord.EmbedThumbnail{URL: m.CoverImage.Large},
		// Anilist blue
		Color: discord.Color(ColorToInt(m.CoverImage.Color)),
		Footer: &discord.EmbedFooter{
			Text: "View on anilist",
			Icon: "https://anilist.co/img/icons/favicon-32x32.png",
		},
		Image: &discord.EmbedImage{URL: m.BannerImage},
	}
}

func buildUserEmbed(u search.User) discord.Embed {
	return discord.Embed{
		Title:       u.Name,
		Description: Sanitize(u.About),
		URL:         u.Siteurl,
		Thumbnail:   &discord.EmbedThumbnail{URL: u.Avatar.Large},
		// Anilist blue
		Color: 0x19212d,
		Footer: &discord.EmbedFooter{
			Text: "View on anilist",
			Icon: "https://anilist.co/img/icons/favicon-32x32.png",
		},
		Image: &discord.EmbedImage{
			URL: fmt.Sprintf("https://img.anili.st/user/%d", u.ID),
		},
	}
}

func buildCharEmbed(c search.Character) discord.Embed {
	return discord.Embed{
		Title:       FixString(c.Name.Full),
		Description: Sanitize(c.Description),
		URL:         c.Siteurl,
		Thumbnail:   &discord.EmbedThumbnail{URL: c.Image.Large},
		// Anilist blue
		Color: 0x19212d,
		Footer: &discord.EmbedFooter{
			Text: "View on anilist",
			Icon: "https://anilist.co/img/icons/favicon-32x32.png",
		},
	}
}

// SanitizeHTML removes all HTML tags from the given string.
// It also removes double newlines and double || characters
var SanitizeHTML = regexp.MustCompile(`<[^>]*>|\|\|[^|]*\|\||\s{2,}|img[\d\%]*\([^)]*\)|[#~*]{2,}|\n`)

// Sanitize removes all HTML tags from the given string.
// It also removes double newlines and double || characters
func Sanitize(s string) string {
	return SanitizeHTML.ReplaceAllString(s, " ")
}

// FixString removes eventual
// double space or any whitespace possibly in a string
// and replace it with a space
func FixString(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// ColorToInt
// Turn an hex color string beginning with a # into a uint32 representing a color
func ColorToInt(s string) uint32 {
	s = strings.Trim(s, "#")

	u, _ := strconv.ParseUint(s, 16, 32)
	return uint32(u)
}
