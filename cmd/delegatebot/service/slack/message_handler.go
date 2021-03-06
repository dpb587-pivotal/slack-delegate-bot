package slack

import (
	"fmt"
	"regexp"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/handler"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/message"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

var reChannelMention = regexp.MustCompile(`<#([^|]+)|([^>]+)>`)

type MessageHandler struct {
	rtm             *slack.RTM
	delegateHandler handler.Handler
}

func NewMessageHandler(rtm *slack.RTM, delegateHandler handler.Handler) *MessageHandler {
	return &MessageHandler{
		rtm:             rtm,
		delegateHandler: delegateHandler,
	}
}

func (h *MessageHandler) GetResponse(request message.Message, ev *slack.MessageEvent) (*slack.OutgoingMessage, error) {
	response, err := h.delegateHandler.Execute(&request)
	if err != nil {
		return nil, errors.Wrap(err, "executing handler")
	}

	var msg string

	if len(response.Delegates) > 0 {
		msg = delegates.Join(response.Delegates, " ")

		if request.OriginType == message.ChannelOriginType {
			msg = fmt.Sprintf("^ %s", msg)
		}
	} else if response.EmptyMessage != "" {
		msg = response.EmptyMessage
	}

	if msg == "" {
		return nil, nil
	}

	outgoing := h.rtm.NewOutgoingMessage(msg, ev.Msg.Channel)

	if request.OriginType == message.ChannelOriginType {
		outgoing.ThreadTimestamp = ev.Msg.Timestamp
	}

	return outgoing, nil
}
