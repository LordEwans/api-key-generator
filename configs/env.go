package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Helps retrieve the MongoDB URI fron the env file
func EnvMongoURI() string {
	env, err := os.LookupEnv("MONGOURI")
	if !err {
		err := godotenv.Load()
		if err != nil {
			err := godotenv.Load(".env.test")
			if err != nil {
				log.Fatal("Error loading .env file: ", err)
			}
		}
		return os.Getenv("MONGOURI")
	}
	return env
}

// Helps retrieve the RabbitMQ AMQP address fron the env file
func EnvRabbitMQ() string {
	env, err := os.LookupEnv("RABBITMQ")
	if !err {
		err := godotenv.Load()
		if err != nil {
			err := godotenv.Load(".env.test")
			if err != nil {
				log.Fatal("Error loading .env file: ", err)
			}
		}
		return os.Getenv("RABBITMQ")
	}
	return env
}
