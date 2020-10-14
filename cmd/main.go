package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gohouse/gotgbot"
	"github.com/gohouse/gotgbot/bot"
	"strings"
)

func main() {
	var token = "1323450187:xxxxxxxxxxxxxxxxxxxxxxxx"
	r := gotgbot.Default(token)

	r.GET("/help", func(ctx *bot.Context) {
		ctx.BotAPI.Send(tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, strings.Join(r.BuildHelpList(bot.String2ModuleType(ctx.Update.Message.Chat.Type)), "\n")))
	}).SetTitle("命令列表查看")

	r.GET("/start", Start)

	r.Run()
	//r.RunWebhook()
}
func Start(ctx *bot.Context) {
	ctx.Send(tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, "欢迎光临,发送 /help 命令解锁更多操作"))
}