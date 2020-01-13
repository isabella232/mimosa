package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"
	"net/http"
	"bytes"

	jwt "github.com/dgrijalva/jwt-go"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type responseMessage struct {
	Topic    string
	Qos      byte
	Retained bool
	Payload  interface{}
}

type Command struct {
	Command string   `json:"command"`
	Targets []string `json:"targets"`
}

type Host struct {
	Name        string `json:"name"`
	PrivateIPv4 string `json:"privateIPv4"`
	PrivateIPv6 string `json:"privateIPv6"`
}

func main() {
	log.Println("[main] Entered")

	responseChannel := make(chan responseMessage)

	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("./roots.pem")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	log.Println("[main] Creating TLS Config")
	config := &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{},
		MinVersion:         tls.VersionTLS12,
	}

	projectID := "mimosa-andy-260909"
	region := "us-central1"
	registryID := "edges"
	deviceID := "source-5d4718f9-f2c2-493d-99f8-7598a54892cc"
	clientID := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s",
		projectID,
		region,
		registryID,
		deviceID)

	log.Println("[main] Creating MQTT Client Options")
	opts := MQTT.NewClientOptions()

	broker := fmt.Sprintf("ssl://%v:%v", "mqtt.googleapis.com", "8883")
	log.Printf("[main] Broker '%v'", broker)

	opts.AddBroker(broker)
	opts.SetClientID(clientID).SetTLSConfig(config)

	opts.SetUsername("unused")

	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = jwt.StandardClaims{
		Audience:  projectID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	log.Println("[main] Load Private Key")
	keyBytes, err := ioutil.ReadFile("./key.pem")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[main] Parse Private Key")
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[main] Sign String")
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Fatal(err)
	}

	opts.SetPassword(tokenString)

	// Incoming
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("[handler] Topic: %v\n", msg.Topic())
		fmt.Printf("[handler] Payload: %v\n", msg.Payload())
	})

	log.Println("[main] MQTT Client Connecting")
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	go subscribeForCommands(client, deviceID, responseChannel)

	wg := sync.WaitGroup{}
	wg.Add(1)
	sendMessages(client, responseChannel, &wg)
	wg.Wait()

	log.Println("[main] MQTT Client Disconnecting")
	client.Disconnect(250)
}

func sendMessages(client MQTT.Client, responseChannel chan responseMessage, wg *sync.WaitGroup) {
	log.Println("[sendMessages] waiting to send messages")
	defer wg.Done()
	for message := range responseChannel {
		log.Println("Sending message")
		token := client.Publish(
			message.Topic,
			message.Qos,
			message.Retained,
			message.Payload)
		token.WaitTimeout(5 * time.Second)
	}
}

func sendBoltCommand(cmdToSend []byte) (string, error) {
	req, err := http.NewRequest("POST", "http://localhost:4567/v1/command", bytes.NewBuffer(cmdToSend))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	// saving data
	log.Println("sending update")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

func subscribeForCommands(client MQTT.Client, deviceID string, responseChannel chan responseMessage) {
	topic := struct {
		commands  string
		telemetry string
	}{
		commands:  fmt.Sprintf("/devices/%v/commands/#", deviceID),
		telemetry: fmt.Sprintf("/devices/%v/events", deviceID),
	}

	log.Println("[subscribeForCommands] Creating Subscription")
	client.Subscribe(topic.commands, 0, func(client MQTT.Client, msg MQTT.Message) {
		log.Printf("[handler] Topic: %v\n", msg.Topic())
		log.Printf("[handler] Payload: %v\n", string(msg.Payload()))
		log.Printf("[handler] MessageID: %d\n", msg.MessageID())
		cmd := Command{}
		err := json.Unmarshal(msg.Payload(), &cmd)
		if err != nil {
			log.Println(err)
		}
		output, err := sendBoltCommand(msg.Payload())
		if err != nil {
			log.Println(err)
		}
		responseChannel <- responseMessage{
			Topic:    topic.telemetry,
			Qos:      0,
			Retained: false,
			Payload:  string(output),
		}
	})

}
