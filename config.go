package main

import (
	"log"
)

// Contents of discord_config.json
type DiscordConfig struct {
	BotSecret   string `json:"bot_secret"`
	Prefix      string `json:"command_prefix"`
	ByondServer string `json:"byond_server"`
	Channels    []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"channels,omitempty"`
	Webhooks []struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		Token string `json:"token"`
	} `json:"webhooks,omitempty"`
	Poly string `json:"poly"`
}

var (
	discordConfig DiscordConfig
)

func GetWebhook(name string) (string, string) {
	for _, webhook := range discordConfig.Webhooks {
		if webhook.Name == name {
			return webhook.ID, webhook.Token
		}
	}

	log.Fatalf("Trying to get invalid webhook URL for %s!\n", name)

	return "<UNK>", "<UNK>"
}

func IsVisibleChannel(id string) bool {
	for _, channel := range discordConfig.Channels {
		if channel.ID == id {
			return true
		}
	}

	return false
}
