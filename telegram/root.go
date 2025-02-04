package telegram

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/assimon/captcha-bot/util/config"
	ulog "github.com/assimon/captcha-bot/util/log"
	tb "gopkg.in/telebot.v3"
)

var Bot *tb.Bot

// BotStart 机器人启动
func BotStart() {
	setting := tb.Settings{
		Token: func() string {
			if token := os.Getenv("bot_token"); len(token) > 0 {
				return token
			} else {
				return config.TelegramC.BotToken
			}
		}(),
		Updates: 100,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second, AllowedUpdates: []string{
			"message",
			"chat_member",
			"inline_query",
			"callback_query",
		}},
		OnError: func(err error, context tb.Context) {
			ulog.Sugar.Error(err)
		},
	}
	// 代理
	if config.TelegramC.ApiProxy != "" {
		trans := &http.Transport{
			Proxy: func(_ *http.Request) (*url.URL, error) {
				return url.Parse(config.TelegramC.ApiProxy)
			},
		}
		setting.Client = &http.Client{
			Timeout:   time.Minute,
			Transport: trans,
		}
	}
	var err error
	Bot, err = tb.NewBot(setting)
	if err != nil {
		log.Fatal(err)
	}
	RegisterHandle()
	go RunSyncTask()
	Bot.Start()
}

// RegisterHandle 注册处理器
func RegisterHandle() {
	Bot.Handle(PING_CMD, func(c tb.Context) error {
		return c.Send("pong")
	})
	Bot.Handle(START_CMD, StartCaptcha)
	Bot.Handle(tb.OnChatMember, UserJoinGroup)
	Bot.Handle(tb.OnText, OnTextMessage)
	// 广告
	Bot.Handle(ADD_AD, AddAd, isRootMiddleware)
	Bot.Handle(ALL_AD, AllAd, isRootMiddleware)
	Bot.Handle(DEL_AD, DelAd, isRootMiddleware)
}

// isManageMiddleware 管理员中间件
func isManageMiddleware(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		if isManage(c.Chat(), c.Sender().ID) {
			return next(c)
		}
		return c.Respond(&tb.CallbackResponse{
			Text:      "您未拥有管理员权限，请勿点击！",
			ShowAlert: true,
		})
	}
}

// isRootMiddleware 超管中间件
func isRootMiddleware(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		if !c.Message().Private() || !isRoot(c.Sender().ID) {
			return nil
		}
		return next(c)
	}
}

// isManage 判断是否为管理员
func isManage(chat *tb.Chat, userId int64) bool {
	adminList, err := Bot.AdminsOf(chat)
	if err != nil {
		return false
	}
	for _, member := range adminList {
		if member.User.ID == userId {
			return true
		}
	}
	return false
}

// isRoot 判断是否为超管
func isRoot(userid int64) bool {
	for _, id := range config.TelegramC.GetManageUsers() {
		if userid == id {
			return true
		}
	}
	return false
}
