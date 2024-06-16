package common

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/samber/lo"
)

type ConfEnvSetting struct {
	YamlFilePath []string
	EnvFilePath  string
}

type Config struct {
	LogLevel   string `config:"logLevel"`
	ServerPort string `config:"serverPort"`

	YellowCardCredentials struct {
		ApiKey    string `config:"apiKey"`
		SecretKey string `config:"secretKey"`
		BaseUrl   string `config:"baseUrl"`
	}

	JWTCredentials struct {
		AccessTokenSecret string `config:"accessTokenSecret"`
		AccessTokenClaim  struct {
			Issuer   string `config:"issuer"`
			Audience string `config:"audience"`
		}
		AccessTokenTTL time.Duration `config:"accessTokenTTL"`
	}

	AppCredentials struct {
		BusinessID  string `config:"businessID"`
		UserEmail   string `config:"userEmail"`
		AdminEmails string `config:"adminEmails"`
	}

	RedisAddr string `config:"redisAddr"`

	MongoDB struct {
		DBUri        string `config:"dbUri"`
		DatabaseName string `config:"databaseName"`
	}

	AllowedCorsOrigin []string `config:"allowedCorsOrigin"`

	SmtpCredentials struct {
		BaseUrl       string `config:"baseUrl"`
		ProjectSecret string `config:"projectSecret"`
	}
}

var (
	cfg    = &Config{}
	loaded = false
)
var once sync.Once

func LoadConfiguration(cfgEnvSetting ConfEnvSetting) (*Config, error) {

	if loaded {
		return cfg, nil
	}

	config.WithOptions(func(opt *config.Options) {
		opt.DecoderConfig.TagName = "config"
	})

	config.WithOptions(config.ParseEnv)

	if len(cfgEnvSetting.YamlFilePath) != 0 {
		config.AddDriver(yaml.Driver)
	}
	files := []string{}
	if !lo.IsEmpty(cfgEnvSetting.EnvFilePath) {
		files = append(files, cfgEnvSetting.EnvFilePath)
	}
	files = append(files, cfgEnvSetting.YamlFilePath...)
	err := config.LoadFiles(files...)
	if err != nil {
		return nil, err
	}

	once.Do(func() {
		err := config.Decode(cfg)
		if err != nil {
			log.Panic(err)
		}
		loaded = true
	})

	if cfg.LogLevel == "debug" {
		cfgJson, err := json.MarshalIndent(cfg, "", "	")
		if err == nil {
			fmt.Printf("Cudium-Backend API Config: %s\n", cfgJson)
		}
	}

	return cfg, nil
}

func GetConfig() *Config {
	return cfg
}

func GetTestConfig(setting ConfEnvSetting) error {
	_, err := LoadConfiguration(setting)
	return err
}
