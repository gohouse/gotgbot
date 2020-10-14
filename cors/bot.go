package cors

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gohouse/gotgbot/bot"
	"time"
)

func MessageTimeout(ctx *bot.Context) {
	TimeoutWithSeconds(15, ctx.Message)(ctx)
}

func TimeoutWithSeconds(s int, msg *botapi.Message) bot.HandleFunc {
	return func(ctx *bot.Context)  {
		go func() {
			select {
			case <-time.After(time.Duration(s)*time.Second):
				if ctx.SendMessageId == 0 {
					log.Infof("%ds后关闭操作触发失败, 因为已经有更早的触发了", s)
					return
				}
				log.Infof("%ds后关闭操作触发", s)
				// 删除回复消息
				ctx.Send(botapi.NewDeleteMessage(msg.Chat.ID, ctx.SendMessageId))
				// 删除命令消息
				ctx.Send(botapi.NewDeleteMessage(msg.Chat.ID, ctx.Message.MessageID))
				ctx.SendMessageId = 0
			}
		}()
	}
}
