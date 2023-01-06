package godiscordbot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Map slash command handlers for their alias
func (d *DiscordBot) MapLegacyCommands(commands map[string]func(d *DiscordBot, m *discordgo.MessageCreate, arguments []string)) {
	for alias, command := range commands {
		d.legacyCommandHandlers[alias] = command
	}
}

// Process messages to detect legacy command usage, and invoke the mapped command's handler when applicable
func (d *DiscordBot) legacyCommandProcessor(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.ID == d.s.State.User.ID {
		return
	}

	// Ignore commands without prefix
	if !strings.HasPrefix(m.Content, d.LegacyCommandPrefix) {
		return
	}

	// Clean up user's command message
	d.DeleteMessageWithTimeout(m.ChannelID, m.Message.ID, d.LegacyCommandDeletionTimeout)

	// Match command if exists
	commandParts := strings.Split(m.Content, " ")
	for command := range d.legacyCommandHandlers {
		if strings.ToLower(commandParts[0]) == command {
			d.legacyCommandHandlers[command](d, m, commandParts[1:])
			return
		}
	}

	// If no command exists, respond to the user and log the usage
	d.SendMessageWithDeletionTimeout(m.ChannelID, responseUnknownLegacyCommand, d.LegacyCommandDeletionTimeout)

	log.Info().
		Str("member", m.Author.Username).
		Str("channel", m.ChannelID).
		Str("command", m.Content).
		Msg("Unknown Legacy Command")
}
