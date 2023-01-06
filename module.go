package godiscordbot

import "github.com/bwmarrin/discordgo"

// Module defines the different handlers that a discord bot feature would implement
type Module interface {
	GetLegacyCommandHandlers() map[string]func(d *DiscordBot, m *discordgo.MessageCreate, arguments []string)
	GetSlashCommandHandlers() map[*discordgo.ApplicationCommand]func(d *DiscordBot, i *discordgo.InteractionCreate)
	GetMessageComponentHandlers() map[string]func(d *DiscordBot, i *discordgo.InteractionCreate)
}

// LoadModule is a helper to map any handlers defined by a module
func (d *DiscordBot) LoadModule(module Module) {
	d.MapLegacyCommands(module.GetLegacyCommandHandlers())
	d.MapSlashCommands(module.GetSlashCommandHandlers())
	d.MapMessageComponentHandlers(module.GetMessageComponentHandlers())
}
