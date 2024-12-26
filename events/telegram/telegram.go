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
	admins  []string // admin usernames
}

type Meta struct {
	ChatID       int
	Username     string
	FileID       string
	Filename     string
	CallbackID   string
	CallbackData string
}

var (
	ErrUnknownEventType = errors.New("unknown message type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage, admins []string) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
		admins:  admins,
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
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process event", err)
	}

	switch event.Type {
	case events.Message:
		return p.processMessage(event.Text, meta.ChatID, meta.Username)
	case events.Document:
		return p.processDocument(event.Text, meta.ChatID, meta.Username, meta.FileID, meta.Filename)
	case events.Callback:
		return p.processCallback(meta.ChatID, meta.CallbackID, meta.CallbackData)
	default:
		return e.Wrap("can't process event", ErrUnknownEventType)
	}
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
			ChatID:       upd.Callback.Message.Chat.ID,
			CallbackID:   upd.Callback.ID,
			CallbackData: upd.Callback.Data,
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
	return events.Unknown
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	} else if upd.Message.Text == "" {
		return upd.Message.Caption
	}
	return upd.Message.Text
}
