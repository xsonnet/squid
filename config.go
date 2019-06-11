package squid

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	// 文件路径
	File string
	// 配置内容
	Data map[string]string
}

func InitConfig(file string) Config {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	var result map[string]string
	err = json.Unmarshal(content, &result)
	if err != nil {
		log.Fatal(err)
	}
	return Config{Data: result}
}

func (config Config) Read(name string) string {
	return config.Data[name]
}

func (config Config) Write(name string, value string) error {
	content, err := ioutil.ReadFile(config.File)
	result := make(map[string]string)
	err = json.Unmarshal(content, &result)
	result[name] = value
	data, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	return ioutil.WriteFile(config.File, data, 0777)
}