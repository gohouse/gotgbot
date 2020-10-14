package gotgbot

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gohouse/gorose/v2"
	"github.com/gohouse/gotgbot/config"
	"github.com/gohouse/t"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func GetLargePhotoFromResponse(pic *[]tgbotapi.PhotoSize) (photo tgbotapi.PhotoSize) {
	var width int
	for _, v := range *pic {
		width = t.Max(v.Width, width).Int()
		if v.Width == width {
			photo = v
		}
	}
	return
}

func PathExists(path string) (bool, error) {

	_, err := os.Stat(path)

	if err == nil {

		return true, nil

	}

	if os.IsNotExist(err) {

		return false, nil

	}

	return false, err

}
func SimpleNewlog(msg interface{}) {
	var filename = time.Now().Format("0102-150405")
	var msgFilename = fmt.Sprintf("%s-msg-%v", filename, time.Now().Nanosecond())
	newlog(msgFilename, msg)
}
func newlog(file string, msg interface{}) {
	openFile, _ := os.OpenFile(fmt.Sprintf("log/tg_%s.log", file), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	marshal, _ := json.MarshalIndent(msg, "", "    ")
	fmt.Fprintln(openFile, string(marshal))
}

func Redis() *redis.Client {
	return config.Redis()
}

func Logger() *logrus.Logger {
	return config.Logger()
}

func DB() gorose.IOrm {
	return config.DB()
}
