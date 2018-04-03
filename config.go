package main

import (
	"fmt"
	"os"

	toml "github.com/BurntSushi/toml"
	hierr "github.com/reconquest/hierr-go"
)

type Config struct {
	Trello struct {
		AppKey string `toml: "appkey"`
		Token  string `toml: "token"`
		Member string `toml: "member"`
		Board  string `toml: "board"`
		List   string `toml: "list"`
	} `toml:"trello"`
}

func loadConfig(path string) (*Config, error) {

	config := &Config{}
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, hierr.Errorf(err, "can't load config %s", path)
	}
	if len(config.Trello.AppKey) == 0 {
		return nil, fmt.Errorf("AppKey is null")
	}
	if len(config.Trello.Token) == 0 {
		return nil, fmt.Errorf("Token is null")
	}
	if len(config.Trello.Member) == 0 {
		return nil, fmt.Errorf("Member is null")
	}
	if len(config.Trello.Board) == 0 {
		return nil, fmt.Errorf("Board is null")
	}
	if len(config.Trello.List) == 0 {
		return nil, fmt.Errorf("List is null")
	}

	return config, nil
}

func saveConfig(path string, config *Config) error {

	outputFile, err := os.OpenFile(
		path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, OutputFileMode,
	)
	if err != nil {
		return hierr.Errorf(err, "can't open output file: %s", path)
	}

	err = toml.NewEncoder(outputFile).Encode(config)
	if err != nil {
		return hierr.Errorf(err, "can't save config %s", path)
	}

	err = outputFile.Close()
	if err != nil {
		return hierr.Errorf(err, "can't close output file: %s", path)
	}

	return nil
}
