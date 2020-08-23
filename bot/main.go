package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/vitorbaraujo/buschebot/repliers" // used for register custom repliers
	"github.com/vitorbaraujo/buschebot/reply"
)

func main() {
	var apiToken string
	if val, ok := os.LookupEnv("BOT_API_TOKEN"); ok {
		apiToken = val
	}

	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Panicf("creating bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %v", bot.Self.UserName)

	if err := queryMessages(bot); err != nil {
		log.Panicf("querying messages: %v", err)
	}
}

func queryMessages(bot *tgbotapi.BotAPI) error {
	cfg := tgbotapi.NewUpdate(0)
	cfg.Timeout = 60

	updates, err := bot.GetUpdatesChan(cfg)
	if err != nil {
		return err
	}

	// This runs indefinitely.
	for update := range updates {
		if update.Message == nil {
			continue
		}

		response := reply.GetReply(&reply.MessagePayload{
			Text:   update.Message.Text,
			UserId: fmt.Sprint(update.Message.From.ID),
		})

		if response.Text == "" {
			// bot did not come up with an answer.
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response.Text)
		if response.Reply {
			msg.ReplyToMessageID = update.Message.MessageID
		}

		bot.Send(msg)
	}

	return nil
}
