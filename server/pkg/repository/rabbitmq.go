package repository

import (
	"context"
	"encoding/json"
	"log"
	entity "server/model/entity"
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

func NewRabbitMQ(cfg ConfigMQ) (chan entity.Task, chan entity.Worker, error) {
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

	_, err = channelRabbitMQ.QueueDeclare(
		"SendTask", // queue name
		true,       // durable
		false,      // auto delete
		false,      // exclusive
		false,      // no wait
		nil,        // arguments
	)
	if err != nil {
		return nil, nil, err
	}
	_, err = channelRabbitMQ.QueueDeclare(
		"SendResult", // queue name
		true,         // durable
		false,        // auto delete
		false,        // exclusive
		false,        // no wait
		nil,          // arguments
	)
	if err != nil {
		return nil, nil, err
	}
	_, err = exchangeChannel.QueueDeclare(
		"SendHeartBeat", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		return nil, nil, err
	}

	err = exchangeChannel.ExchangeDeclare(
		"DurationExchange", // exchange name
		"fanout",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)

	if err != nil {
		return nil, nil, err
	}

	channel := make(chan entity.Task, 15)
	go func() {
		messages, err := channelRabbitMQ.Consume(
			"SendResult", // queue
			"",           // consumer
			true,         // auto-ack
			false,        // exclusive
			false,        // no-local
			false,        // no-wait
			nil,          // args
		)
		if err != nil {
			log.Fatalf("failed to register a consumer. Error: %s", err)
		}

		var forever chan struct{}

		go func() {
			for message := range messages {
				var task entity.Task
				json.Unmarshal(message.Body, &task)
				channel <- task
			}
		}()

		<-forever
	}()

	heartBeatChannel := make(chan entity.Worker, 15)

	go func() {
		messages, err := channelRabbitMQ.Consume(
			"SendHeartBeat", // queue
			"",              // consumer
			true,            // auto-ack
			false,           // exclusive
			false,           // no-local
			false,           // no-wait
			nil,             // args
		)
		if err != nil {
			log.Fatalf("failed to register a consumer. Error: %s", err)
		}

		var forever chan struct{}

		go func() {
			for message := range messages {
				var worker entity.Worker
				json.Unmarshal(message.Body, &worker)
				heartBeatChannel <- worker
			}
		}()

		<-forever
	}()

	return channel, heartBeatChannel, nil
}

func SendMessage(bytes []byte, queueName string) error {
	message := amqp091.Publishing{
		ContentType: "text/plain",
		ReplyTo:     "SendResult",
		Body:        []byte(bytes),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if queueName == "SendDuration" {
		if err := channelRabbitMQ.PublishWithContext(ctx,
			"DurationExchange", // exchange
			"",                 // routing key
			false,              // mandatory
			false,              // immediate
			message,            // message to publish
		); err != nil {
			return err
		}
		return nil
	}

	if err := channelRabbitMQ.PublishWithContext(ctx,
		"",        // exchange
		queueName, // queue name
		false,     // mandatory
		false,     // immediate
		message,   // message to publish
	); err != nil {
		return err
	}
	return nil
}

func Close() {
	connectRabbitMQ.Close()
	channelRabbitMQ.Close()
	exchangeChannel.Close()
}
