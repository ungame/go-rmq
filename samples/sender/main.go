package main

import (
	"flag"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/ungame/go-rmq/client"
	"github.com/ungame/go-rmq/samples/models"
	"github.com/ungame/go-rmq/samples/utils"
	"log"
	"time"
)

var spam int

func init() {
	flag.IntVar(&spam, "spam", 10, "spam messages to message broker")
	client.LoadConfigsFromFlags(flag.CommandLine)
	flag.Parse()
}

func main() {
	configs := client.NewConfigs()

	rmqClient, err := client.NewClient(configs)
	if err != nil {
		log.Panicln(err)
	}
	defer rmqClient.Close()

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

	send(rmqClient, spam)

	log.Println("Sender finished!")
}

func send(rmqClient client.Client, qty int) {
	for _, user := range newUsers(qty) {

		publisherOpts := utils.GetDefaultPublisherOptions(user)

		err := rmqClient.Send(publisherOpts)
		if err != nil {
			log.Panicln(err)
		}

		log.Println(" [*] Message sent!")
	}
}

func newUsers(qty int) []*models.User {
	users := make([]*models.User, 0, qty)
	for i := 0; i < qty; i++ {
		users = append(users, newUser())
	}
	return users
}

func newUser() *models.User {
	return &models.User{
		ID:        uuid.NewString(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}
