package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math"
	"time"
)

const abortIndex int8 = math.MaxInt8 / 2

type HandleFunc func(ctx *Context)
type HandleBot func(botapi *tgbotapi.BotAPI)
type Context struct {
	BotAPI        *tgbotapi.BotAPI
	Update        *tgbotapi.Update
	Param         string
	SendMessageId int
	index         int8
	handlers      []HandleFunc
}

func NewContext(botAPI *tgbotapi.BotAPI, update *tgbotapi.Update, handlers []HandleFunc) *Context {
	return &Context{BotAPI: botAPI, Update: update, handlers: handlers}
}
func (c *Context) AnswerCallbackQuerySuccess(id, text string) (tgbotapi.APIResponse, error) {
	// ⚙️⚒🛠⚒🔗🖊🖋❌✅
	return c.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(id, fmt.Sprintf("✅%s", text)))
}
func (c *Context) AnswerCallbackQueryFail(id, text string) (tgbotapi.APIResponse, error) {
	return c.BotAPI.AnswerCallbackQuery(tgbotapi.NewCallback(id, fmt.Sprintf("❌%s", text)))
}
func (c *Context) Send(ct tgbotapi.Chattable) (send tgbotapi.Message, err error) {
	send, err = c.BotAPI.Send(ct)
	if err != nil {
		log.Error("发送消息错误:", err.Error())
		return
	}
	c.SendMessageId = send.MessageID
	return
}
func (c *Context) SendWithExpire(ct tgbotapi.Chattable, expire int) (send tgbotapi.Message, err error) {
	send, err = c.BotAPI.Send(ct)
	if err != nil {
		log.Error("发送消息错误:", err.Error())
		return
	}
	go func() {
		time.Sleep(time.Duration(expire) * time.Second)
		c.BotAPI.Send(tgbotapi.NewDeleteMessage(send.Chat.ID, send.MessageID))
	}()
	return
}

/************************************/
/*********** FLOW CONTROL ***********/
/************************************/

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		if c.IsAborted() {
			break
		}
		(c.handlers)[c.index](c)
		c.index++
	}
}
func (c *Context) Start() {
	c.index = 0
	if len(c.handlers) > 0 {
		(c.handlers)[c.index](c)
		c.Next()
	}
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = abortIndex
}
