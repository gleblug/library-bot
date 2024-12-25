package telegram

type UpdateQuery struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type AnswerCallbackQuery struct {
	CallbackID string `json:"callback_query_id"`
}

type UpdateResponce struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID       int              `json:"update_id"`
	Message  *IncomingMessage `json:"message"`
	Callback *CallbackQuery   `json:"callback_query"`
}

type IncomingMessage struct {
	Text     string    `json:"text"`
	From     User      `json:"from"`
	Chat     Chat      `json:"chat"`
	Document *Document `json:"document"`
	Caption  string    `json:"caption"`
}

type Document struct {
	ID       string `json:"file_id"`
	Filename string `json:"file_name"`
}

type User struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type CallbackQuery struct {
	ID      string           `json:"id"`
	Message *IncomingMessage `json:"message"`
	Data    string           `json:"data"`
}

type CallbackMessage struct {
	Chat Chat `json:"chat"`
}

// outcoming

type OutcomingMessage struct {
	ChatID int                   `json:"chat_id"`
	Text   string                `json:"text"`
	Markup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type InlineKeyboardMarkup struct {
	Keyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text     string `json:"text"`
	Callback string `json:"callback_data"`
}

type OutcomingDocument struct {
	ChatID int    `json:"chat_id"`
	FileID string `json:"document"`
}
