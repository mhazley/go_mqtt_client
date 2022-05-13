package go_mqtt_client

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"testing"
	"time"
)

var lastPayload string

var handleMqttMessages mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Debug().Str("topic", msg.Topic()).Str("payload", string(msg.Payload())).Msg("Received MQTT message")
	lastPayload = string(msg.Payload())
}

func TestMqttClientInit(t *testing.T) {
	subTopic := []string{"a", "b", "c"}
	broker := "tcp://127.0.0.1:1883"
	client := MqttClientInit(handleMqttMessages, subTopic, broker)
	err := client.Start()
	if err != nil {
		t.Error(err)
	}
	client.Destroy()
}

func TestMqttClientInitFailsToConnect(t *testing.T) {
	subTopic := []string{"a", "b", "c"}
	broker := "tcp://9.99.9.99:1883"
	client := MqttClientInit(handleMqttMessages, subTopic, broker)
	client.Options.ConnectTimeout = 1 * time.Second
	err := client.Start()
	if err.Error() != "network Error : dial tcp 9.99.9.99:1883: i/o timeout" {
		t.Error(err)
	}
	client.Destroy()
}

func TestMqttClientInitSameClientIds(t *testing.T) {
	subTopic := []string{"a", "b", "c"}
	broker := "tcp://127.0.0.1:1883"
	clientA := MqttClientInit(handleMqttMessages, subTopic, broker)
	clientA.Options.SetClientID("id1")
	err := clientA.Start()
	if err != nil {
		t.Error(err)
	}

	clientB := MqttClientInit(handleMqttMessages, subTopic, broker)
	clientB.Options.SetClientID("id1")
	err = clientB.Start()
	if err != nil {
		t.Error(err)
	}

	clientA.Destroy()
	clientB.Destroy()
}

func TestMqttClientPublish(t *testing.T) {
	broker := "tcp://127.0.0.1:1883"
	subTopicA := []string{"b"}
	clientA := MqttClientInit(handleMqttMessages, subTopicA, broker)
	err := clientA.Start()
	if err != nil {
		t.Error(err)
	}
	subTopicB := []string{"a"}
	clientB := MqttClientInit(handleMqttMessages, subTopicB, broker)
	err = clientB.Start()
	if err != nil {
		t.Error(err)
	}

	payload1 := "hello"
	clientA.Publish([]byte(payload1), "b")
	time.Sleep(1 * time.Millisecond)
	if lastPayload != payload1 {
		t.Errorf("Expected payload %s got payload %s", payload1, lastPayload)
	}

	payload2 := "world"
	clientB.Publish([]byte(payload2), "a")
	time.Sleep(1 * time.Millisecond)
	if lastPayload != payload2 {
		t.Errorf("Expected payload %s got payload %s", payload2, lastPayload)
	}

	clientA.Destroy()
	clientB.Destroy()
}
