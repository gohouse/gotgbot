package gotgbot

import (
	"encoding/json"
	"fmt"
	"github.com/gohouse/gotgbot/config"
	"io/ioutil"
	"os"
)

var option config.ConfigOption
func buildConfigWithFile(files string) {
	//var conf config.ConfigOption
	var file = "config.json"
	if files != "" {
		file = files
	}
	switch len(os.Args) {
	case 2:
		file = os.Args[1]
	}
	stat, err2 := PathExists(file)
	if err2!=nil {
		panic(err2.Error())
	}
	if !stat {
		fmt.Println("配置文件缺失, config.json")
		os.Exit(1)
	}

	readFile, err := ioutil.ReadFile(file)
	if err!=nil {
		panic(err.Error())
	}
	//logrus.Infof("%s", readFile)
	err = json.Unmarshal(readFile, &option)
	if err!=nil {
		panic(err.Error())
	}
	//logrus.Infof("%#v", option)
}
