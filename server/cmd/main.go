package main

import (
	"log"
	"os"

	"server/pkg/handler"
	"server/pkg/repository"
	"server/pkg/service"
	"server/utility"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("can not initialize configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("can not load env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:         viper.GetString("db.host"),
		Port:         viper.GetString("db.port"),
		Username:     viper.GetString("db.username"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		DBName:       viper.GetString("db.dbname"),
		SSLMode:      viper.GetString("db.sslmode"),
		MigrationURL: viper.GetString("db.migration_url"),
	})

	if err != nil {
		log.Fatalf("can not initialize db: %s", err.Error())
	}

	resultChannel, heartBeatChannel, err := repository.NewRabbitMQ(repository.ConfigMQ{
		Url: viper.GetString("message-broker.url"),
	})

	if err != nil {
		log.Fatalf("can not initialize MQ: %s", err.Error())
	}
	defer repository.Close()

	repos := repository.NewRepository(db, resultChannel, heartBeatChannel)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(utility.Server)

	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("can not run http server: %s", err.Error())
	}
}

//TODO: swagger

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
