package discord

import (
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/rs/zerolog/log"
)

// Command is a discord command.
type Command struct {
	callback   func(e *gateway.InteractionCreateEvent)
	apiCommand api.CreateCommandData
}

func (b *bot) Register(commands ...Command) {
	for _, command := range commands {
		c, err := b.s.CreateCommand(b.appID, command.apiCommand)
		if err != nil {
			log.Err(err).Interface("command", command).Msg("failed to register")
			continue
		}

		b.commands[c.ID] = command.callback
	}
}

func (b *bot) RegisterGuild(guildID discord.GuildID, commands ...Command) {
	for _, command := range commands {
		c, err := b.s.CreateGuildCommand(b.appID, guildID, command.apiCommand)
		if err != nil {
			log.Err(err).Interface("command", command).Msg("failed to register")
			continue
		}

		b.commands[c.ID] = command.callback
	}
}

func (b *bot) route() {
	for e := range b.eventChan {
		fn, ok := b.commands[e.Data.ID]
		if ok {
			fn(e)
		}

		log.Debug().Str("ID", e.Data.ID.String()).Str("name", e.Data.Name).Msg("unknown command")

		_ = b.s.DeleteCommand(b.appID, e.Data.ID)
		_ = b.s.DeleteGuildCommand(b.appID, e.GuildID, e.Data.ID)
	}
}
