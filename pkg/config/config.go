package config

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
)

// Config is the global proxy config
var Config struct {
	Handlers []struct {
		Regexp  string `json:"regexp,omitempty"`
		Handler string `json:"handler,omitempty"`
	} `json:"handlers,omitempty"`
}

// Read reads the global proxy config
func Read(log *logrus.Entry) error {
	log.Print("reading config.yml")

	b, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, &Config)
}
