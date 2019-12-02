package command

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)


type Yaml struct {
	Server map[string]Server
}

// Parsing 解析配置文件
func ParseConf(file string) *Yaml {
	conf := new(Yaml)
	yamlFile, err := ioutil.ReadFile(file)

	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	log.Println("conf", conf)
	return conf
}