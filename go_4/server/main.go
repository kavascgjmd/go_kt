package main

import (
	"log"
	"net"
	pb "go4/proto"

	"google.golang.org/grpc"
)

type helloServer struct{
	pb.GreetServiceServer
}

const (
	port = ":8080"
)

func main(){
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to start the server %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreetServiceServer(grpcServer, &helloServer{})
	err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("ji %v",err)
	}
}