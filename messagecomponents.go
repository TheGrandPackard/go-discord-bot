package godiscordbot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Map message component handlers for their alias
func (d *DiscordBot) MapMessageComponentHandlers(handlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	for alias, handler := range handlers {
		d.messageComponentHandlers[alias] = handler
	}
}

// Process messages to detect /command usage, and invoke the mapped command's handler when applicable
func (d *DiscordBot) messageComponentProcessor(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		// Parse message component action and ID
		action := strings.Split(i.MessageComponentData().CustomID, ":")[0]

		if h, ok := d.messageComponentHandlers[action]; ok {
			h(s, i)
		}
	case discordgo.InteractionModalSubmit:
		// Parse modal submit action and ID
		action := strings.Split(i.ModalSubmitData().CustomID, ":")[0]

		if h, ok := d.messageComponentHandlers[action]; ok {
			h(s, i)
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		data := i.ApplicationCommandData()

		if h, ok := d.messageComponentHandlers[data.Name]; ok {
			h(s, i)
		}
	}
}
