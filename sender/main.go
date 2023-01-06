package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"

	"github.com/streadway/amqp"
)

const queueService = "QueueService1"

func main() {
	// Define RabbitMQ server URL
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new RabbitMQ connection
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}

	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err = channelRabbitMQ.QueueDeclare(
		queueService, // queue name
		true,         // durable
		false,        // auto delete
		false,        // exclusive
		false,        // no waits
		nil,          // arguments
	)
	if err != nil {
		panic(err)
	}

	// Create a new Fiber instance
	app := fiber.New()

	// add middleware
	app.Use(
		logger.New(),
	)

	// add route
	app.Get("/send", func(ctx *fiber.Ctx) error {
		// Create a message to publish
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(ctx.Query("msg")),
		}

		err := channelRabbitMQ.Publish(
			"", // exchange
			queueService,
			false,
			false,
			message,
		)
		if err != nil {
			return err
		}
		return nil
	})

	// start fiber API server
	log.Fatal(app.Listen(":3000"))
}
