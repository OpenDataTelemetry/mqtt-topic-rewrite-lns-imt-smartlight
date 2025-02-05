package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func connLostHandler(c MQTT.Client, err error) {
	fmt.Printf("Connection lost, reason: %v\n", err)
	os.Exit(1)
}

func main() {
	// MQTT.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	MQTT.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	// MQTT.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	// MQTT.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	id := uuid.New().String()
	var sbMqttSubClientId strings.Builder
	var sbMqttPubClientId strings.Builder
	var sbPubTopic strings.Builder
	sbMqttSubClientId.WriteString("mqtt-topic-rewrite-lns-imt-")
	sbMqttSubClientId.WriteString(id)
	sbMqttPubClientId.WriteString("mqtt-topic-rewrite-lns-imt-")
	sbMqttPubClientId.WriteString(id)

	mqttSubBroker := "mqtt://networkserver.maua.br:1883"
	mqttSubClientId := sbMqttSubClientId.String()
	mqttSubUser := "PUBLIC"
	mqttSubPassword := "public"
	mqttSubQos := 0

	mqttSubOpts := MQTT.NewClientOptions()
	mqttSubOpts.AddBroker(mqttSubBroker)
	mqttSubOpts.SetClientID(mqttSubClientId)
	mqttSubOpts.SetUsername(mqttSubUser)
	mqttSubOpts.SetPassword(mqttSubPassword)
	mqttSubOpts.SetConnectionLostHandler(connLostHandler)

	mqttSubTopics := map[string]byte{
		"application/1/node/+/rx":  byte(mqttSubQos),
		"application/6/node/+/rx":  byte(mqttSubQos),
		"application/9/node/+/rx":  byte(mqttSubQos),
		"application/13/node/+/rx": byte(mqttSubQos),
		"application/18/node/+/rx": byte(mqttSubQos),
		"application/19/node/+/rx": byte(mqttSubQos),
		"application/20/node/+/rx": byte(mqttSubQos),
	}

	mqttPubBroker := "mqtt://mqtt.maua.br:1883"
	mqttPubClientId := sbMqttPubClientId.String()
	mqttPubUser := "public"
	mqttPubPassword := "public"
	mqttPubQos := 0

	mqttPubOpts := MQTT.NewClientOptions()
	mqttPubOpts.AddBroker(mqttPubBroker)
	mqttPubOpts.SetClientID(mqttPubClientId)
	mqttPubOpts.SetUsername(mqttPubUser)
	mqttPubOpts.SetPassword(mqttPubPassword)

	c := make(chan [2]string)

	mqttSubOpts.SetDefaultPublishHandler(func(mqttSubClient MQTT.Client, msg MQTT.Message) {
		c <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	mqttSubClient := MQTT.NewClient(mqttSubOpts)
	if token := mqttSubClient.Connect(); token.Wait() && token.Error() != nil {
		// fmt.Printf("Error %s\n", token.Error())
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", mqttSubBroker)
	}

	pClient := MQTT.NewClient(mqttPubOpts)
	if token := pClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", mqttPubBroker)
	}

	// SubscribeMultiple(filters map[string]byte, callback MessageHandler) Token
	if token := mqttSubClient.SubscribeMultiple(mqttSubTopics, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for {
		incoming := <-c
		s := strings.Split(incoming[0], "/")
		goTopic := false
		var measurement string
		switch s[1] {
		case "1":
			if s[3] == "0004a30b001e0b53" || s[3] == "0004a30b01000200" || s[3] == "0004a30b0028ddd9" {
				measurement = "SoilMoisture3DepthLevels"
				goTopic = true
			} else if s[3] == "0004a30b001a1d6f" {
				measurement = "VibrationAverage"
				goTopic = true
			} else if s[3] == "0004a30b0023580e" {
				measurement = "Temperature8Point"
				goTopic = true
			} else {
				break
			}
		case "6":
			if s[3] == "0004a30b00e94314" {
				measurement = "Sprinkler"
				goTopic = true
			} else {
				measurement = "SmartLight"
				goTopic = true
			}
		case "9":
			measurement = "EnergyMeter"
			goTopic = true
		case "13":
			measurement = "WeatherStation"
			goTopic = true
		case "18":
			measurement = "WaterTankLevel"
			goTopic = true
		case "19":
			measurement = "GaugePressure"
			goTopic = true
		case "20":
			measurement = "Hydrometer"
			goTopic = true
		}
		deviceId := s[3]

		if goTopic == true {
			sbPubTopic.Reset()
			sbPubTopic.WriteString("OpenDataTelemetry/IMT/LNS/")
			sbPubTopic.WriteString(measurement)
			sbPubTopic.WriteString("/")
			sbPubTopic.WriteString(deviceId)
			sbPubTopic.WriteString("/up/imt")
			fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
			token := pClient.Publish(sbPubTopic.String(), byte(mqttPubQos), false, incoming[1])
			token.Wait()
		}

	}
}
