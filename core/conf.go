package core

import (
	"fmt"
	"github.com/kkvbAugust/fiscobcos-go/config"
	"github.com/kkvbAugust/fiscobcos-go/whole"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

// InitConf 读取yaml文件配置
func InitConf(path string) {
	//const ConfigFile = "resources/settings.yaml"
	c := &config.Config{}
	yamlConf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("get yamlConf error: %s", err))
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		//log.Fatalf("config Init Unmarshal: %v", err)
		logrus.Errorf("config Init Unmarshal: %v", err)
		panic("config Init Unmarshal:" + err.Error())
	}
	log.Println("config yamlFile load Init success")
	whole.Config = c
}
