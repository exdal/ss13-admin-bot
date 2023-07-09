package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CommandCallback func(session *discordgo.Session, message *discordgo.MessageCreate, args []string) string
type Command struct {
	Name     string
	Args     []string
	Help     string
	Callback CommandCallback
}

var (
	discordSession *discordgo.Session
	commands       []Command
)

func init() {
	RegisterCommand("ping", nil, "pong", CommandPing)
}

func DiscordSessionStart(token string) bool {
	newDiscordSession, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("error creating Discord session, %v\n", err)
		return false
	}

	newDiscordSession.AddHandler(OnMessage)
	newDiscordSession.Identify.Intents = discordgo.IntentsGuildMessages

	err = newDiscordSession.Open()
	if err != nil {
		log.Fatalf("error opening connection, %v\n", err)
		return false
	}

	discordSession = newDiscordSession

	return true
}

func DiscordSessionEnd() {
	discordSession.Close()
}

func RegisterCommand(name string, args []string, help string, callback CommandCallback) {
	commands = append(commands, Command{name, args, help, callback})
}

func ProcessCommand(session *discordgo.Session, message *discordgo.MessageCreate) string {
	if len(message.Content) != 0 && !strings.HasPrefix(message.Content, discordConfig.Prefix) {
		return ""
	}

	if !IsVisibleChannel(message.ChannelID) {
		return ""
	}

	clearMessage := strings.TrimPrefix(message.Content, discordConfig.Prefix)
	args := strings.Split(clearMessage, " ")
	if len(args) < 1 {
		return "Invalid message"
	}

	for _, command := range commands {
		if command.Name == args[0] {
			return command.Callback(session, message, args)
		}
	}

	return fmt.Sprintf("Unknown command, use %shelp for list of commands.", discordConfig.Prefix)
}

func CommandPing(session *discordgo.Session, message *discordgo.MessageCreate, args []string) string {
	return "pong"
}
