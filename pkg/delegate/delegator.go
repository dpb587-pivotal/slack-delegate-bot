package delegate

import "github.com/dpb587/slack-delegate-bot/cmd/delegatebot/message"

//go:generate counterfeiter . Delegator
type Delegator interface {
	Delegate(message.Message) ([]Delegate, error)
}
