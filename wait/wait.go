package wait

import (
	"fmt"
	"sync"
)

type MarkupKeyType int

const (
	MKTNil MarkupKeyType = iota
)

var markupKeyString = map[MarkupKeyType]string{}

var lockMarkupKeyString sync.Mutex

func SetMarkupKeyString(muks map[MarkupKeyType]string) {
	lockMarkupKeyString.Lock()
	defer lockMarkupKeyString.Unlock()
	for k, v := range muks {
		markupKeyString[k] = v
	}
}

func (muk MarkupKeyType) String() string {
	return markupKeyString[muk]
}

type Wait struct {
	//msg tgbotapi.Message
	lock     *sync.RWMutex
	waitList map[int]waitData
}

type waitData struct {
	it   MarkupKeyType
	data interface{}
}

func NewWaitData(it MarkupKeyType, data interface{}) waitData {
	return waitData{it: it, data: data}
}

var waition *Wait

func init() {
	waition = &Wait{lock: &sync.RWMutex{}, waitList: make(map[int]waitData)}
}

func NewWait() *Wait {
	return waition
}

func (w *Wait) Quit(userId int) {
	w.lock.Lock()
	defer w.lock.Unlock()
	//if _, ok := w.waitList[userId]; ok {
	//	w.waitList[userId] = waitData{}
	//}
	delete(w.waitList, userId)
}
func (w *Wait) QuitAndAdd(userId int, it MarkupKeyType, data interface{}) error {
	log.Infof("新增状态: %d - %s - %v", userId, it.String(), data)
	w.Quit(userId)
	return w.Add(userId, it, data)
}
func (w *Wait) Add(userId int, it MarkupKeyType, data interface{}) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	if v, ok := w.waitList[userId]; ok {
		if v.it > MKTNil && v.it != it {
			log.Infof("您当前处于%s状态", v.it.String())
			return fmt.Errorf("您当前处于%s状态, 输入 /cancel 命令退出", v.it.String())
		}
	}
	w.waitList[userId] = NewWaitData(it, data)
	return nil
}

// IsWaiting 会话状态
func (w *Wait) IsWaiting(userId int) bool {
	log.Infof("等待状态: %s", w.waitList[userId].it.String())
	if v, ok := w.waitList[userId]; ok {
		return v.it != MKTNil
	}
	return false
}
func (w *Wait) WaitType(userId int) MarkupKeyType {
	return w.waitList[userId].it
}
func (w *Wait) WaitData(userId int) interface{} {
	return w.waitList[userId].data
}
