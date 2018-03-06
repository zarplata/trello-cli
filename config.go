package main

import (
	"os"

	toml "github.com/BurntSushi/toml"
	hierr "github.com/reconquest/hierr-go"
)

type Config struct {
	Trello struct {
		AppKey string `toml: "appkey" required:"true"`
		Token  string `toml: "token" required:"true"`
		Member string `toml: "member" required:"true"`
		Board  string `toml: "board" required:"true"`
		List   string `toml: "list" required:"true"`
	} `toml:"trello"`
}

func loadConfig(path string) (*Config, error) {

	config := &Config{}
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return config, hierr.Errorf(err, "can't load config %s", path)
	}

	logger.Debugf("successfully load config file %s", path)

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

	logger.Debugf("successfully save config file %s", path)

	return nil

}
