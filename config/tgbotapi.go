package config

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"net/http"
)

var tgbot *tgServer

func TgServer() *tgServer {
	return tgbot
}
func tgbotapiInit(conf *TgOption) {
	if conf == nil {
		return
	}
	tgbot = newTgServer(conf)
}

type TgOption struct {
	Token     string
	HttpsUrl  string
	HttpsPort int
	Timeout   int
	Debug     bool

	BotName     string
	ChannelName string
}

type tgServer struct {
	opt *TgOption
	bot *tgbotapi.BotAPI
}

func newTgServer(opt *TgOption) *tgServer {
	if opt == nil {
		log.Panic("params needed")
		return nil
	} else if opt.Token == "" {
		log.Panic("bot Token needed")
		return nil
	} else if opt.HttpsPort == 0 {
		opt.HttpsPort = 443
	}
	bot, err := tgbotapi.NewBotAPI(opt.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = opt.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &tgServer{opt: opt, bot: bot}
}

func (s *tgServer) RunLongPoll() (*tgbotapi.BotAPI, *tgbotapi.UpdatesChannel) {
	s.bot.RemoveWebhook()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = s.opt.Timeout

	updates, err := s.bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	return s.bot, &updates
}

func (s *tgServer) RunWebhook() (*tgbotapi.BotAPI, *tgbotapi.UpdatesChannel) {
	//_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://www.google.com:8443/"+bot.Token, "cert.pem"))
	resp, err := s.bot.SetWebhook(tgbotapi.NewWebhook(s.opt.HttpsUrl + "/" + s.bot.Token))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("resp: %s, %#v, %#v", resp.Result, resp.Parameters, resp)
	info, err := s.bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	log.Printf("info: %#v", resp)
	updates := s.bot.ListenForWebhook("/" + s.bot.Token)
	//go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		get, _ := http.Get("https://www.baidu.com")
		defer get.Body.Close()
		all, _ := ioutil.ReadAll(get.Body)
		fmt.Fprintln(writer, string(all))
	})
	go http.ListenAndServe(fmt.Sprintf(":%d", s.opt.HttpsPort), nil)

	return s.bot, &updates
}
