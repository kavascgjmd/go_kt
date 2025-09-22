package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main(){
   conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
   failedOnError(err, "failed to connect Rabbitmq")
   defer conn.Close();
   ch , err := conn.Channel()
   failedOnError(err, "failed to create channel");
   defer ch.Close()
   q, err := ch.QueueDeclare(
	"hello",
	false,
	false,
	false,
	false,
	nil,
   )

   failedOnError(err, "failed declare a queue")
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

   defer cancel()

   body := "Hello World!"
   ch.PublishWithContext(ctx, "", q.Name, false, false,
   amqp.Publishing{
	ContentType: "text/plain",
	Body: []byte(body),
   } )

}

func failedOnError(err error, msg string){
	log.Fatalf("%s %s", msg, err);
}