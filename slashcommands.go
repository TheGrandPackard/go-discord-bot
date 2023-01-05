package godiscordbot

import (
	"github.com/rs/zerolog/log"

	"github.com/bwmarrin/discordgo"
)

// Map slash command handlers for their alias
func (d *DiscordBot) MapSlashCommands(commands map[*discordgo.ApplicationCommand]func(d *DiscordBot, i *discordgo.InteractionCreate)) {
	for command, handler := range commands {
		d.slashCommandHandlers[command] = handler
	}
}

// Process messages to detect /command usage, and invoke the mapped command's handler when applicable
func (d *DiscordBot) slashCommandProcessor(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		// This would be more efficient as a map instead of an array, but is simpler this way for how few commands there are to process
		for cmd, handler := range d.slashCommandHandlers {
			if cmd.Name == i.ApplicationCommandData().Name {
				handler(d, i)
				return
			}
		}
	}
}

// Register slash commands with the discord server, unregistering any existing commands
// Note: this only needs to be done if the ApplicationCommand definition changes, and not the handler
func (d *DiscordBot) RegisterSlashCommands() {
	d.UnregisterSlashCommands()

	log.Info().Msg("Registering Slash Commands")

	for command := range d.slashCommandHandlers {
		cmd, err := d.s.ApplicationCommandCreate(d.s.State.User.ID, d.GuildID, command)
		if err != nil {
			log.Error().Msgf("Cannot create '%v' command: %v", command.Name, err)
		}
		log.Info().Msgf("Registered Slash Command: %s", command.Name)
		d.registeredSlashCommands = append(d.registeredSlashCommands, cmd)
	}
}

// Unregister slash commands with the discord server
func (d *DiscordBot) UnregisterSlashCommands() {
	log.Info().Msg("Unregistering Slash Commands")

	if len(d.registeredSlashCommands) == 0 {
		var err error
		d.registeredSlashCommands, err = d.s.ApplicationCommands(d.s.State.User.ID, d.GuildID)
		if err != nil {
			log.Error().Msgf("Could not fetch registered commands: %v", err)
			return
		}
	}

	for _, v := range d.registeredSlashCommands {
		err := d.s.ApplicationCommandDelete(d.s.State.User.ID, d.GuildID, v.ID)
		if err != nil {
			log.Error().Msgf("Cannot delete '%v' command: %v", v.Name, err)
		}
		log.Info().Msgf("Unregistered Slash Command: %s: %s", v.Name, v.ID)
	}
}
