package main

import (
	"fmt"
	"os"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func main() {
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
		s := strings.Split(incoming[0], "/")
		devEUI := s[3]
		sbPubTopic.Reset()
		sbPubTopic.WriteString("lns_imt/")
		sbPubTopic.WriteString(devEUI)
		sbPubTopic.WriteString("/rx")
		// fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
		token := pClient.Publish(sbPubTopic.String(), byte(pQos), false, incoming[1])
		token.Wait()
	}

}
