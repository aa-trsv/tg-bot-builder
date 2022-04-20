package telegram

import (
	"github.com/aa-trsv/telegram-bot-otrs-builder/pkg/config"
	"github.com/aa-trsv/telegram-bot-otrs-builder/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Bot struct {
	bot              *tgbotapi.BotAPI
	accessRepository repository.AccessRepository
	config           config.Config
}

func NewBot(bot *tgbotapi.BotAPI, accessRepository repository.AccessRepository, cfg config.Config) *Bot {
	return &Bot{bot: bot, accessRepository: accessRepository, config: cfg}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				log.Println(err)
			}
			continue
		}

		b.handleMessage(update.Message)
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
