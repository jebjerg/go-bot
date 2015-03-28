package config

import (
	"testing"
)

type twitter_conf struct {
	Channels       []string `json:"channels"`
	ConsumerKey    string   `json:"consumer_key"`
	ConsumerSecret string   `json:"consumer_secret"`
}

func TestConfig(t *testing.T) {
	conf := &twitter_conf{}
	if err := NewConfig(conf, "fixtures/twitter.json"); err != nil {
		t.Errorf("unable to load twitter.json config: %v", err)
	}
	if err := Save(conf, "/tmp/null.json"); err != nil {
		t.Errorf("unable to save config: %v", err)
	}
}
