package telegram

import (
	"log"
	"strings"

	"github.com/gleblug/library-bot/clients/telegram"
	"github.com/gleblug/library-bot/lib/e"
)

const (
	StartCmd   = "/start"
	HelpCmd    = "/help"
	TopicalCmd = "/topical"
	ChoirCmd   = "/choir"
)

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func (p *Processor) processMessage(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new cmd '%s' from '%s'", text, username)

	if isSearchQuery(text) {
		return p.searchBooks(chatID, text)
	}

	switch text {
	case StartCmd:
		return p.sendHello(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	case TopicalCmd:
		return p.sendTopical(chatID)
	case ChoirCmd:
		return p.sendChoir(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) searchBooks(chatID int, query string) (err error) {
	const count = 1

	defer func() { err = e.WrapIfErr("can't do cmd search", err) }()

	books, err := p.storage.Search(query, count)

	log.Printf("find these: %v", books)

	if err != nil {
		return err
	}

	names := make([]string, 0, len(books))
	for _, book := range books {
		names = append(names, book.Filename)
	}

	return p.tg.SendMessageWithKeyboard(chatID, msgWhitImFind, names)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendTopical(chatID int) error {
	return nil // TODO: implement
}

func (p *Processor) sendChoir(chatID int) error {
	return p.tg.SendMessage(chatID, msgChoirgGroup)
}

func isSearchQuery(text string) bool {
	return (text != "") && (text[0] != '/')
}
