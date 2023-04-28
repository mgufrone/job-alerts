package handlers

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mgufrone.dev/job-alerts/pkg/payload"
)

type Help struct {
	bot *tgbotapi.BotAPI
}

func NewHelp(bot *tgbotapi.BotAPI) *Help {
	return &Help{bot: bot}
}

func (h *Help) Name() string {
	return "help"
}

func (h *Help) Run(ctx context.Context, py payload.Payload) error {
	var msg tgbotapi.Message
	err := py.As("", &msg)
	if err != nil {
		return err
	}
	reply := tgbotapi.NewMessage(msg.Chat.ID, `Here's what you can do with the bot.
/subscriptions - get your current subscriptions to the job alerts
/newsubscription - create new subscription within the criteria you set
/help - really? I need to explain this command?

For direct chat, you can do these
list recent jobs - it will show you recent jobs
list 
`)
	if msg.MessageID > 0 {
		reply.ReplyToMessageID = msg.MessageID
	}

	_, err = h.bot.Send(reply)
	return err
}

func (h *Help) When(py payload.Payload) bool {
	var msg tgbotapi.Message
	err := py.As("", &msg)
	return err == nil && ((msg.IsCommand() && msg.Command() == "help") || msg.Text == "help")
}

func (h *Help) SubscribeAt() []string {
	return []string{
		"chat",
		"command",
	}
}
