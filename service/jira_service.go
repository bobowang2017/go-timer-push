package service

import (
	"encoding/json"
	"github.com/andygrunwald/go-jira"
	"go-timer-push/config"
	"go-timer-push/logger"
	"go-timer-push/utils/encrypt"
	"io/ioutil"
)

var jiraClient *jira.Client

type Field struct {
	FieldAsHtml   string
	fieldCssClass string
}

type Issue struct {
	Key      string
	Position int
	Fields   []Field
}

type UnSignedSupport struct {
	Displayed int
	Total     int
	Start     int
	End       int
	IssueHash string
	Issues    []Issue
	//Columns               []map[string]string
	IsUsingDefaultSorting bool
}

type JiraService struct {
}

func NewJiraService() *JiraService {
	return &JiraService{}
}

func (c *JiraService) GetJiraClient() (*jira.Client, error) {
	aesClient := encrypt.NewAesCryptor()
	username, _ := aesClient.AESBase64Decrypt(config.Cfg.JiraCfg.UserName)
	password, _ := aesClient.AESBase64Decrypt(config.Cfg.JiraCfg.Password)
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	var err error
	if jiraClient == nil {
		jiraClient, err = jira.NewClient(tp.Client(), config.Cfg.JiraCfg.HostUrl)
	}
	if err != nil {
		return nil, err
	}
	return jiraClient, nil
}

func (c *JiraService) respToStr(resp *jira.Response) string {
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		logger.Logger.Error(err)
		return ""
	} else {
		return string(body)
	}
}

func (c *JiraService) Refresh() (res *UnSignedSupport, err error) {
	jiraClient, err := c.GetJiraClient()
	res = &UnSignedSupport{}
	if err != nil {
		return nil, err
	}
	req, _ := jiraClient.NewRequest("GET", config.Cfg.JiraCfg.UnSupportUrl, nil)
	resp, err := jiraClient.Do(req, nil)
	respStr := c.respToStr(resp)
	if err != nil {
		logger.Logger.Error(respStr)
		return nil, err
	}
	logger.Logger.Info(respStr)
	if err := json.Unmarshal([]byte(respStr), res); err != nil {
		logger.Logger.Error(err)
	}
	return res, nil
}
