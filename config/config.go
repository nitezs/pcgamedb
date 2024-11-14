package config

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type config struct {
	LogLevel           string    `env:"LOG_LEVEL" json:"log_level"`
	Server             server    `json:"server"`
	Database           database  `json:"database"`
	Redis              redis     `json:"redis"`
	OnlineFix          onlinefix `json:"online_fix"`
	Twitch             twitch    `json:"twitch"`
	Webhooks           webhooks  `json:"webhooks"`
	DatabaseAvaliable  bool
	OnlineFixAvaliable bool
	MegaAvaliable      bool
	RedisAvaliable     bool
}

type webhooks struct {
	CrawlTask []string `env:"WEBHOOKS_ERROR_TASK" json:"crawl_task"`
}

type server struct {
	Port      string `env:"SERVER_PORT" json:"port"`
	SecretKey string `env:"SERVER_SECRET_KEY" json:"secret_key"`
	AutoCrawl bool   `env:"SERVER_AUTO_CRAWL" json:"auto_crawl"`
}

type database struct {
	Host     string `env:"DATABASE_HOST" json:"host"`
	Port     int    `env:"DATABASE_PORT" json:"port"`
	User     string `env:"DATABASE_USER" json:"user"`
	Password string `env:"DATABASE_PASSWORD" json:"password"`
	Database string `env:"DATABASE_NAME" json:"database"`
}

type twitch struct {
	ClientID     string `env:"TWITCH_CLIENT_ID" json:"client_id"`
	ClientSecret string `env:"TWITCH_CLIENT_SECRET" json:"client_secret"`
}

type redis struct {
	Host     string `env:"REDIS_HOST" json:"host"`
	Port     int    `env:"REDIS_PORT" json:"port"`
	Password string `env:"REDIS_PASSWORD" json:"password"`
	DBIndex  int    `env:"REDIS_DB" json:"db_index"`
}

type onlinefix struct {
	User     string `env:"ONLINEFIX_USER" json:"user"`
	Password string `env:"ONLINEFIX_PASSWORD" json:"password"`
}

type runtimeConfig struct {
	ServerStartTime time.Time
}

var Config config
var Runtime runtimeConfig

func init() {
	Config = config{
		LogLevel: "info",
		Database: database{
			Port:     27017,
			User:     "root",
			Password: "password",
		},
		MegaAvaliable: TestMega(),
	}
	if _, err := os.Stat("config.json"); err == nil {
		configData, err := os.ReadFile("config.json")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(configData, &Config)
		if err != nil {
			panic(err)
		}
	}
	loadEnvVariables(&Config)
	Config.OnlineFixAvaliable = Config.OnlineFix.User != "" && Config.OnlineFix.Password != ""
	Config.RedisAvaliable = Config.Redis.Host != ""
	Config.DatabaseAvaliable = Config.Database.Database != "" && Config.Database.Host != ""
}

func loadEnvVariables(cfg interface{}) {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" || envTag == "-" {
			if field.Type.Kind() == reflect.Struct {
				loadEnvVariables(v.Field(i).Addr().Interface())
			}
			continue
		}
		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}
		switch field.Type.Kind() {
		case reflect.String:
			v.Field(i).SetString(envValue)
		case reflect.Int:
			if value, err := strconv.Atoi(envValue); err == nil {
				v.Field(i).SetInt(int64(value))
			}
		case reflect.Bool:
			if value, err := strconv.ParseBool(envValue); err == nil {
				v.Field(i).SetBool(value)
			}
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.String {
				envValueSlice := strings.Split(envValue, ",")
				v.Field(i).Set(reflect.ValueOf(envValueSlice))
			}
		}
	}
}

func TestMega() bool {
	cmd := exec.Command("mega-get", "--help")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return err == nil
}
