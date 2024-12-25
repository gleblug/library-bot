package telegram

import (
	"errors"

	"github.com/gleblug/library-bot/clients/telegram"
	"github.com/gleblug/library-bot/events"
	"github.com/gleblug/library-bot/lib/e"
	"github.com/gleblug/library-bot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID     int
	Username   string
	FileID     string
	Filename   string
	CallbackID string
}

var (
	ErrUnknownEventType = errors.New("unknown message type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.Document:
		return p.processDocument(event)
	case events.Callback:
		return p.processCallback(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func (p *Processor) processDocument(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process document", err)
	}

	if err := p.saveBook(event.Text, meta.ChatID, meta.Filename, meta.FileID); err != nil {
		return e.Wrap("can't process document", err)
	}

	return nil
}

func (p *Processor) processCallback(event events.Event) (err error) {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process callback", err)
	}

	book, err := p.storage.Read(meta.Filename)
	if err != nil {
		return e.Wrap("can't process callback", err)
	}

	if err := p.sendBook(meta.ChatID, book.FileID); err != nil {
		return e.Wrap("can't process callback", err)
	}

	p.tg.AnswerCallback(meta.CallbackID)

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	switch updType {
	case events.Message:
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	case events.Document:
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
			FileID:   upd.Message.Document.ID,
			Filename: upd.Message.Document.Filename,
		}
	case events.Callback:
		res.Meta = Meta{
			ChatID:     upd.Callback.Message.Chat.ID,
			Filename:   upd.Callback.Data,
			CallbackID: upd.Callback.ID,
		}
	default:
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Callback != nil {
		return events.Callback
	} else if upd.Message != nil {
		if upd.Message.Document != nil {
			return events.Document
		}
		return events.Message
	}
	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	} else if upd.Message.Text == "" {
		return upd.Message.Caption
	}
	return upd.Message.Text
}
