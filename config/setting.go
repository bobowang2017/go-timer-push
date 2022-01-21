package config

import (
	"go-timer-push/logger"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type SettingConfig struct {
	JiraHostUrl string `yaml:"JiraHostUrl"`
	Url         string `yaml:"Url"`
	UserName    string `yaml:"UserName"`
	Password    string `yaml:"Password"`
}

var Cfg = &SettingConfig{}

func LoadConfig() {
	logger.Logger.Info("重新加载配置")
	var source = "./cfg.yml"
	if f, err := os.Open(source); err != nil {
		log.Fatalf("打开配置文件失败: %v", err)
	} else {
		if err := yaml.NewDecoder(f).Decode(Cfg); err != nil {
			log.Fatalf("反序列化配置文件失败: %v", err)
		}
	}
}
