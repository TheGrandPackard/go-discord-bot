package godiscordbot

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// Delete a message from a channel after a duration
func (d *DiscordBot) DeleteMessageWithTimeout(channelID, messageID string, timeout time.Duration) {
	// Skip delete if timeout is empty
	if timeout == 0 {
		return
	}

	// Delete message after specified timeout
	// TODO: Persist the message ID and expiry in the database to delete on startup in case the bot restarts or crashes
	time.AfterFunc(timeout, func() {
		if err := d.s.ChannelMessageDelete(channelID, messageID); err != nil {
			log.Error().Stack().Err(err).
				Str("channel", channelID).
				Str("messageID", messageID).
				Msg("Error deleting message")
		}
	})
}

// Send a message with content to the specified channel ID that will automatically delete itself after the timeout
func (d *DiscordBot) SendMessageWithDeletionTimeout(channelID string, content string, timeout time.Duration) error {
	msg, err := d.s.ChannelMessageSend(channelID, validateMessage(content))
	if err != nil {
		log.Error().Stack().Err(err).
			Str("channel", channelID).
			Str("msg", content).
			Msg("Error sending message")
		return err
	}

	// Skip delete if timeout is empty
	if timeout == 0 {
		return nil
	}

	// Delete message after specified timeout
	// TODO: Persist the message ID and expiry in the database to delete on startup in case the bot restarts or crashes
	time.AfterFunc(timeout, func() {
		err = d.s.ChannelMessageDelete(channelID, msg.ID)
		if err != nil {
			log.Error().Stack().Err(err).
				Str("channel", channelID).
				Str("messageID", msg.ID).
				Str("msg", content).
				Msg("Error deleting message")
		}
	})

	return nil
}

// Validate message to ensure that it is not empty and cause an error
func validateMessage(message string) string {
	if message == "" {
		return "\u200B"
	}

	return message
}

// Send an embed message with content to the specified channel ID that will automatically delete itself after the timeout
func (d *DiscordBot) SendMessageEmbedWithDeletionTimeout(channelID string, embed *discordgo.MessageEmbed, timeout time.Duration) error {
	msg, err := d.s.ChannelMessageSendEmbed(channelID, validateEmbedMessage(embed))
	if err != nil {
		log.Error().Stack().Err(err).
			Str("channel", channelID).
			Str("title", embed.Title).
			Msg("Error sending message embed")
		return err
	}

	// Skip delete if timeout is empty
	if timeout == 0 {
		return nil
	}

	// Delete message after specified timeout
	// TODO: Persist the message ID and expiry in the database to delete on startup in case the bot restarts or crashes
	time.AfterFunc(timeout, func() {
		err = d.s.ChannelMessageDelete(channelID, msg.ID)
		if err != nil {
			log.Error().Stack().Err(err).
				Str("channel", channelID).
				Str("messageID", msg.ID).
				Str("title", embed.Title).
				Msg("Error deleting message")
		}
	})

	return nil
}

// Validate embed message fields to ensure that they are not empty and cause an error
func validateEmbedMessage(embedMessage *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	for _, field := range embedMessage.Fields {
		if field.Name == "" {
			field.Name = "\u200B"
		}
		if field.Value == "" {
			field.Value = "\u200B"
		}
	}

	return embedMessage
}

// Validate embed message fields to ensure that they are not empty and cause an error
func validateEmbedMessages(embedMessages []*discordgo.MessageEmbed) []*discordgo.MessageEmbed {
	for _, embedMessage := range embedMessages {
		validateEmbedMessage(embedMessage)
	}

	return embedMessages
}

func (d *DiscordBot) InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse) error {
	return d.s.InteractionRespond(interaction, resp)
}
