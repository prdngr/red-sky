package core

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type config struct {
	Terraform struct {
		ExecutablePath string    `koanf:"executablePath"`
		Initialized    bool      `koanf:"initialized"`
		UpdatedAt      time.Time `koanf:"updatedAt"`
	}
}

const (
	NodVersion    = "0.1.0"
	ConfigFile    = "config.yaml"
	KeysDirectory = "keys"
)

var (
	k      = koanf.New(".")
	parser = yaml.Parser()
	Config config
)

func ReadConfig() {
	configFile := path.Join(GetNodDirectory(), ConfigFile)

	if err := k.Load(file.Provider(configFile), parser); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	if err := k.Unmarshal("", &Config); err != nil {
		log.Fatalf("error parsing config file: %s", err)
	}
}

func WriteConfig() {
	b, _ := k.Marshal(parser)
	fmt.Println(string(b))
}
