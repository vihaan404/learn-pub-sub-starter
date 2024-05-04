package main

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/vihaan404/learn-pub-sub-starter/internal/gamelogic"
	"github.com/vihaan404/learn-pub-sub-starter/internal/pubsub"
	"github.com/vihaan404/learn-pub-sub-starter/internal/routing"
	"log"
	"os"
	"os/signal"
)

const URL = "amqp://guest:guest@localhost:5672/"

func main() {
	conn, err := amqp091.Dial(URL)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	fmt.Println("Peril game connect to rabbitmq")
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not open a channel: %v", err)
	}

	defer conn.Close()

	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		command := input[0]
		switch command {
		case "pause":
			fmt.Println("publishing the pause message")
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				log.Fatalf("could not publish a message: %v", err)
			}
		case "resume":
			fmt.Println("publishing the resume message")
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				log.Fatalf("could not publish a message: %v", err)
			}
		case "quit":
			fmt.Println("quiting the application ......")
			return
		default:
			fmt.Println("Unknown command")
			continue

		}

	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
