package telegram

import (
	"github.com/gleblug/library-bot/lib/e"
)

func (p *Processor) processCallback(chatID int, callbackID string, data string) error {
	book, err := p.storage.Read(data)
	if err != nil {
		return err
	}

	if err := p.sendBook(chatID, book.FileID); err != nil {
		return err
	}

	p.tg.AnswerCallback(callbackID)

	return nil
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
