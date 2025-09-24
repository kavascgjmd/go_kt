package main

import (
	"log"
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
   msgs , err := ch.Consume(
	q.Name, "", true, false, false, false, nil,
   )

   failedOnError(err, "failed to recive")
   var forever chan struct {}
   go func (){
	  for d := range msgs{
		log.Printf("Recieved a message: %s", d.Body)
	  }
   }()

   <-forever
}


func failedOnError(err error, msg string){
	log.Fatalf("%s %s", msg, err);
}