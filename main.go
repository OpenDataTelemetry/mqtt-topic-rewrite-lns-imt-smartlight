package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type LnsImtTopicRewrite struct {
	ApplicationName string `json:"applicationName"`
	NodeName        string `json:"nodeName"`
	DevEUI          string `json:"devEUI"`
}

func main() {
	var litr LnsImtTopicRewrite

	id := uuid.New().String()
	var sbMqttClientId strings.Builder
	var sbPubTopic strings.Builder
	sbMqttClientId.WriteString("mqtt-topic-rewrite-lns-imt-")
	sbMqttClientId.WriteString(id)

	sTopic := "application/+/node/+/rx"
	sBroker := "mqtt://networkserver.maua.br:1883"
	sClientId := sbMqttClientId.String()
	sUser := "PUBLIC"
	sPassword := "public"
	sQos := 0

	sOpts := MQTT.NewClientOptions()
	sOpts.AddBroker(sBroker)
	sOpts.SetClientID(sClientId)
	sOpts.SetUsername(sUser)
	sOpts.SetPassword(sPassword)

	pBroker := "mqtt://mqtt.maua.br:1883"
	pClientId := sbMqttClientId.String()
	pUser := "public"
	pPassword := "public"
	pQos := 0

	pOpts := MQTT.NewClientOptions()
	pOpts.AddBroker(pBroker)
	pOpts.SetClientID(pClientId)
	pOpts.SetUsername(pUser)
	pOpts.SetPassword(pPassword)

	c := make(chan [2]string)

	sOpts.SetDefaultPublishHandler(func(sMqttClient MQTT.Client, msg MQTT.Message) {
		c <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	sClient := MQTT.NewClient(sOpts)
	if token := sClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", sBroker)
	}

	pClient := MQTT.NewClient(pOpts)
	if token := pClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", pBroker)
	}

	if token := sClient.Subscribe(sTopic, byte(sQos), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for {
		incoming := <-c
		json.Unmarshal([]byte(incoming[1]), &litr)
		var sdt strings.Builder // string device type
		organizationName := "SmartCampusMaua"
		applicationName := litr.ApplicationName
		deviceId := litr.DevEUI
		deviceName := litr.NodeName
		s := strings.Split(deviceName, "_") // string device type
		sdt.WriteString(s[0])
		sdt.WriteString("s")
		deviceType := sdt.String()

		sbPubTopic.Reset()
		sbPubTopic.WriteString("OpenDataTelemetry/")
		sbPubTopic.WriteString(organizationName)
		sbPubTopic.WriteString("/")
		sbPubTopic.WriteString(applicationName)
		sbPubTopic.WriteString("/")
		sbPubTopic.WriteString(deviceType)
		sbPubTopic.WriteString("/")
		sbPubTopic.WriteString(deviceId)
		sbPubTopic.WriteString("/rx/lns_imt")
		// fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
		token := pClient.Publish(sbPubTopic.String(), byte(pQos), false, incoming[1])
		token.Wait()
	}

}
