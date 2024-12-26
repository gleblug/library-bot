package telegram

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/gleblug/library-bot/lib/e"
	"github.com/gleblug/library-bot/storage"
)

func (p *Processor) processDocument(text string, chatID int, username string, fileID string, filename string) error {
	if p.role(username) >= Admin {
		return p.saveBook(text, chatID, username, fileID, filename)
	}
	return p.notEnoughRights(chatID)
}

func (p *Processor) saveBook(text string, chatID int, username string, fileID string, filename string) (err error) {
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

	log.Printf("user '%s' save '%+v' book", username, book)
	return sendMsg(msgSaved)
}
