package openapi

import (
	"github.com/tencent-connect/botgo/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// Type token 类型
type Type string

const (
	TypeBot    Type = "Bot"
	TypeNormal Type = "Bearer"
)

// Token 用于调用接口的 token 结构
type Token struct {
	AppID       uint64
	AccessToken string
	Type        Type
}

// LoadFromConfig 从配置中读取 appid 和 token
func (t *Token) LoadFromConfig(file string) error {
	var conf struct {
		AppID uint64 `yaml:"appid"`
		Token string `yaml:"token"`
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Errorf("read token from file failed, err: %v", err)
		return err
	}
	if err = yaml.Unmarshal(content, &conf); err != nil {
		log.Errorf("parse config failed, err: %v", err)
		return err
	}
	t.AppID = conf.AppID
	t.AccessToken = conf.Token
	return nil
}
