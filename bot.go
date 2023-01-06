package godiscordbot

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Options struct {
	// Token to authenticate the bot with the Discord API
	DiscordToken string
	// Unused
	GuildID string

	// Any intents necessary for the bot to operate
	// * discordgo.IntentsGuildMessages - required for the bot to process legacy commands
	// * discordgo.IntentsGuildVoiceStates - required for the bot to join voice channels and play audio
	Intents discordgo.Intent

	// Legacy Commands
	// The timeout before deleting user messages (empty or zero means messages are never deleted)
	LegacyCommandDeletionTimeout time.Duration
	// The timeout before deleting bot responses (empty or zero means messages are never deleted)
	LegacyCommandResponseTimeout time.Duration
	// Prefix for legacy commands, such as "!"
	LegacyCommandPrefix string

	// Slash Commands
	// Whether or not to register slash commands when starting the bot
	RegisterSlashCommands bool
	// Whether or not to unregister slash commands when starting the bot
	UnregisterSlashCommands bool
	// The timeout before deleting bot responses (empty or zero means messages are never deleted)
	SlashCommandResponseTimeout time.Duration
}

type DiscordBot struct {
	s *discordgo.Session

	// Configurations
	GuildID string

	// Legacy Commands
	LegacyCommandDeletionTimeout time.Duration
	LegacyCommandResponseTimeout time.Duration
	LegacyCommandPrefix          string
	legacyCommandHandlers        map[string]func(d *DiscordBot, m *discordgo.MessageCreate, arguments []string)

	// Slash Commands
	RegisterSlashCommands       bool
	UnregisterSlashCommands     bool
	SlashCommandResponseTimeout time.Duration
	slashCommandHandlers        map[*discordgo.ApplicationCommand]func(d *DiscordBot, i *discordgo.InteractionCreate)
	registeredSlashCommands     []*discordgo.ApplicationCommand

	// Message Components
	messageComponentHandlers map[string]func(d *DiscordBot, i *discordgo.InteractionCreate)
}

func New(options Options) (*DiscordBot, error) {
	bot := &DiscordBot{
		// Configurations
		GuildID: options.GuildID,

		// Legacy Commands
		LegacyCommandDeletionTimeout: options.LegacyCommandDeletionTimeout,
		LegacyCommandResponseTimeout: options.LegacyCommandResponseTimeout,
		LegacyCommandPrefix:          options.LegacyCommandPrefix,
		legacyCommandHandlers:        map[string]func(d *DiscordBot, m *discordgo.MessageCreate, arguments []string){},

		// Slash Commands
		RegisterSlashCommands:   options.RegisterSlashCommands,
		UnregisterSlashCommands: options.UnregisterSlashCommands,
		slashCommandHandlers:    map[*discordgo.ApplicationCommand]func(d *DiscordBot, i *discordgo.InteractionCreate){},
		registeredSlashCommands: []*discordgo.ApplicationCommand{},

		// Message Components
		messageComponentHandlers: map[string]func(d *DiscordBot, i *discordgo.InteractionCreate){},
	}

	// Initialize Discord Chat API
	var err error
	bot.s, err = discordgo.New("Bot " + options.DiscordToken)
	if err != nil {
		return nil, err
	}

	// Register Discord API Handlers
	bot.s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Info().Msg("Bot is up!") })
	// Register the legacy command processor func as a callback for legacy command invocation
	bot.s.AddHandler(bot.legacyCommandProcessor)
	// Register the slash command processor func as a callback for /command invocation
	bot.s.AddHandler(bot.slashCommandProcessor)
	// Register the message component processor func as a callback for message component handlers
	bot.s.AddHandler(bot.messageComponentProcessor)
	// Declare intents necessary for the bot (reading messages, interacting with voice channels, etc.)
	bot.s.Identify.Intents = options.Intents

	return bot, nil
}

func (d *DiscordBot) Start() error {
	// Open websocket to Discord API
	err := d.s.Open()
	if err != nil {
		return err
	}

	// Register commands asynchronously
	go d.registerSlashCommands()

	return nil
}

func (d *DiscordBot) Stop() error {
	// Close websocket to Discord API
	err := d.s.Close()
	if err != nil {
		return err
	}

	// Unregister commands synchronously
	d.unregisterSlashCommands()

	return nil
}
