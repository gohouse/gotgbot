package wait

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (w *Wait) Run(msg *tgbotapi.Message) (ct tgbotapi.Chattable, err error) {
	//if msg.Chat.ID == int64(msg.From.ID) { // 私聊
	//	log.Info("收到私聊: ", msg.From.UserName)
	if w.IsWaiting(msg.From.ID) {
		var keepSession bool
		ct = tgbotapi.NewMessage(msg.Chat.ID, "success, 发送 /mycar 即可查看车辆信息了")
		if !keepSession {
			w.Quit(msg.From.ID)
		}
	}
	//}
	return
}
