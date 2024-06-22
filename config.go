package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path"
	"runtime"
	"time"
)

type Config struct {
	MqttBroker      string        `env:"MQTT_BROKER, default=tcp://localhost:1883"`
	MqttUsername    string        `env:"MQTT_USERNAME"`
	MqttPassword    string        `env:"MQTT_PASSWORD"`
	CommandsTopic   string        `env:"MQTT_COMMANDS_TOPIC, default=measure_commands"`
	MeasuresTopic   string        `env:"MQTT_MEASURES_TOPIC, default=measurements/#"`
	ClientId        string        `env:"MQTT_CLIENT_ID, default=go_sensor_reporting"`
	Sensors         []string      `env:"SENSORS, default=bmp,dust"`
	ResponseTimeout time.Duration `env:"RESPONSE_TIMEOUT, default=10s"`
	DbConString     string        `env:"DATABASE_URL"`
	DbMigrationPath string        `env:"DB_MIGRATION_PATH, default=migrations"`
}

func ReadConfig(ctx context.Context) (Config, error) {
	var cfg Config

	if err := envconfig.Process(ctx, &cfg); err != nil {
		return cfg, errors.Wrap(err, "error processing config")
	}

	return cfg, nil
}

func LoadDotEnv() error {
	var pathErr *fs.PathError

	if err := godotenv.Load(".env"); errors.As(err, &pathErr) {
		log.Info("starting with no .env file")
	} else if err != nil {
		return err
	}

	return nil
}

func InitBasePath() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
