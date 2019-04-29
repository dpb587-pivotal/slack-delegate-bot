package user

import (
	"github.com/dpb587/slack-delegate-bot/logic/delegate"
	"github.com/dpb587/slack-delegate-bot/delegatebot/message"
)

type Delegator struct {
	ID string
}

var _ delegate.Delegator = &Delegator{}

func (i Delegator) Delegate(_ message.Message) ([]delegate.Delegate, error) {
	return []delegate.Delegate{delegate.User{ID: i.ID}}, nil
}