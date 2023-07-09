package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type GameInfo struct {
	PlayerCount uint32
}

var (
	gameInfo GameInfo
)

func LoadJsonFile(content any, pathToFile string) bool {
	file, err := os.Open(pathToFile)
	defer file.Close()

	if err != nil {
		log.Fatalf("An error occured while reading %s, error: %v\n", pathToFile, err)
		return false
	}

	fileData, _ := ioutil.ReadAll(file)

	err = json.Unmarshal(fileData, &content)
	if err != nil {
		log.Fatalf("An error occured while parsing json file %s, error: %v\n", pathToFile, err)
		return false
	}

	return true
}

func main() {
	log.SetOutput(os.Stdout)
	rand.Seed(time.Now().Unix())

	if !LoadJsonFile(&discordConfig, "config/discord_config.json") {
		return
	}

	byond_session_start(discordConfig.ByondServer)
	DiscordSessionStart(discordConfig.BotSecret)
	InitTickers()

	log.Println("Started bot.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Stopping bot...")

	DiscordSessionEnd()
	byond_session_shutdown()
}

func OnMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	response := ProcessCommand(session, message)
	if len(response) != 0 {
		session.ChannelMessageSend(message.ChannelID, response)
	}
}
