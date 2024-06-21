package main

import (
	"context"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	if err := LoadDotEnv(); err != nil {
		log.WithError(err).Fatal("error loading .env file")
	}

	cfg, err := ReadConfig(ctx)
	if err != nil {
		log.WithError(err).Fatal("error reading config from env")
	}

	mqttClient := NewMqttClient(cfg)
	if err := mqttClient.Init(); err != nil {
		log.WithError(err).Fatal("error initializing mqtt client")
	}

	engine := NewEngine(cfg, mqttClient)
	if err := engine.Work(ctx); err != nil {
		log.WithError(err).Fatal("error running work")
	}
}
