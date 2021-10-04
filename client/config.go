package client

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	rmqUser string
	rmqPass string
	rmqHost string
	rmqPort int
)

type Config interface{
	URL() string
}

func LoadConfigsFromEnv() error {
	rmqUser = os.Getenv("RABBITMQ_USER")
	rmqPass = os.Getenv("RABBITMQ_PASS")
	rmqHost = os.Getenv("RABBITMQ_HOST")
	var err error
	rmqPort, err = strconv.Atoi(os.Getenv("RABBITMQ_PORT"))
	if err != nil {
		return err
	}
	return nil
}

func LoadConfigsFromFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&rmqUser, "rmq_user", "guest", "rabbitmq username")
	flagSet.StringVar(&rmqPass, "rmq_pass", "guest", "rabbitmq password")
	flagSet.StringVar(&rmqHost, "rmq_host", "localhost", "rabbitmq host")
	flagSet.IntVar(&rmqPort, "rmq_port", 5672, "rabbitmq port")
}

type config struct {
	user string
	pass string
	host string
	port int
}

func NewConfigs() Config {
	return &config{
		user: rmqUser,
		pass: rmqPass,
		host: rmqHost,
		port: rmqPort,
	}
}

func (c *config) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", c.user, c.pass, c.host, c.port)
}

func (c *config) String() string {
	return c.URL()
}
