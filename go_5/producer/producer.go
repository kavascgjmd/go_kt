package main

import (
	"encoding/json"
	"log"
	"github.com/gofiber/fiber/v3"
	"github.com/IBM/sarama"
	"fmt"
)

type Comment struct {
	Text string `form:"text" json:"text"`
}

func main(){
	app := fiber.New();
    api := app.Group("/api/v1")
	api.Post("/comments",createComment)
	app.Listen(":3000")
}

func createComment(c  fiber.Ctx) error{
	cmt := &Comment{};
	err := c.Bind().Body(cmt); if err != nil {
         log.Println(err);
		 c.Status(400).JSON(&fiber.Map{
            "sucess":false,
			"message":err,
		 })
		 return err
	}
	cmtInBytes , err := json.Marshal(cmt);
	PushCommentToQueue("comments",cmtInBytes);
	err = c.JSON(&fiber.Map{
		"sucess": true,
		"message": "comment push successfully",
		"comment" : cmt,
	})
	if err != nil{
		c.Status(400).JSON(&fiber.Map{
			"sucess": false,
			"message" : "error", 
		})
		return err
	}
	return nil
}

func PushCommentToQueue(topic string, message []byte) error {
     brokerUrl := []string{"localhost:29092"}
	 producer , err := ConnectProducer(brokerUrl)
	 if err != nil {return err}
	 defer producer.Close()
	 msg := &sarama.ProducerMessage{
		Topic: topic,
		Value : sarama.StringEncoder(message),
	 }
	 partition, offset ,err := producer.SendMessage(msg)
	 if err != nil{
		return err
	 }

	 fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n" , topic, partition, offset)
	 return nil
}

func ConnectProducer(brokerUrl [] string) (sarama.SyncProducer, error){
     config := sarama.NewConfig()
	 config.Producer.Return.Successes = true
	 config.Producer.RequiredAcks = sarama.WaitForAll
	 config.Producer.Retry.Max = 5

	 conn, err := sarama.NewSyncProducer(brokerUrl, config)
	 if err != nil {
		return nil, err
	 }
	 return conn, nil
}
