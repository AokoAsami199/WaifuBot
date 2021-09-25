package discord

import (
	"context"

	"github.com/Karitham/WaifuBot/service/anilist"
	"github.com/davecgh/go-spew/spew"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/rs/zerolog/log"
)

type bot struct {
	commands  map[string]Commander
	s         *session.Session
	eventChan chan *gateway.InteractionCreateEvent
	appID     discord.AppID
}

// ListenAndServe registers the commands and wait for events.
func LS(appID, token string) (close func()) {
	s, err := session.New("Bot " + token)
	if err != nil {
		log.Fatal().Err(err).Msg("Session failed")

		return
	}

	b := New(s, appIDFromString(appID))
	if err != nil {
		log.Fatal().Err(err).Msg("Session failed")

		return
	}

	if err = b.Open(context.TODO()); err != nil {
		log.Fatal().Err(err).Msg("failed to open")
	}

	b.Register(b.Commands())

	log.Info().Msg("Gateway connected")

	return func() {
		err := s.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to close")
		}
	}
}

// appID from string returns a discord.AppID from a string.
func appIDFromString(s string) discord.AppID {
	a, _ := discord.ParseSnowflake(s)
	return discord.AppID(a)
}

func New(sess *session.Session, appID discord.AppID) *bot {
	b := &bot{
		s:         sess,
		eventChan: make(chan *gateway.InteractionCreateEvent),
		commands:  make(map[string]Commander),
		appID:     appID,
	}

	sess.AddHandler(b.eventChan)

	sess.AddIntents(gateway.IntentGuilds)
	sess.AddIntents(gateway.IntentGuildMessages)

	return b
}

func (b *bot) Open(ctx context.Context) error {
	err := b.s.Open(ctx)
	go b.route()
	return err
}

func (b *bot) Commands() map[string]Commander {
	anilist := anilist.New()

	return map[string]Commander{
		// TODO: Implement roll with a store
		"roll": Command{
			fn: Unimplemented,
			cmd: api.CreateCommandData{
				Name:        "roll",
				Description: "roll a character, obtain it for yourself",
			},
		},
		// TODO: Implement verify
		"verify": Command{
			fn: Unimplemented,
			cmd: api.CreateCommandData{
				Name:        "verify",
				Description: "check if a user has a character",
				Options: []discord.CommandOption{
					{
						Name:        "charID",
						Description: "ID of the character you to check for",
						Type:        discord.IntegerOption,
						Required:    true,
					},
					{
						Name:        "user",
						Description: "user you want to check the character agaisnt",
						Type:        discord.UserOption,
						Required:    true,
					},
				},
			},
		},
		// TODO: Implement give, and the store
		"give": Command{
			fn: Unimplemented,
			cmd: api.CreateCommandData{
				Name:        "give",
				Description: "give a character to a user",
				Options: []discord.CommandOption{
					{
						Name:        "charID",
						Description: "ID of the character you want to gift",
						Type:        discord.IntegerOption,
						Required:    true,
					},
					{
						Name:        "user",
						Description: "user you want to gift the character to",
						Type:        discord.UserOption,
						Required:    true,
					},
				},
			},
		},
		// TODO: Implement profile and the store
		"profile": SubCommand{
			"edit": SubCommand{
				"favorite": Command{
					fn: Unimplemented,
					cmd: api.CreateCommandData{
						Name:        "favorite",
						Description: "set your favorite character",
						Options: []discord.CommandOption{{
							Name:        "characterID",
							Description: "the id of the character you want to set as your favorite",
							Type:        discord.IntegerOption,
							Required:    true,
						}},
					},
				},
				"quote": Command{
					fn: Unimplemented,
					cmd: api.CreateCommandData{
						Name:        "quote",
						Description: "set your profile quote character",
						Options: []discord.CommandOption{{
							Name:        "quote",
							Description: "the id of the character you want to set as your favorite",
							Type:        discord.StringOption,
							Required:    true,
						}},
					},
				},
			},
			"view": Command{
				fn: Unimplemented,
				cmd: api.CreateCommandData{
					Name:        "view",
					Description: "view the user's profile",
					Options: []discord.CommandOption{
						{
							Name:        "user",
							Type:        discord.UserOption,
							Description: "user to view the profile of",
							Required:    false,
						},
					},
				},
			},
		},
		// TODO: implement list
		"list": Command{
			fn: Unimplemented,
			cmd: api.CreateCommandData{
				Name:        "list",
				Description: "list all characters",
				Options: []discord.CommandOption{
					{
						Name:        "user",
						Description: "user to list the characters of",
						Type:        discord.UserOption,
						Required:    false,
					},
				},
			},
		},
		"search": SubCommand{
			"anime": Command{
				cmd: api.CreateCommandData{
					Name:        "anime",
					Description: "search for an anime",
					Options: []discord.CommandOption{{
						Type:        discord.StringOption,
						Name:        "title",
						Description: "title of the anime",
						Required:    true,
					}},
				},
				fn: b.SearchAnime(anilist),
			},
			"manga": Command{
				cmd: api.CreateCommandData{
					Name:        "manga",
					Description: "search for a manga",
					Options: []discord.CommandOption{{
						Type:        discord.StringOption,
						Name:        "title",
						Description: "title of the manga",
						Required:    true,
					}},
				},
				fn: b.SearchManga(anilist),
			},
			"char": Command{
				cmd: api.CreateCommandData{
					Name:        "char",
					Description: "search for a character",
					Options: []discord.CommandOption{{
						Type:        discord.StringOption,
						Name:        "name",
						Description: "name of the character",
						Required:    true,
					}},
				},
				fn: b.SearchChar(anilist),
			},
			"user": Command{
				cmd: api.CreateCommandData{
					Name:        "user",
					Description: "search for a user",
					Options: []discord.CommandOption{{
						Type:        discord.StringOption,
						Name:        "name",
						Description: "name of the user",
						Required:    true,
					}},
				},
				fn: b.SearchUser(anilist),
			},
		},
	}
}

func Unimplemented(e *gateway.InteractionCreateEvent, d discord.CommandInteractionData) api.InteractionResponseData {
	return api.InteractionResponseData{
		Flags:   api.EphemeralResponse,
		Content: option.NewNullableString("This command is not yet implemented"),
		Embeds: &[]discord.Embed{{
			Title:       "Command",
			Description: "```dump\n" + spew.Sdump(d) + "\n```",
			Timestamp:   discord.NowTimestamp(),
			Color:       discord.Color(0xFF0000),
		}},
	}
}
