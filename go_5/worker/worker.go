package main

import (
	"fmt"
	"os"
	"syscall"
	"github.com/IBM/sarama"
	"os/signal"
)

func main(){
	topic := "comments";
	consumer ,err := connectConsumer([]string{"localhost:29092"})
	if err != nil{
		panic(err)
	}
	partition, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil{
		panic(err)
	}
	sigchain := make(chan os.Signal, 1)
	signal.Notify(sigchain, syscall.SIGINT, syscall.SIGTERM)
	msgCount := 0
	doneCh := make(chan struct{})

	go func(){
		for {
			select {
			case err := <-partition.Errors():
			     fmt.Println(err)
		    case msg := <-partition.Messages():
			     msgCount++;
				 fmt.Printf("recivd message %d , topic %s , message %s \n", msgCount, string(msg.Topic), string(msg.Value))
            case <-sigchain:
				 fmt.Println("interruption detected")
				 doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Proccessed", msgCount, "messages")
	if err = consumer.Close(); err != nil{
		panic(err)
	}

}

func connectConsumer(brokerUrl []string)(sarama.Consumer, error){
	config := sarama.NewConfig();
	config.Consumer.Return.Errors = true
	conn, err := sarama.NewConsumer(brokerUrl, config);
	if err != nil{
		return nil, err
	}
	return conn, nil
}