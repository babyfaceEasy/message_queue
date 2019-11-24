package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	pb "github.com/me/message_queue/messagequeue"
	"github.com/me/message_queue/models"
	"github.com/subosito/gotenv"
	"google.golang.org/grpc"
)

type messageQueueServer struct {
	pb.UnimplementedMessageQueueServer
}

type tempStruct struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func (m *messageQueueServer) CreateMessage(ctx context.Context, queueReq *pb.QueueMessage) (*pb.Response, error) {

	// implement the functionality to create messages
	// sanitize and validate data
	// call Queue and Message models
	queueDetails := new(models.Queue)
	err := queueDetails.GetQueueByName(queueReq.Queue.GetName())
	if err != nil {
		return &pb.Response{Status: pb.Response_ERROR}, nil
	}
	// queue found
	message := models.Message{Message: queueReq.GetMessageJson()}
	_, err = message.CreateMessage(queueReq.GetQueue().GetName())

	if err != nil {
		return &pb.Response{Status: pb.Response_ERROR}, nil
	}

	return &pb.Response{Status: pb.Response_SUCCESS}, nil
}

func (m *messageQueueServer) GetMessage(ctx context.Context, queueReq *pb.QueueName) (*pb.QueueMessage, error) {
	// return the oldest message in the queue
	// get the queue then return the message
	queueDetails := models.Queue{}
	err := queueDetails.GetQueueByName(queueReq.GetName())

	if err != nil {
		return &pb.QueueMessage{
			Queue:       &pb.QueueName{},
			MessageJson: "",
		}, err
	}

	message, err := queueDetails.GetMessage()
	if err != nil {
		return &pb.QueueMessage{
			Queue:       &pb.QueueName{},
			MessageJson: "",
		}, err
	}

	return &pb.QueueMessage{
		Queue:       queueReq,
		MessageJson: message.Message,
	}, nil

}

func newServer() *messageQueueServer {
	m := &messageQueueServer{}
	return m
}

func init() {
	err := gotenv.Load()
	if err != nil {
		panic("you need a .env file in order to run this project.")
	}
}

func main() {
	/*
			db, err := gorm.Open("mysql", "root:root@/message_queue?charset=utf8&parseTime=True&loc=Local")
			if err != nil {
				panic("Failed to connect to database")
			}

			// create new queue

			newQueue := models.Queue{Name: "resend_orders"}
			response := db.NewRecord(newQueue)
			db.Create(&newQueue)

			fmt.Println(response)

			// get first queue details
			queue := new(models.Queue)
			db.First(&queue)
			fmt.Printf("%v", queue)

			defer db.Close()

		newQueue := models.Queue{Name: "another_try"}
		data, err := newQueue.CreateQueue()

		if err != nil {
			panic(err)
		}

		fmt.Printf("%v", data)
	*/

	/*
		db, err := gorm.Open("mysql", "root:root@/message_queue?charset=utf8&parseTime=True&loc=Local")
		defer db.Close()
		if err != nil {
			panic("Failed to connect to database")
		}
		tempData := tempStruct{Name: "Olakunle Odegbaro", ID: 2}
		stringData, _ := json.Marshal(tempData)
		fmt.Println(string(stringData))
		// create message
		newMessage := models.Message{
			QueueID: 1,
			Message: string(stringData),
		}
		db.Create(&newMessage)

		message := new(models.Message)
		db.First(&message)
		fmt.Printf("%+v", message)
	*/
	/*
		queueDetails := models.Queue{}
		err := queueDetails.GetQueueByName("resend_orders")

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%+v", queueDetails)
	*/

	/*
		db, err := gorm.Open("mysql", "root:root@/message_queue?charset=utf8&parseTime=True&loc=Local")
		defer db.Close()
		if err != nil {
			panic("Failed to connect to database")
		}
		queue := new(models.Queue)
		db.First(&queue)
		message, err := queue.GetMessage()
		if err != nil {
			fmt.Println(err)
		}
		//messages := []models.Message

		//fmt.Printf("%+v\n", queue)
		fmt.Printf("%+v\n", message)
	*/
	/*
		newQueue := models.Queue{Name: "failed_emails"}
		data, err := newQueue.CreateQueue()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v", data)
	*/

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", os.Getenv("GRPC_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Listening on %s\n", os.Getenv("GRPC_PORT"))
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMessageQueueServer(grpcServer, newServer())
	grpcServer.Serve(lis)

}
