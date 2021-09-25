package discord

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Karitham/WaifuBot/service/anilist"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/rs/zerolog/log"
)

//go:generate go-mockgen -f github.com/Karitham/WaifuBot/discord -i animeSearcher -o animeSearchMock_test.go
type animeSearcher interface {
	Anime(string) ([]anilist.Media, error)
}

func (b *bot) SearchAnime(s animeSearcher) slashCmd {
	return func(e *gateway.InteractionCreateEvent, d discord.CommandInteractionData) api.InteractionResponseData {
		if len(d.Options) < 1 {
			return api.InteractionResponseData{
				Content: option.NewNullableString(
					"Please provide a search term",
				),
				Flags: api.EphemeralResponse,
			}
		}
		log.Trace().Interface("d", d).Msg("Searching for anime")

		anime, err := s.Anime(d.Options[0].String())
		if err != nil {
			return api.InteractionResponseData{
				Content: option.NewNullableString("Error searching for this anime, either it doesn't exist or something went wrong"),
			}
		}

		return api.InteractionResponseData{
			Embeds: &[]discord.Embed{buildMediaEmbed(anime[0])},
		}
	}
}

//go:generate go-mockgen -f github.com/Karitham/WaifuBot/discord -i mangaSearcher -o mangaSearchMock_test.go
type mangaSearcher interface {
	Manga(string) ([]anilist.Media, error)
}

func (b *bot) SearchManga(s mangaSearcher) slashCmd {
	return func(e *gateway.InteractionCreateEvent, d discord.CommandInteractionData) api.InteractionResponseData {
		if len(d.Options) < 1 {
			return api.InteractionResponseData{
				Content: option.NewNullableString(
					"Please provide a search term",
				),
				Flags: api.EphemeralResponse,
			}
		}
		manga, err := s.Manga(d.Options[0].String())

		if err != nil || len(manga) < 1 {
			return api.InteractionResponseData{
				Content: option.NewNullableString(
					"Error searching for this manga, either it doesn't exist or something went wrong",
				),
				Flags: api.EphemeralResponse,
			}
		}

		return api.InteractionResponseData{
			Embeds: &[]discord.Embed{buildMediaEmbed(manga[0])},
		}
	}
}

//go:generate go-mockgen -f github.com/Karitham/WaifuBot/discord -i userSearcher -o userSearchMock_test.go
type userSearcher interface {
	User(string) ([]anilist.User, error)
}

func (b *bot) SearchUser(s userSearcher) slashCmd {
	return func(e *gateway.InteractionCreateEvent, d discord.CommandInteractionData) api.InteractionResponseData {
		if len(d.Options) < 1 {
			return api.InteractionResponseData{
				Content: option.NewNullableString(
					"Please provide a search term",
				),
				Flags: api.EphemeralResponse,
			}
		}
		user, err := s.User(d.Options[0].String())

		if err != nil && len(user) < 1 {
			return api.InteractionResponseData{
				Content: option.NewNullableString(
					"Error searching for this user, either it doesn't exist or something went wrong",
				),
				Flags: api.EphemeralResponse,
			}
		}

		return api.InteractionResponseData{
			Embeds: &[]discord.Embed{buildUserEmbed(user[0])},
		}
	}
}

//go:generate go-mockgen -f github.com/Karitham/WaifuBot/discord -i charSearcher -o charSearchMock_test.go
type charSearcher interface {
	Character(string) ([]anilist.Character, error)
}

func (b *bot) SearchChar(s charSearcher) slashCmd {
	return func(e *gateway.InteractionCreateEvent, d discord.CommandInteractionData) api.InteractionResponseData {
		if len(d.Options) < 1 {
			return api.InteractionResponseData{
				Content: option.NewNullableString(
					"Please provide a search term",
				),
				Flags: api.EphemeralResponse,
			}
		}
		char, err := s.Character(d.Options[0].String())

		if err != nil || len(char) < 1 {
			return api.InteractionResponseData{
				Content: option.NewNullableString(
					"Error searching for this character, either it doesn't exist or something went wrong",
				),
				Flags: api.EphemeralResponse,
			}
		}

		return api.InteractionResponseData{
			Embeds: &[]discord.Embed{buildCharEmbed(char[0])},
		}
	}
}

func buildMediaEmbed(m anilist.Media) discord.Embed {
	return discord.Embed{
		Title:       FixString(m.Title.Romaji),
		Description: Sanitize(m.Description),
		URL:         m.Siteurl,
		Thumbnail:   &discord.EmbedThumbnail{URL: m.CoverImage.Large},
		// Anilist blue
		Color: discord.Color(ColorToInt(m.CoverImage.Color)),
		Footer: &discord.EmbedFooter{
			Text: "View on anilist",
			Icon: anilist.IconURL,
		},
		Image: &discord.EmbedImage{URL: m.BannerImage},
	}
}

func buildUserEmbed(u anilist.User) discord.Embed {
	return discord.Embed{
		Title:       u.Name,
		Description: Sanitize(u.About),
		URL:         u.Siteurl,
		Thumbnail:   &discord.EmbedThumbnail{URL: u.Avatar.Large},
		// Anilist blue
		Color: anilist.Color,
		Footer: &discord.EmbedFooter{
			Text: "View on anilist",
			Icon: anilist.IconURL,
		},
		Image: &discord.EmbedImage{
			URL: fmt.Sprintf("https://img.anili.st/user/%d", u.ID),
		},
	}
}

func buildCharEmbed(c anilist.Character) discord.Embed {
	return discord.Embed{
		Title:       FixString(c.Name.Full),
		Description: Sanitize(c.Description),
		URL:         c.Siteurl,
		Thumbnail:   &discord.EmbedThumbnail{URL: c.Image.Large},
		// Anilist blue
		Color: anilist.Color,
		Footer: &discord.EmbedFooter{
			Text: "View on anilist",
			Icon: anilist.IconURL,
		},
	}
}

// SanitizeHTML removes all HTML tags from the given string.
// It also removes double newlines and double || characters.
var SanitizeHTML = regexp.MustCompile(`<[^>]*>|\|\|[^|]*\|\||\s{2,}|img[\d\%]*\([^)]*\)|[#~*]{2,}|\n`)

// Sanitize removes all HTML tags from the given string.
// It also removes double newlines and double || characters.
func Sanitize(s string) string {
	return SanitizeHTML.ReplaceAllString(s, " ")
}

// FixString removes eventual
// double space or any whitespace possibly in a string
// and replace it with a space.
func FixString(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// ColorToInt
// Turn an hex color string beginning with a # into a uint32 representing a color.
func ColorToInt(s string) uint32 {
	s = strings.Trim(s, "#")

	u, _ := strconv.ParseUint(s, 16, 32)
	return uint32(u)
}
