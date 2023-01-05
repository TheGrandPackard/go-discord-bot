package godiscordbot

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Options struct {
	DiscordToken           string
	GuildID                string
	Intents                discordgo.Intent
	CommandDeletionTimeout time.Duration
	CommandResponseTimeout time.Duration
	LegacyCommandPrefix    string
}

type DiscordBot struct {
	s       *discordgo.Session
	GuildID string

	// Configurations
	CommandDeletionTimeout time.Duration
	CommandResponseTimeout time.Duration

	// Legacy Commands
	legacyCommandPrefix   string
	legacyCommandHandlers map[string]func(d *DiscordBot, m *discordgo.MessageCreate, arguments []string)

	// Slash Commands
	slashCommandHandlers    map[*discordgo.ApplicationCommand]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	registeredSlashCommands []*discordgo.ApplicationCommand

	// Message Components
	messageComponentHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func New(options Options) (*DiscordBot, error) {
	bot := &DiscordBot{
		GuildID:                options.GuildID,
		CommandDeletionTimeout: options.CommandDeletionTimeout,
		CommandResponseTimeout: options.CommandResponseTimeout,

		// Handlers
		legacyCommandPrefix:      options.LegacyCommandPrefix,
		legacyCommandHandlers:    map[string]func(d *DiscordBot, m *discordgo.MessageCreate, arguments []string){},
		slashCommandHandlers:     map[*discordgo.ApplicationCommand]func(s *discordgo.Session, i *discordgo.InteractionCreate){},
		messageComponentHandlers: map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){},
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
	// Declare intents necessary for reading messages and interacting with voice channels
	bot.s.Identify.Intents = options.Intents

	return bot, nil
}

func (d *DiscordBot) Start() error {
	// Open websocket to Discord API
	return d.s.Open()
}

func (d *DiscordBot) Stop() error {
	return d.s.Close()
}
