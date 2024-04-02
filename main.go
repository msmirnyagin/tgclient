package main

import (
	"tgclient/config"
	"tgclient/internal/botcode"
	"tgclient/internal/webhook"

	"context"
	"encoding/json"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/updates"
	updhook "github.com/gotd/td/telegram/updates/hook"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	cfg := config.ReadEnv()
	var chatID int64 = 0
	logger, _ := zap.NewDevelopment(zap.IncreaseLevel(zapcore.InfoLevel), zap.AddStacktrace(zapcore.FatalLevel))
	defer func() { _ = logger.Sync() }()

	d := tg.NewUpdateDispatcher()

	gaps := updates.New(updates.Config{
		Handler: d,
		Logger:  logger.Named("gaps"),
	})

	sessionStorage := &telegram.FileSessionStorage{
		Path: filepath.Join("session.json"),
	}

	// Filling client options.
	options := telegram.Options{
		Logger:         logger,         // Passing logger for observability.
		SessionStorage: sessionStorage, // Setting up session sessionStorage to store auth data.
		UpdateHandler:  gaps,           // Setting up handler for updates from server.
		Middlewares: []telegram.Middleware{
			updhook.UpdateHook(gaps.Handle),
			ratelimit.New(rate.Every(time.Second*5), 5),
		},
	}

	// Setup message update handlers.

	d.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
		m, _ := update.Message.(*tg.Message)
		jitem, _ := json.Marshal(m)
		if err := webhook.GetPost(cfg.Url, string(jitem)); err != nil {
			botcode.SendLog(cfg.Bot, chatID, "Err webhook")
			return errors.Wrap(err, "webhook")
		}
		return nil
	})
	d.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		m, _ := update.Message.(*tg.Message)
		jitem, _ := json.Marshal(m)

		if err := webhook.GetPost(cfg.Url, string(jitem)); err != nil {
			botcode.SendLog(cfg.Bot, chatID, "Err webhook")
			return errors.Wrap(err, "webhook")
		}
		return nil
	})
	//flow Auth
	codePrompt := func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
		code, cid, err := botcode.GetCode(cfg.Bot)
		chatID = cid
		if err != nil {
			return "", err
		}
		logger.Info(strings.TrimSpace(code))
		return strings.TrimSpace(code), nil
	}
	flow := auth.NewFlow(
		auth.Constant(cfg.Phone, cfg.Password, auth.CodeAuthenticatorFunc(codePrompt)),
		auth.SendCodeOptions{})

	//create client
	client := telegram.NewClient(cfg.AppId, cfg.AppHash, options)

	return client.Run(ctx, func(ctx context.Context) error {
		// Perform auth if no session is available. ---bot
		status, err := client.Auth().Status(ctx)
		if err != nil {
			botcode.SendLog(cfg.Bot, chatID, "Err status")
			return errors.Wrap(err, "status")
		}
		if !status.Authorized {
			if err := client.Auth().IfNecessary(ctx, flow); err != nil {
				botcode.SendLog(cfg.Bot, chatID, "Err auth")
				return errors.Wrap(err, "auth")
			}
		}

		// Fetch user info.
		user, err := client.Self(ctx)
		if err != nil {
			botcode.SendLog(cfg.Bot, chatID, "Err call self")
			return errors.Wrap(err, "call self")
		}

		return gaps.Run(ctx, client.API(), user.ID, updates.AuthOptions{
			OnStart: func(ctx context.Context) {
				botcode.SendLog(cfg.Bot, chatID, "Gaps started")
				logger.Info("Gaps started")
			},
		})
	})
}
