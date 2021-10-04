package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ungame/go-rmq/client"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
	}

	err = client.LoadConfigsFromEnv()
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	configs := client.NewConfigs()
	fmt.Println(configs)

	rmqClient, err := client.NewClient(configs)
	if err != nil {
		log.Panicln(err)
	}

	defer rmqClient.Close()
}
