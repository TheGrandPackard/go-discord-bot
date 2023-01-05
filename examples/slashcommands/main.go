package main

import (
	"os"
	"os/signal"
	"time"

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
		DiscordToken:           "<DISCORD_TOKEN_GOES_HERE>",
		GuildID:                "",
		Intents:                discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates,
		CommandDeletionTimeout: time.Second * 30,
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

	// Register Slash Commands
	bot.RegisterSlashCommands()

	// Capture Ctrl-c to shut down bot
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	// Unregister Slash Commands
	bot.UnregisterSlashCommands()

	log.Info().Msg("Gracefully shutting down")
}
