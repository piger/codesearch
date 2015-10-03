package server

import (
	"encoding/json"
	"io/ioutil"
)

type ServerConfig struct {
	Repositories []Repository
}

type Repository struct {
	Name   string
	Dir    string
	Remote string
}

func ReadConfig(filename string) (*ServerConfig, error) {
	cfgdata, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config ServerConfig
	if err := json.Unmarshal(cfgdata, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
