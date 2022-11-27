package main

import (
	"log"

	"github.com/BurntSushi/toml"
	ms "github.com/PalPalych7/OtusProjectWork/internal/mainstructs"
)

type Config ms.StatSenderConfig

func NewConfig(configFile string) Config {
	var myConf Config
	_, err := toml.DecodeFile(configFile, &myConf)
	if err != nil {
		log.Fatal("err Decode config File=", err)
	}
	return myConf
}
