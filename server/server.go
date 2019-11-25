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
	// TODO: sanitize and validate data
	queueDetails := new(models.Queue)
	err := queueDetails.GetQueueByName(queueReq.Queue.GetName())
	if err != nil {
		return &pb.Response{Status: pb.Response_ERROR}, nil
	}
	// queue found
	message := models.Message{Message: queueReq.GetMessageJson()}
	_, err = message.CreateMessage(queueReq.GetQueue().GetName())

	if err != nil {
		return &pb.Response{Status: pb.Response_ERROR}, err
	}

	return &pb.Response{Status: pb.Response_SUCCESS}, nil
}

func (m *messageQueueServer) GetMessage(ctx context.Context, queueReq *pb.QueueName) (*pb.QueueMessage, error) {
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

// UpdateMessageStatus sets the message status based on the provided logic
func (m *messageQueueServer) UpdateMessageStatus(ctx context.Context, messageID *pb.MessageID) (*pb.Response, error) {
	message := new(models.Message)
	err := message.GetMessageByID(messageID.GetId())
	if err != nil {
		return &pb.Response{Status: pb.Response_ERROR}, err
	}

	err = message.UpdateStatus()
	if err != nil {
		return &pb.Response{Status: pb.Response_ERROR}, err
	}

	return &pb.Response{Status: pb.Response_SUCCESS}, nil
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
