package config

import (
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

//func logInit(l *logrus.Logger) {
//	if l != nil {
//		*log = *l
//	}
//}
func Logger() *logrus.Logger {
	if log == nil {
		log = logrus.New()
		//file, err := os.OpenFile(fmt.Sprintf("log/tg_%s.log", time.Now().Format("20060102")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		//if err!=nil {
		//	panic(err.Error())
		//}
		//log.Out = file
	}
	return log
}
