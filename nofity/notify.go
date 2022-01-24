package nofity

import (
	"encoding/json"
	"fmt"
	"github.com/kirinlabs/HttpRequest"
	"go-timer-push/config"
	"go-timer-push/logger"
	"net/http"
	"time"
)

type IMsg interface {
	GetMsgTemplate(level string) (map[string]string, error)
	Send(map[string]string) error
}

type Msg struct {
	MsgInfo map[string]int
	IMsg
}

func (b *Msg) SendMessage() {
	var (
		msg map[string]string
		err error
	)
	msg, err = b.IMsg.GetMsgTemplate("levelA")
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	err = b.IMsg.Send(msg)
	if err != nil {
		logger.Logger.Error(err)
	}
	isSendLevelTwo := false
	for _, v := range b.MsgInfo {
		if v >= 5 {
			isSendLevelTwo = true
			break
		}
	}
	if isSendLevelTwo {
		logger.Logger.Info("启动二级通知")
		msg, err = b.IMsg.GetMsgTemplate("levelB")
		err = b.IMsg.Send(msg)
		if err != nil {
			logger.Logger.Error(err)
		}
	}
}

type PushPlusMsg struct {
	Token    string
	Title    string
	Content  string
	Topic    string
	Template string
}

type PushPlus struct {
	ClientName    string
	SendUrl       string
	LevelOneToken string
	LevelTwoToken string
	Msg
}

func (p *PushPlus) GetMsgTemplate(level string) (map[string]string, error) {
	useToken := p.LevelOneToken
	content := ""
	switch level {
	case "levelA":
		useToken = p.LevelOneToken
		content = "请及时处理工单"
	case "levelB":
		useToken = p.LevelTwoToken
		content = "工单超时，请及时处理"
	}
	res := map[string]string{
		"token":    useToken,
		"title":    fmt.Sprintf("您有%d个工单待分配", len(p.Msg.MsgInfo)),
		"content":  content,
		"topic":    "",
		"template": "html",
	}
	return res, nil
}

func (p *PushPlus) Send(msg map[string]string) error {
	var (
		req  *HttpRequest.Request
		resp *HttpRequest.Response
		err  error
	)
	body := map[string]string{
		"token":   msg["token"],
		"title":   msg["title"],
		"content": msg["content"],
		//"topic":    p.Token,
		"template": msg["template"],
	}
	req = HttpRequest.NewRequest().SetTimeout(3 * time.Second).SetHeaders(
		map[string]string{"Content-Type": "application/json"},
	)
	tempMsg, _ := json.Marshal(body)
	logger.Logger.Info(string(tempMsg))
	if resp, err = req.JSON().Post(p.SendUrl, body); err != nil {
		logger.Logger.Error(err)
		return err
	}
	if resp != nil {
		defer resp.Close()
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Logger.Error(err)
		return err
	}
	return nil
}

func NewPushPlusClient() *PushPlus {
	return &PushPlus{
		ClientName:    "PushPlus",
		SendUrl:       config.Cfg.PushPlus.SendUrl,
		LevelTwoToken: config.Cfg.NoticeUsers[1].PushPlusToken,
		LevelOneToken: config.Cfg.NoticeUsers[0].PushPlusToken,
	}
}
