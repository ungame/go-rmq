package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ungame/go-rmq/client"
	"github.com/ungame/go-rmq/samples/models"
	"github.com/ungame/go-rmq/samples/utils"
	"log"
	"os"
	"os/signal"
)

func init() {
	client.LoadConfigsFromFlags(flag.CommandLine)
	flag.Parse()
}

func main() {
	configs := client.NewConfigs()

	rmqClient, err := client.NewClient(configs)
	if err != nil {
		log.Panicln(err)
	}

	exchangeOpts := utils.GetDefaultExchangeOptions()

	err = rmqClient.ExchangeDeclare(exchangeOpts)
	if err != nil {
		log.Panicln(err)
	}

	queueOpts := utils.GetDefaultQueueOptions()

	err = rmqClient.QueueDeclare(queueOpts)
	if err != nil {
		log.Panicln(err)
	}

	rmqClient.HandleFunc(models.User{}, func(data []byte) error {
		var user models.User
		err := json.Unmarshal(data, &user)
		if err != nil {
			return fmt.Errorf("error on decode message: %s", err)
		}
		log.Printf("Handled successfully: %v\n", user)
		return nil
	})

	consumerOpts := utils.GetDefaultConsumerOptions()

	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)

	go func() {
		<-ctrlC
		log.Println("Closing consumer...")
		rmqClient.Close()
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	err = rmqClient.Listen(consumerOpts)
	if err != nil {
		log.Panicln(err)
	}
}
