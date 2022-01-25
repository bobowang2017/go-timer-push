package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Server struct {
	LogSavePath     string `yaml:"LogSavePath"`
	LogMaxAge       string `yaml:"LogMaxAge"`
	LogRotationTime string `yaml:"LogRotationTime"`
	LogLevel        string `yaml:"LogLevel"`
}

type JiraCfg struct {
	HostUrl        string `yaml:"HostUrl"`
	UserName       string `yaml:"UserName"`
	Password       string `yaml:"Password"`
	UnSupportUrl   string `yaml:"UnSupportUrl"`
	MyUnSupportUrl string `yaml:"MyUnSupportUrl"`
}

type NoticeUser struct {
	Username      string `yaml:"Username"`
	Password      string `yaml:"Password"`
	PushPlusToken string `yaml:"PushPlusToken"`
}

type SettingConfig struct {
	Server   Server  `yaml:"Server"`
	JiraCfg  JiraCfg `yaml:"JiraCfg"`
	PushPlus struct {
		SendUrl string `yaml:"SendUrl"`
	} `yaml:"PushPlus"`
	LevelOneNoticeUsers []NoticeUser `yaml:"LevelOneNoticeUsers"`
	LevelTwoNoticeUsers []NoticeUser `yaml:"LevelTwoNoticeUsers"`
}

var Cfg = &SettingConfig{}

func LoadConfig() {
	var source = "./cfg-dev.yml"
	if f, err := os.Open(source); err != nil {
		log.Fatalf("打开配置文件失败: %v", err)
	} else {
		if err := yaml.NewDecoder(f).Decode(Cfg); err != nil {
			log.Fatalf("反序列化配置文件失败: %v", err)
		}
	}
}

func LoadTestConfig() {
	var source = "../cfg-dev.yml"
	if f, err := os.Open(source); err != nil {
		log.Fatalf("打开配置文件失败: %v", err)
	} else {
		if err := yaml.NewDecoder(f).Decode(Cfg); err != nil {
			log.Fatalf("反序列化配置文件失败: %v", err)
		}
	}
}
