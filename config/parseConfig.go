// 读取配置文件信息
package config

import (
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	Config conf
	once   sync.Once
)

type Account struct {
	MachineId string `yaml:"machineid"`
	Password  string `yaml:"password"`
}

type conf struct {
	DetectUrl string    `yaml:"detecturl"`
	Accounts  []Account `yaml:"accounts"`
	Goods     []string  `yaml:"goods"`
}

func (c *conf) getConf() {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func init() {
	once.Do(readConfig)
}

func readConfig() {
	Config.getConf()
}
