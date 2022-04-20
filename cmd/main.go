package main

import (
	"github.com/aa-trsv/telegram-bot-otrs-builder/pkg/config"
	"github.com/aa-trsv/telegram-bot-otrs-builder/pkg/repository"
	"github.com/aa-trsv/telegram-bot-otrs-builder/pkg/repository/boltdb"
	"github.com/aa-trsv/telegram-bot-otrs-builder/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = cfg.Debug

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	accessRepository := boltdb.NewAccessRepository(db)

	telegramBot := telegram.NewBot(bot, accessRepository, *cfg)
	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessList))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
