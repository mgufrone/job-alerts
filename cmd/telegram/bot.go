package main

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/cmd/telegram/handlers"
	"mgufrone.dev/job-alerts/internal/bootstrap"
	"mgufrone.dev/job-alerts/internal/usecases/job"
	"mgufrone.dev/job-alerts/internal/usecases/notification"
	"mgufrone.dev/job-alerts/pkg/event"
	"mgufrone.dev/job-alerts/pkg/payload"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func startBot(lc fx.Lifecycle, bot *tgbotapi.BotAPI, mgr event.IManager, lgr log.FieldLogger) {
	hook := fx.StartStopHook(func() error {
		bot.Debug = true
		// Create a new UpdateConfig struct with an offset of 0. Offsets are used
		// to make sure Telegram knows we've handled previous values and we don't
		// need them repeated.
		updateConfig := tgbotapi.NewUpdate(0)

		// Tell Telegram we should wait up to 30 seconds on each request for an
		// update. This way we can get information just as quickly as making many
		// frequent requests without having to send nearly as many.
		updateConfig.Timeout = 30

		// Start polling Telegram for updates.
		updates := bot.GetUpdatesChan(updateConfig)
		//bot.Send(tgbotapi.MessageConfig{
		//	Text:     "Sending via notifications",
		//	BaseChat: tgbotapi.BaseChat{ChatID: 5686199640},
		//})

		// Let's go through each update that we're getting from Telegram.
		for update := range updates {
			// Telegram can send many types of updates depending on what your Bot
			// is up to. We only want to look at messages for now, so we can
			// discard any other updates.
			if update.Message == nil {
				continue
			}
			lgr.Debugln("message received", update.Message.Text)
			lgr.Debugln("received message from", update.SentFrom().ID)
			by, _ := json.Marshal(update.Message)
			py := payload.New(by)
			ctx := context.TODO()
			eventType := "chat"
			if update.Message.IsCommand() {
				eventType = "command"
			}
			if err := mgr.Publish(ctx, eventType, py); err != nil {
				lgr.Error("failed to handle event", eventType, py)
			}

			// Now that we know we've gotten a new message, we can construct a
			// reply! We'll take the Chat ID and Text from the incoming message
			// and use it to create a new message.
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "still work in progress, but thank you for the interest")
			//msg.ReplyToMessageID = update.Message.MessageID
			// We'll also say that this message is a reply to the previous message.
			// For any other specifications than Chat ID or Text, you'll need to
			// set fields on the `MessageConfig`.
			//msg.ReplyToMessageID = update.Message.MessageID

			// Okay, we're sending our message off! We don't care about the message
			// we just sent, so we'll discard it.
			//if _, err := bot.Send(msg); err != nil {
			//	// Note that panics are a bad way to handle errors. Telegram can
			//	// have service outages or network errors, you should retry sending
			//	// messages or more gracefully handle failures.
			//	//panic(err)
			//}
		}
		return nil
	}, func() {
		bot.StopReceivingUpdates()
	})
	lc.Append(hook)
}
func main() {
	godotenv.Load()
	app := fx.New(
		bootstrap.AppModule,
		job.Module,
		notification.Module,
		event.Module,
		fx.Provide(
			func() (*tgbotapi.BotAPI, error) {
				return tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
			},
		),
		handlers.Module,
		fx.Invoke(startBot),
	)
	ctx := context.TODO()
	if err := app.Start(ctx); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	<-app.Done()

}
