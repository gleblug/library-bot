package telegram

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/gleblug/library-bot/clients/telegram"
	"github.com/gleblug/library-bot/lib/e"
	"github.com/gleblug/library-bot/storage"
)

const (
	SaveCmd = "/save"
)

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func (p *Processor) doCmd(text string, chatID int, username string) error {
	sendMsg := NewMessageSender(chatID, p.tg)
	text = strings.TrimSpace(text)

	log.Printf("got new cmd '%s' from '%s'", text, username)

	if isSearchQuery(text) {
		return p.searchBooks(chatID, text)
	}

	switch text {
	default:
		return sendMsg(msgUnknownCommand)
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

func (p *Processor) saveBook(text string, chatID int, filename string, fileID string) (err error) {
	defer func() { err = e.WrapIfErr("can't save book", err) }()
	sendMsg := NewMessageSender(chatID, p.tg)

	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	book := &storage.Book{
		Filename: filename,
		FileID:   fileID,
		Tags:     storage.Analyze(text + " " + name),
	}

	isExists, err := p.storage.IsExists(book)
	if err != nil {
		return err
	}
	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err := p.storage.Save(book); err != nil {
		return err
	}

	log.Printf("'%+v' book saved", book)
	return sendMsg(msgSaved)
}

func (p *Processor) sendBook(chatID int, fileID string) (err error) {
	sendMsg := NewMessageSender(chatID, p.tg)

	err = p.tg.SendDocument(chatID, fileID)
	if err != nil {
		sendMsg(msgCantFind)
		return e.Wrap("can't send book", err)
	}

	return nil
}

func isSearchQuery(text string) bool {
	return (text != "") && (text[0] != '/')
}
