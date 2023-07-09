package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

type TickerCallback func()
type Ticker struct {
	Name        string
	Interval    time.Duration
	Callback    TickerCallback
	Handle      *time.Ticker
	ShouldClose chan bool
}

var (
	tickers []*Ticker
)

func RegisterTicker(name string, interval time.Duration, callback TickerCallback) {
	ticker := new(Ticker)
	ticker.Name = name
	ticker.Interval = interval
	ticker.Callback = callback
	ticker.Handle = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ticker.ShouldClose:
				return
			case <-ticker.Handle.C:
				callback()
			}
		}
	}()

	tickers = append(tickers, ticker)
}

func InitTickers() {
	RegisterTicker("player_count", 15*time.Second, TopicPlayerCount)
	RegisterTicker("poly_speech", 15*time.Minute, PolySpeak)

	for _, v := range tickers {
		v.Callback()
	}
}

func TopicPlayerCount() {
	_, dataType, data := byond_on_topic("?playing")
	if dataType != TOPIC_TYPE_DECIMAL {
		gameInfo.PlayerCount = 0
	} else {
		gameInfo.PlayerCount = uint32(byond_read_number(data))
	}

	discordSession.UpdateGameStatus(0, fmt.Sprintf("%v oyuncu ile SS13", gameInfo.PlayerCount))
}

type NPCSavePoly struct {
	Phrases []string `json:"phrases,omitempty"`
}

func PolySpeak() {
	speech := new(NPCSavePoly)
	if !LoadJsonFile(speech, discordConfig.Poly) {
		return
	}

	webhookParams := new(discordgo.WebhookParams)
	webhookParams.Content = speech.Phrases[rand.Intn(len(speech.Phrases))]

	id, token := GetWebhook("poly")
	discordSession.WebhookExecute(id, token, false, webhookParams)
}
