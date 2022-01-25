package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"go-timer-push/config"
	"go-timer-push/logger"
	"go-timer-push/nofity"
	"go-timer-push/service"
	"net/http"
)

var UnSupportCache = map[string]int{}

func IntervalRefresh() {
	var (
		jiraSvc        *service.JiraService
		unSupportIssue *service.UnSignedSupport
		err            error
	)
	defer func() {
		if panicErr := recover(); panicErr != nil {
			logger.Logger.Error(panicErr)
		}
	}()
	logger.Logger.Info("开始刷新请求数据")
	jiraSvc = service.NewJiraService()
	unSupportIssue, err = jiraSvc.Refresh()
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	//所有工单已分配，直接清空缓存
	if unSupportIssue.Total == 0 {
		UnSupportCache = map[string]int{}
		return
	}
	unSupportIssueMap := map[string]int{}
	for _, issue := range unSupportIssue.Issues {
		unSupportIssueMap[issue.Key] = 1
		if _, ok := UnSupportCache[issue.Key]; ok {
			UnSupportCache[issue.Key] += 1
		} else {
			UnSupportCache[issue.Key] = 1
		}
	}
	for k, _ := range UnSupportCache {
		if _, ok := unSupportIssueMap[k]; !ok {
			delete(UnSupportCache, k)
		}
	}
	msg := make([]string, 1)
	for k, v := range UnSupportCache {
		msg = append(msg, fmt.Sprintf("工单%s未响应，超时时间%d分钟", k, v*2))
		logger.Logger.Infof("工单%s未响应，超时时间%d分钟", k, v*2)
	}
	pushPlusClient := nofity.NewPushPlusClient()
	msgClient := nofity.Msg{
		IMsg:    pushPlusClient,
		MsgInfo: UnSupportCache,
	}
	pushPlusClient.Msg = msgClient
	msgClient.SendMessage()
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("I am ok"))
}

func main() {
	config.LoadConfig()
	logger.Setup()
	logger.Logger.Info("程序启动")
	IntervalRefresh()
	cronTimer := cron.New()
	_, _ = cronTimer.AddFunc("*/2 * * * *", IntervalRefresh)
	cronTimer.Start()
	//sigChan := make(chan os.Signal, 1)
	//signal.Notify(sigChan, syscall.SIGINT, os.Interrupt, syscall.SIGTERM)
	//<-sigChan
	http.HandleFunc("/", HealthHandler)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		logger.Logger.Error(err)
		logger.Logger.Panic("系统服务启动异常")
	}
}
