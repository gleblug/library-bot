package telegram

import "slices"

type Role int

const (
	User Role = iota
	Admin
)

func (p *Processor) role(username string) Role {
	if slices.Contains(p.admins, username) {
		return Admin
	}
	return User
}

func (p *Processor) notEnoughRights(chatID int) error {
	return p.tg.SendMessage(chatID, msgNotEnoughRights)
}
