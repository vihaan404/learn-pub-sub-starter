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

//exchange: peril_direct (this is a constant in the internal/routing package)
//queueName: pause.username where username is the user's input. The pause section of the name is the routing key constant in the internal/routing package.
//routingKey: pause (this is a constant in the internal/routing package)
//simpleQueueType: transient

func main() {
	conn, err := amqp091.Dial(URL)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game connect to rabbitmq")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could get the username: %v", err)
	}
	queueName := fmt.Sprintf("pause.%s", username)
	_, _, _ = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, 1)

	gameState := gamelogic.NewGameState(username)
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		command := input[0]
		switch command {
		case "spawn":
			err := gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			_, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println(err)
				continue

			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("dont spam messages")
		case "quit":
			fmt.Println("quiting the game ....")
			return
		default:
			fmt.Println("unknown command")
			continue
		}

	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
