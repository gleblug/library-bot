package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/gleblug/library-bot/lib/e"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod     = "getUpdates"
	sendMessageMethod    = "sendMessage"
	sendDocumentMethod   = "sendDocument"
	answerCallbackMethod = "answerCallbackQuery"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return fmt.Sprintf("bot%s", token)
}

func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("can't get updates", err) }()

	jsonStr, err := json.Marshal(UpdateQuery{
		Offset: offset,
		Limit:  limit,
	})

	data, err := c.doRequest(getUpdatesMethod, jsonStr)
	if err != nil {
		return nil, err
	}

	var res UpdateResponce

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	jsonStr, err := json.Marshal(OutcomingMessage{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	_, err = c.doRequest(sendMessageMethod, jsonStr)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) SendMessageWithKeyboard(chatID int, text string, buttons []string) error {
	keyboard := make([][]InlineKeyboardButton, 0, len(buttons))
	for _, btn := range buttons {
		keyboard = append(keyboard, []InlineKeyboardButton{{btn, btn}})
	}

	jsonStr, err := json.Marshal(OutcomingMessage{
		ChatID: chatID,
		Text:   text,
		Markup: &InlineKeyboardMarkup{keyboard},
	})
	if err != nil {
		return e.Wrap("can't send message with keyboard", err)
	}

	_, err = c.doRequest(sendMessageMethod, jsonStr)
	if err != nil {
		return e.Wrap("can't send message with keyboard", err)
	}

	return nil
}

func (c *Client) SendDocument(chatID int, fileID string) error {
	jsonStr, err := json.Marshal(OutcomingDocument{
		ChatID: chatID,
		FileID: fileID,
	})
	if err != nil {
		return e.Wrap("can't send document", err)
	}

	_, err = c.doRequest(sendDocumentMethod, jsonStr)
	if err != nil {
		return e.Wrap("can't send document", err)
	}

	return nil
}

func (c *Client) AnswerCallback(callbackID string) error {
	jsonStr, err := json.Marshal(AnswerCallbackQuery{
		CallbackID: callbackID,
	})
	if err != nil {
		return e.Wrap("can't answer callback", err)
	}

	_, err = c.doRequest(answerCallbackMethod, jsonStr)
	if err != nil {
		return e.Wrap("can't answer callback", err)
	}

	return nil
}

func (c *Client) doRequest(method string, jsonStr []byte) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
