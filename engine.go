package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

type Engine interface {
	Work(ctx context.Context) error
}

func NewEngine(cfg Config, mqttClient MqttClient) Engine {
	return &engine{cfg: cfg, mqttClient: mqttClient}
}

type engine struct {
	cfg        Config
	mqttClient MqttClient
}

func (e *engine) Work(ctx context.Context) error {
	responseChan := make(chan MeasurePayload)
	err := e.mqttClient.Subscribe(func(msg MeasurePayload) {
		responseChan <- msg
	})
	if err != nil {
		return err
	}

	if err := e.mqttClient.Publish(CommandMeasure, e.cfg.Sensors); err != nil {
		//if err := e.mqttClient.Publish(CommandMeasure, []string{"bmp", "dust"}); err != nil {
		return err
	}

	msgs, err := e.waitForMessages(responseChan)
	log.WithField("msgs", msgs).Info("received messages")

	return nil
}

func (e *engine) waitForMessages(responseChan chan MeasurePayload) ([]MeasurePayload, error) {
	timeout := time.After(e.cfg.ResponseTimeout)
	var messages []MeasurePayload
	for {
		select {
		case <-timeout:
			return messages, nil
		case msg := <-responseChan:
			messages = append(messages, msg)
		}
	}
}
