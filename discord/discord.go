package discord

import (
	"github.com/Karitham/WaifuBot/service/anilist"
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/session"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type bot struct {
	commands  map[discord.CommandID]func(e *gateway.InteractionCreateEvent)
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

	if err = b.Open(); err != nil {
		log.Fatal().Err(err).Msg("failed to open")
	}

	b.Register(b.Search(anilist.New()))
	b.Register()

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
		commands:  make(map[discord.CommandID]func(e *gateway.InteractionCreateEvent)),
		appID:     appID,
	}

	sess.AddHandler(b.eventChan)

	sess.Gateway.AddIntents(gateway.IntentGuilds)
	sess.Gateway.AddIntents(gateway.IntentGuildMessages)

	return b
}

func (b *bot) Open() error {
	err := b.s.Open()
	go b.route()

	return err
}

func (b *bot) RespondWithEmbed(e *gateway.InteractionCreateEvent, embed discord.Embed) *zerolog.Event {
	err := b.s.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Embeds: []discord.Embed{embed},
		},
	})
	if err != nil {
		return log.Err(err)
	}
	return log.Trace()
}

func (b *bot) RespondWithMessage(e *gateway.InteractionCreateEvent, s string) *zerolog.Event {
	err := b.s.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: "Error searching for this char, either it doesn't exist or something went wrong",
		},
	})
	if err != nil {
		return log.Err(err)
	}

	return log.Trace()
}

func (b *bot) RespondWithError(e *gateway.InteractionCreateEvent, s string, err error) *zerolog.Event {
	err2 := b.s.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: "Error searching for this char, either it doesn't exist or something went wrong",
		},
	})
	if err2 != nil {
		return log.Warn().AnErr("default_err", err).AnErr("response_err", err2)
	}

	if err == nil {
		return log.Trace()
	}

	return log.Debug().Err(err)
}
