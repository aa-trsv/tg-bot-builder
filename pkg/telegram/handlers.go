package telegram

import (
	"github.com/aa-trsv/telegram-bot-otrs-builder/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

const (
	commandStart = "start"
	buildStage   = "stage"
	rebuildStage = "rebuild"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handlerStartCommand(message)
	case buildStage:
		return b.handlerStageCommand(message)
	case rebuildStage:
		return b.handlerRebuildStageCommand(message)
	default:
		return b.handlerUnknownCommand(message)
	}
}

// TODO: Добавить обработку проверки авторизации. В БД сохранять результат создания сессии, если успешно то rw иначе ro
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	b.bot.Send(msg)
}

func (b *Bot) handlerStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.config.Messages.Start)

	if err := b.accessRepository.Save(message.Chat.ID, "rw", repository.AccessList); err != nil {
		return err
	}

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handlerStageCommand(message *tgbotapi.Message) error {
	err := b.buildPackage()

	path := b.config.DirPath
	filePath, err := b.getFilePath(path)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewDocumentUpload(message.Chat.ID, filePath)

	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}

	err = os.Remove(filePath)
	return err
}

func (b *Bot) handlerRebuildStageCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.config.Messages.Rebuild)

	err := b.rebuildPackage()

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handlerUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.config.Messages.UnknownCommand)

	access, err := b.accessRepository.Get(message.Chat.ID, repository.AccessList)
	if err != nil {
		return err
	}

	msg.Text = access
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) getFilePath(path string) (string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	filePath := path
	for _, file := range files {
		filePath += file.Name()
	}

	return filePath, nil
}

func (b *Bot) buildPackage() error {
	c := exec.Command("ssh", b.config.SSHCommand.Build)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Run()

	return nil
}

func (b *Bot) rebuildPackage() error {
	c := exec.Command("ssh", b.config.SSHCommand.Rebuild)
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Run()

	return nil
}
