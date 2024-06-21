package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type MqttClient interface {
	Init() error
	Subscribe(handler MessageHandler) error
	Publish(cmd string, sensors []string) error
	Disconnect()
}

type MeasurePayload struct {
	Value  float64 `json:"value"`
	Sensor string  `json:"sensor"`
}

type MessageHandler func(msg MeasurePayload)

func NewMqttClient(cfg Config) MqttClient {
	return &mqttClient{cfg: cfg}
}

const CommandMeasure = "measure"

type mqttClient struct {
	client mqtt.Client
	cfg    Config
}

func (m *mqttClient) Init() error {
	// Create MQTT client options
	opts := mqtt.NewClientOptions().AddBroker(m.cfg.MqttBroker)
	opts.SetClientID(m.cfg.ClientId)
	opts.SetUsername(m.cfg.MqttUsername)
	opts.SetPassword(m.cfg.MqttPassword)

	// Create and connect MQTT client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	m.client = client

	return nil
}

func (m *mqttClient) Subscribe(msgHandler MessageHandler) error {
	// Subscribe to the response topic
	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		var payload MeasurePayload
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.WithField("payload", string(msg.Payload())).Error("error unmarshalling payload")
			return
		}
		msgHandler(payload)
	}
	if token := m.client.Subscribe(m.cfg.MeasuresTopic, 1, messageHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (m *mqttClient) Publish(cmd string, sensors []string) error {
	payload := struct {
		Command string   `json:"cmd"`
		Sensors []string `json:"sensors"`
	}{
		Command: cmd,
		Sensors: sensors,
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	token := m.client.Publish(m.cfg.CommandsTopic, 0, false, bytes)
	token.Wait()

	return token.Error()
}

func (m *mqttClient) Disconnect() {
	m.client.Disconnect(250)
}
