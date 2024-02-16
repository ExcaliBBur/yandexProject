package main

import (
	"agent/pkg/repository"
	"agent/pkg/service"
	"log"

	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("can not initialize configs: %s", err.Error())
	}

	taskChannel, durationChannel, err := repository.NewRabbitMQ(repository.ConfigMQ{
		Url: viper.GetString("message-broker.url"),
	})

	if err != nil {
		log.Fatalf("can not initialize MQ: %s", err.Error())
	}
	defer repository.Close()

	repository := repository.NewRepository(taskChannel, durationChannel)
	service.NewService(repository)

	var forever chan bool
	<-forever
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
