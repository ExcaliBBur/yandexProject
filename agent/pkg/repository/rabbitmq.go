package repository

import (
	entity "agent/model/entity"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type ConfigMQ struct {
	Url string
}

var (
	connectRabbitMQ *amqp091.Connection
	channelRabbitMQ *amqp091.Channel
	exchangeChannel *amqp091.Channel
)

func NewRabbitMQ(cfg ConfigMQ) (<-chan entity.Task, <-chan entity.Duration, error) {
	var err error
	connectRabbitMQ, err = amqp091.Dial(cfg.Url)

	if err != nil {
		return nil, nil, err
	}

	channelRabbitMQ, err = connectRabbitMQ.Channel()
	if err != nil {
		return nil, nil, err
	}

	exchangeChannel, err = connectRabbitMQ.Channel()
	if err != nil {
		return nil, nil, err
	}

	taskChannel := make(chan entity.Task, 15)

	go func() {
		messages, err := channelRabbitMQ.Consume(
			"SendTask", // queue
			"",         // consumer
			true,       // auto-ack
			false,      // exclusive
			false,      // no-local
			false,      // no-wait
			nil,        // args
		)
		if err != nil {
			log.Fatalf("failed to register a consumer. Error: %s", err)
		}

		var forever chan struct{}

		go func() {
			for message := range messages {
				var task entity.Task
				json.Unmarshal(message.Body, &task)
				taskChannel <- task
			}
		}()

		<-forever
	}()

	durationChannel := make(chan entity.Duration, 15)

	go func() {
		err = exchangeChannel.ExchangeDeclare(
			"DurationExchange", // name
			"fanout",           // type
			true,               // durable
			false,              // auto-deleted
			false,              // internal
			false,              // no-wait
			nil,                // arguments
		)
		if err != nil {
			log.Fatalf("failed to register a consumer. Error: %s", err)
		}

		q, err := exchangeChannel.QueueDeclare(
			"",    // name
			false, // durable
			false, // delete when unused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)
		// failOnError(err, "Failed to declare a queue")
		if err != nil {
			log.Fatalf("failed to register a consumer. Error: %s", err)
		}

		err = exchangeChannel.QueueBind(
			q.Name,             // queue name
			"",                 // routing key
			"DurationExchange", // exchange
			false,
			nil,
		)
		// failOnError(err, "Failed to bind a queue")
		if err != nil {
			log.Fatalf("failed to register a consumer. Error: %s", err)
		}

		messages, err := exchangeChannel.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		// failOnError(err, "Failed to register a consumer")
		if err != nil {
			log.Fatalf("failed to register a consumer. Error: %s", err)
		}

		if err != nil {
			log.Fatalf("failed to register3 a consumer. Error: %s", err)
		}

		var forever chan struct{}

		go func() {
			for message := range messages {
				var duration entity.Duration
				json.Unmarshal(message.Body, &duration)
				durationChannel <- duration
			}
		}()

		<-forever
	}()

	return taskChannel, durationChannel, nil
}

func SendMessage(bytes []byte, queueName string, delay int) error {
	message := amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(bytes),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := channelRabbitMQ.PublishWithContext(ctx,
		"",        // exchange
		queueName, // queue name
		false,     // mandatory
		false,     // immediate
		message,   // message to publish
	); err != nil {
		return err
	}

	time.Sleep(time.Second * time.Duration(delay))

	return nil
}

func Close() {
	connectRabbitMQ.Close()
	channelRabbitMQ.Close()
}
