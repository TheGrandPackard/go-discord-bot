package main

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/thegrandpackard/godiscordbot"
)

var (
	commands = map[*discordgo.ApplicationCommand]func(d *godiscordbot.DiscordBot, i *discordgo.InteractionCreate){
		// Help Command
		{
			Name:        "help",
			Description: "Get Bot Help",
		}: func(d *godiscordbot.DiscordBot, i *discordgo.InteractionCreate) {
			log.Info().
				Str("member", i.Member.User.Username).
				Str("channel", i.ChannelID).
				Str("command", "/help").
				Msg("Slash Command Received")

			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Help command response goes here.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		},
		// Info Command
		{
			Name:        "info",
			Description: "Get Bot Info",
		}: func(d *godiscordbot.DiscordBot, i *discordgo.InteractionCreate) {
			log.Info().
				Str("member", i.Member.User.Username).
				Str("channel", i.ChannelID).
				Str("command", "/info").
				Msg("Slash Command Received")

			d.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This is an example discord bot.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		},
	}
)

func main() {
	bot, err := godiscordbot.New(godiscordbot.Options{
		DiscordToken:            "<DISCORD_BOT_TOKEN>",
		RegisterSlashCommands:   true,
		UnregisterSlashCommands: true,
	})
	if err != nil {
		panic(err)
	}

	// Map slash commands
	bot.MapSlashCommands(commands)

	// Start the bot
	err = bot.Start()
	if err != nil {
		panic(err)
	}
	defer bot.Stop()

	// Capture Ctrl-c to shut down bot
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Info().Msg("Gracefully shutting down")
}
