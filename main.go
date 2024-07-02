package main

import (
	"fmt"
	"os"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

// type LnsImtTopicRewrite struct {
// 	ApplicationName string `json:"applicationName"`
// 	NodeName        string `json:"nodeName"`
// 	DevEUI          string `json:"devEUI"`
// }

func main() {
	// var litr LnsImtTopicRewrite

	id := uuid.New().String()
	var sbMqttClientId strings.Builder
	var sbPubTopic strings.Builder
	sbMqttClientId.WriteString("mqtt-topic-rewrite-lns-imt-smartlight-")
	sbMqttClientId.WriteString(id)

	sTopic := "application/6/node/+/rx"
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

	// m := make(map[string]string)
	// m["0004a30b00e9442a"] = "SmartLight"
	// m["0004a30b00e9485c"] = "SmartLight"
	// m["0004a30b00e9b7ce"] = "SmartLight"
	// m["0004a30b00e942f2"] = "SmartLight"
	// m["0004a30b00e97b18"] = "SmartLight"
	// m["0004a30b00e95edc"] = "SmartLight"
	// m["0004a30b00e97b62"] = "SmartLight"
	// m["0004a30b00e932d0"] = "SmartLight"
	// m["0004a30b00e9833b"] = "SmartLight"
	// m["0004a30b00e9b664"] = "SmartLight"
	// m["0004a30b00e9b90b"] = "SmartLight"
	// m["0004a30b00e9d3a2"] = "SmartLight"
	// m["0004a30b00e9651e"] = "SmartLight"
	// m["0004a30b00e940d5"] = "SmartLight"
	// m["0004a30b00e9493f"] = "SmartLight"
	// m["0004a30b00e932aa"] = "SmartLight"
	// m["0004a30b00e9325a"] = "SmartLight"
	// m["0004a30b00e94955"] = "SmartLight"
	// m["0004a30b00e99654"] = "SmartLight"
	// m["0004a30b00e94935"] = "SmartLight"
	// m["0004a30b00e96769"] = "SmartLight"
	// m["0004a30b00e9a1e1"] = "SmartLight"
	// m["0004a30b00e9b58b"] = "SmartLight"
	// m["0004a30b00e96567"] = "SmartLight"
	// m["0004a30b00e95b95"] = "SmartLight"
	// m["0004a30b00e9d7f6"] = "SmartLight"
	// m["0004a30b00e999b3"] = "SmartLight"
	// m["0004a30b00e982fb"] = "SmartLight"
	// m["0004a30b00e9867c"] = "SmartLight"
	// m["0004a30b00e9b7ca"] = "SmartLight"
	// m["0004a30b00e964cb"] = "SmartLight"
	// m["0004a30b00e96768"] = "SmartLight"
	// m["0004a30b00e96765"] = "SmartLight"
	// m["0004a30b00e92e56"] = "SmartLight"
	// m["0004a30b00e931a9"] = "SmartLight"
	// m["0004a30b00e94a76"] = "SmartLight"
	// m["0004a30b00e9a22b"] = "SmartLight"
	// m["0004a30b00e95f76"] = "SmartLight"
	// m["0004a30b00e999ba"] = "SmartLight"
	// m["0004a30b00e94b81"] = "SmartLight"
	// m["0004a30b00e98c37"] = "SmartLight"
	// m["0004a30b00e97b2f"] = "SmartLight"
	// m["0004a30b00e9be97"] = "SmartLight"
	// m["0004a30b00e93245"] = "SmartLight"
	// m["0004a30b00e94545"] = "SmartLight"
	// m["0004a30b00e94b99"] = "SmartLight"
	// m["0004a30b00e99b8a"] = "SmartLight"
	// m["0004a30b00e9305d"] = "SmartLight"
	// m["0004a30b00e9a223"] = "SmartLight"
	// m["0004a30b00e9e99a"] = "SmartLight"
	// m["0004a30b00e9e93e"] = "SmartLight"
	// m["0004a30b00e9d69d"] = "SmartLight"
	// m["0004a30b00e94a99"] = "SmartLight"
	// m["0004a30b00e94473"] = "SmartLight"
	// m["0004a30b00e9bf24"] = "SmartLight"
	// m["0004a30b00e947ab"] = "SmartLight"
	// m["0004a30b00e9e9ce"] = "SmartLight"
	// m["0004a30b00e9d364"] = "SmartLight"
	// m["0004a30b00e94314"] = "SmartLight"
	// m["0004a30b00e964d5"] = "SmartLight"
	// m["0004a30b00e933da"] = "SmartLight"
	// m["0004a30b00e94a79"] = "SmartLight"
	// m["0004a30b00e94a9d"] = "SmartLight"
	// m["0004a30b00e9bd4a"] = "SmartLight"
	// m["0004a30b00e97ba5"] = "SmartLight"
	// m["0004a30b00e95a73"] = "SmartLight"
	// m["0004a30b00e95dad"] = "SmartLight"
	// m["0004a30b00e940a1"] = "SmartLight"
	// m["0004a30b00e933c5"] = "SmartLight"
	// m["0004a30b00e93136"] = "SmartLight"
	// m["0004a30b00e99656"] = "SmartLight"
	// m["0004a30b00e942f0"] = "SmartLight"
	// m["0004a30b00e9b90a"] = "SmartLight"
	// m["0004a30b00e92f6f"] = "SmartLight"
	// m["0004a30b00e96650"] = "SmartLight"
	// m["0004a30b00e95f71"] = "SmartLight"

	for {
		incoming := <-c
		s := strings.Split(incoming[0], "/")
		devEUI := s[3]
		// json.Unmarshal([]byte(incoming[1]), &litr)
		// var sdt strings.Builder // string device type
		// applicationName := litr.ApplicationName
		// deviceId := litr.DevEUI
		// _, prs := m[deviceId]
		// if prs {
		// 	fmt.Printf("Device is present in map %s\n", deviceId)
		// }

		// Filter by SmartCampus3.0 applications
		// if applicationName == "WaterTankLevels" ||
		// 	applicationName == "Hydrometers" ||
		// 	applicationName == "GaugePressures" ||
		// 	applicationName == "GMS" ||
		// 	applicationName == "SmartLights" {
		// fmt.Printf("Parse Infrastructure application payload data %s\n", applicationName)
		sbPubTopic.Reset()
		sbPubTopic.WriteString("OpenDataTelemetry/SmartCampusMaua/SmartLight/")
		sbPubTopic.WriteString(devEUI)
		sbPubTopic.WriteString("/rx/lns_imt")
		// fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
		token := pClient.Publish(sbPubTopic.String(), byte(pQos), false, incoming[1])
		token.Wait()
		// } else {
		// 	fmt.Printf("Application  %s is not a SmartCampus3.0 compatible\n", applicationName)
		// }
	}
}
