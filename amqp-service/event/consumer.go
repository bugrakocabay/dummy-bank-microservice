package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// NewConsumer creates a connection for consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// setup creates a channel with topic
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		err = ch.QueueBind(q.Name, s, "logs_topic", false, nil)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("waiting for message on [Exchange, Queue] [logs_topic, %s]", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Println("sendErrorLog error: cant marshal json:", err)
		return err
	}
	request, err := http.NewRequest(http.MethodPost, "http://logger-service/logs/create-error", bytes.NewBuffer(jsonData))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("sendErrorLog error: cant send response:", err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		log.Println("sendErrorLog error: calling auth service:", err)
		return err
	}

	log.Println("logged successfully")
	return nil
}
