package conf

import (
	"bufio"
	"errors"
	"flag"
	"gopkg.in/yaml.v3"
	"net.dalva.GvawSkinSync/logger"
	"os"
	"sync"
)

type Config struct {
	DbHost     string `yaml:"db-host"`
	DbPort     string `yaml:"db-port"`
	DbName     string `yaml:"db-name"`
	DbUser     string `yaml:"db-user"`
	DbPass     string `yaml:"db-pass"`
	DbTimezone string `yaml:"db-timezone"`
	Port       int    `yaml:"port"`
}

const configFile = "./config.yaml"

var (
	cfgInitOnce sync.Once
	Cfg         Config
	LogLevel    string
)

func init() {
	logger.InitializeLoggerOnce()
	InitializeConfigOnce()

}

func InitializeConfigOnce() {
	cfgInitOnce.Do(func() {
		// Process command line parameters
		flag.StringVar(&LogLevel, "l", "debug", "Log Level: trace, debug, info, war, error, fatal, panic")
		flag.Parse()

		err := logger.SetLevel(LogLevel)
		if err != nil {
			logger.Log.Fatal().Err(err).Str("level", LogLevel).Msg("Unknown Log Level")
		}

		_, err = os.Stat(configFile)
		if err != nil {
			logger.Log.Error().Err(err).Msg("No config file. Creating from template.")
			configCreateDefaults()
			logger.Log.Fatal().Msg("Config created, insert correct values then reload App.")
		}

		err = reloadConfig()
		if err != nil {
			logger.Log.Fatal().Err(err).Msg("Fatal")
		}

		logger.Log.Info().Msg("INIT: Config Loaded Successfully")
	})
}

func configCreateDefaults() {
	cfg := Config{
		DbHost:     "http://localhost",
		DbPort:     "5432",
		DbName:     "DbName",
		DbUser:     "DbUser",
		DbPass:     "",
		DbTimezone: "Asia/Jakarta",
		Port:       24003,
	}

	f, err := os.Create(configFile)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("error creating new config file. aborting")
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	encoder := yaml.NewEncoder(w)
	err = encoder.Encode(cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("error creating new config file. aborting")
	}
	err = encoder.Close()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("error creating new config file. aborting")
	}
	err = w.Flush()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("error creating new config file. aborting")
	}
}

func reloadConfig() error {
	f, err := os.Open(configFile)
	if err != nil {
		return errors.New("error loading config")
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Cfg)
	if err != nil {
		return errors.New("error decoding config")
	}

	return nil
}
