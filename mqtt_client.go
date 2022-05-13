package go_mqtt_client

import (
	"crypto/rand"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

type MqttClient struct {
	LocalMqttClient mqtt.Client
	Options         *mqtt.ClientOptions
	Topics          []string
	broker          string
}

// Generates a random string
func tokenGenerator() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		log.Error().Err(err).Msg("Generating a random token.")
		return ""
	}
	return fmt.Sprintf("%x", b)
}

var onConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	r := client.OptionsReader()
	log.Debug().Str("clientId", r.ClientID()).Msg("Connected to MQTT broker.")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, error error) {
	r := client.OptionsReader()
	log.Debug().Str("clientId", r.ClientID()).Err(error).Msg("Connection lost.")
}

// MqttClientInit returns an MqttClient initialised with the MessageHandler function given in mqttMessageHandler, the
// list of subscribed topics given in subTopics, the broker address given in broker, and default options for the
// KeepAlive (5 s), PingTimeout (1 s), AutoReconnect (false) and ClientId (a 4 byte hex string). These options can be
// changed by accessing the mqtt.ClientOptions structure from MqttClient.Options.
func MqttClientInit(mqttMessageHandler mqtt.MessageHandler, subTopics []string, broker string) MqttClient {
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetDefaultPublishHandler(mqttMessageHandler)
	opts.SetClientID(tokenGenerator())
	opts.SetOnConnectHandler(onConnectHandler)
	opts.SetConnectionLostHandler(connectionLostHandler)

	return MqttClient{
		Options: opts,
		Topics:  subTopics,
		broker:  broker,
	}
}

// Destroy disconnects from the broker.
func (mc *MqttClient) Destroy() {
	mc.LocalMqttClient.Disconnect(250)
}

// Publish asynchronously publishes the payload on a topic. Prints out an error if it fails.
func (mc *MqttClient) Publish(payload interface{}, topic string) {
	token := mc.LocalMqttClient.Publish(topic, 0, false, payload)
	go func() {
		token.Wait()
		err := token.Error()
		if err != nil {
			log.Error().Interface("payload", payload).Str("topic", topic).Err(token.Error()).Msg("Failed to Publish")
		}
	}()
}

// SubscribeToTopic Subscribe to an additional topic. Ideally, the user should pass all the topics that they want to
// subscribe to MqttClientInit.
func (mc *MqttClient) SubscribeToTopic(topic string) error {
	if token := mc.LocalMqttClient.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Info().Str("Topic", topic).Msg("Subscribed to topic.")
	return nil
}

// Start starts the MQTT client and subscribes to the list of topics given in MqttClientInit
func (mc *MqttClient) Start() error {
	mc.LocalMqttClient = mqtt.NewClient(mc.Options)
	if token := mc.LocalMqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Info().Str("Broker", mc.broker).Msg("Connected to MQTT broker.")

	for _, topic := range mc.Topics {
		err := mc.SubscribeToTopic(topic)
		if err != nil {
			log.Error().Str("topic", topic).Err(err).Msg("Subscribing to topic.")
		}
	}

	return nil
}
