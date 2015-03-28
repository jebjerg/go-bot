package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Save(c Config, path string) error {
	if data, err := json.MarshalIndent(c, "", "    "); err != nil {
		return err
	} else {
		return ioutil.WriteFile(path, data, os.ModePerm)
	}
}

// NewConfig returns the (unspecified JSON formatted) config
func NewConfig(config Config, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, config)
	return err
}
