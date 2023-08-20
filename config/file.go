package config

import (
	"encoding/json"
	"io"
	"os"
)

type PolicyConfig struct {
	Service string            `json:"service"`
	Cluster string            `json:"cluster"`
	Tenant  string            `json:"tenant"`
	Backend []BackendConfig   `json:"backend"`
	Mode    int               `json:"mode"`
	Url     string            `json:"url"`
	Query   map[string]string `json:"query"`
}

type BackendConfig struct {
	Name   string `json:"name"`
	Addr   string `json:"addr"`
	Weight int    `json:"weight"`
}

var globalConfig []PolicyConfig

func Init() {
	loadFileConfig(os.Getenv("CONFIG_FILE_PATH"))
}

func loadFileConfig(conf string) []PolicyConfig {
	f, err := os.Open(conf)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	var configs []PolicyConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		panic(err)
		return nil
	}
	if len(configs) > 0 {
		globalConfig = configs
	}
	return configs
}

func GetConfigs() []PolicyConfig {
	return globalConfig
}
