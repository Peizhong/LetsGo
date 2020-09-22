package config

import (
	"os"
)

type ENVConfig struct {
	EnableAPM bool
}

var env = &ENVConfig{}

func init() {
	env.EnableAPM = os.Getenv("APM") != ""
}

func ENV() *ENVConfig {
	return env
}
