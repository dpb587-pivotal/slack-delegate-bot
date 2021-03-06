package handler

import "github.com/dpb587/slack-delegate-bot/cmd/delegatebot/message"

//go:generate counterfeiter . Handler
type Handler interface {
	IsApplicable(message.Message) (bool, error)
	Execute(*message.Message) (MessageResponse, error)
}
