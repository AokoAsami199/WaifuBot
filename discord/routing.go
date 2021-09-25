package discord

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/rs/zerolog/log"
)

type Commander interface {
	Type() string
}

type slashCmd = func(e *gateway.InteractionCreateEvent, d discord.CommandInteractionData) api.InteractionResponseData

type Command struct {
	fn  slashCmd
	cmd api.CreateCommandData
}

func (Command) Type() string {
	return "Command"
}

type SubCommand map[string]Commander

func (SubCommand) Type() string {
	return "SubCommand"
}

// Route all the commands
func (b *bot) route() {
	for e := range b.eventChan {
		log.Trace().Interface("event", e).Msg("received event")

		d, iok := e.Data.(*discord.CommandInteractionData)
		cmd, cok := b.commands[d.Name]
		if !cok || !iok {
			log.Debug().Str("ID", e.ID.String()).Msg("unknown command")
			b.s.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: option.NewNullableString("Unknown command, the bot works, discord is just being weird, please retry later"),
					Flags:   api.EphemeralResponse,
				},
			})
			continue
		}

		var fn slashCmd

		// Check for subcommands nested
	outer:
		for i := 0; i <= len(d.Options); i++ {
			switch c := cmd.(type) {
			case Command:
				fn = c.fn

				// set the options as the actual command options
				// and not just subcommand names for routing
				if i-1 >= 0 {
					d.Options = d.Options[i-1].Options
				}

				break outer
			case SubCommand:
				if i >= len(d.Options) {
					b.s.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
						Type: api.MessageInteractionWithSource,
						Data: &api.InteractionResponseData{
							Content: option.NewNullableString(
								"Unknown command, the bot works, discord is just being weird, please retry later",
							),
							Flags: api.EphemeralResponse,
						},
					})
					return
				}
				command, ok := c[d.Options[i].Name]

				if !ok {
					b.s.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
						Type: api.MessageInteractionWithSource,
						Data: &api.InteractionResponseData{
							Content: option.NewNullableString(
								"Unknown command, the bot works, discord is just being weird, please retry later",
							),
							Flags: api.EphemeralResponse,
						},
					})
					return
				}

				cmd = command
			}
		}

		// for some reason discord has issues, and manages to route weird stuff.
		// this is a hack to prevent a panic
		if fn == nil {
			b.s.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: option.NewNullableString(
						"Unknown command, you shouldn't be able to use that. Please retry later",
					),
					Flags: api.EphemeralResponse,
				},
			})
			return
		}
		response := fn(e, *d)

		// Respond to the interaction
		b.s.RespondInteraction(e.ID, e.Token, api.InteractionResponse{
			Type: api.MessageInteractionWithSource,
			Data: &response,
		})

	}
}

func (b *bot) Register(commands map[string]Commander) {
	for name, cmd := range commands {
		_, err := b.s.CreateCommand(
			b.appID,
			b.mapToDiscordCommand(name, cmd),
		)
		if err != nil {
			log.Err(err).Interface("command", cmd).Msg("Error registering command")
		}
	}
}

func (b *bot) RegisterGuild(guildID discord.GuildID, commands map[string]Commander) {
	for name, cmd := range commands {
		_, err := b.s.CreateGuildCommand(
			b.appID,
			guildID,
			b.mapToDiscordCommand(name, cmd),
		)
		if err != nil {
			log.Err(err).Interface("command", cmd).Msg("Error registering command")
		}
	}
}

// Turns the map[string]Commander to a discord payload object, so it can be recorded within the API
func (b *bot) mapToDiscordCommand(name string, cmd Commander) api.CreateCommandData {
	var com Command

outer:
	for {
		switch c := cmd.(type) {
		case Command:
			com = c
			b.commands[c.cmd.Name] = com
			break outer

		case SubCommand:
			com = Command{
				cmd: api.CreateCommandData{
					Name:        name,
					Description: name,
					Options:     []discord.CommandOption{},
				},
			}
			for name, sub := range c {
				opt := discord.CommandOption{}
				switch comm := sub.(type) {
				case Command:
					opt = discord.CommandOption{
						Name:        comm.cmd.Name,
						Description: comm.cmd.Description,
						Options:     comm.cmd.Options,
						Type:        discord.SubcommandOption,
					}
				case SubCommand:
					opt = discord.CommandOption{
						Name:        name,
						Description: name,
						Options:     []discord.CommandOption{},
						Type:        discord.SubcommandGroupOption,
					}
					for _, sub := range comm {
						switch comm := sub.(type) {
						case Command:
							opt.Options = append(opt.Options, discord.CommandOption{
								Name:        comm.cmd.Name,
								Description: comm.cmd.Description,
								Options:     comm.cmd.Options,
								Type:        discord.SubcommandOption,
							})
						default:
							log.Warn().Interface("cmd", comm).Msg("We shouldn't be here")
						}
					}
				}
				com.cmd.Options = append(com.cmd.Options, opt)
			}

			b.commands[name] = c

			cmd = com
		}
	}

	return com.cmd
}
