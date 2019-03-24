package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type(
	Config struct {
		Host string
		Key string
	}
)

var(
	Conf Config
)

func LoadConfig() {
	raw, err := ioutil.ReadFile("config.toml"); if err != nil {
		panic(err)
	}

	_, err = toml.Decode(string(raw), &Conf); if err != nil {
		panic(err)
	}
}