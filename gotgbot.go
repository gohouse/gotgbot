package gotgbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gohouse/gotgbot/bot"
	"github.com/gohouse/gotgbot/config"
	"github.com/gohouse/gotgbot/router"
	"time"
)

var log = config.Logger()

type GoTgBot struct {
	*router.RouteGroup
}

func NewGoTgBot(conf *config.ConfigOption) *GoTgBot {
	config.Init(conf)
	return &GoTgBot{RouteGroup: router.NewRouteGroup()}
}

func NewGoTgBotWithFile(file string) *GoTgBot {
	buildConfigWithFile(file)
	config.Init(&option)
	return &GoTgBot{RouteGroup: router.NewRouteGroup()}
}

func Default(token string) (r *GoTgBot) {
	// 运行
	return NewGoTgBot(&config.ConfigOption{TgOption: config.TgOption{Token: token}})
}

func (gb *GoTgBot) Run() {
	// 加载
	gb.run(config.TgServer().RunLongPoll())
}

func (gb *GoTgBot) RunWebhook() {
	// 加载
	gb.run(config.TgServer().RunWebhook())
}

func (gb *GoTgBot) AnswerCallbackFail(botapi *tgbotapi.BotAPI, callbackId, msgText string) {
	botapi.AnswerCallbackQuery(tgbotapi.NewCallback(callbackId, fmt.Sprintf("❌%s", msgText)))
}

func (gb *GoTgBot) run(botapi *tgbotapi.BotAPI, updates *tgbotapi.UpdatesChannel) {
	go func(ba *tgbotapi.BotAPI) {
		for _, v := range gb.HandBot {
			v(bot.NewContext(botapi, nil, nil))
		}
	}(botapi)
	for ups := range *updates {
		up := ups
		go func(update tgbotapi.Update) {
			//panic回收保护
			defer func() {
				if err := recover(); err != nil {
					log.Error("panic recover:", err)
				}
			}()
			// debug
			gb.debugCallback(&update)
			// 启动 callback
			if update.CallbackQuery != nil {
				log.Info("CallbackQuery act")
				ctx := bot.NewContext(botapi, &update, gb.Middleware)
				ctx.Start()
				for _, v := range gb.CallbackQuery {
					v(ctx)
				}
				return
			}

			if update.Message != nil {
				// 如果存在给定的命令列表中, 则启动 command 流程, 否则走入 wait 流程
				if update.Message.IsCommand() {
					mot := bot.String2ModuleType(update.Message.Chat.Type)
					cmd := update.Message.Command()
					m, h := gb.ExtractRoute(cmd, mot)
					if len(h) > 0 {
						ctx := bot.NewContext(botapi, &update, append(m, h...))
						ctx.Start()
						return
					}
				}
				// 启动 wait
				log.Info("wait act")
				ctx := bot.NewContext(botapi, &update, gb.Middleware)
				ctx.Start()
				for _, v := range gb.Wait {
					v(ctx)
				}
			}
		}(up)
	}
}

func (gb *GoTgBot) debugCallback(update *tgbotapi.Update) {
	var filename = time.Now().Format("0102-150405")
	var msgFilename = fmt.Sprintf("%s-msg", filename)
	if update.CallbackQuery != nil {
		msgFilename = fmt.Sprintf("%s-callback-msg", filename)
		newlog(filename, update.CallbackQuery.Message)
	} else if update.Message != nil {
		newlog(filename, update.Message)
	}
	newlog(msgFilename, update)
}
