package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Gandalf-Rus/distributed-calc2.0/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	host := "localhost"
	port := "5000"

	addr := fmt.Sprintf("%s:%s", host, port) // используем адрес сервера
	// установим соединение
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}
	// закроем соединение, когда выйдем из функции
	defer conn.Close()

	grpcClient := proto.NewNodeServiceClient(conn)
	nodes, err := grpcClient.GetNodes(context.Background(), &proto.GetNodesRequest{
		AgentId:     1,
		FreeWorkers: 3,
	})
	fmt.Printf("%v\n%v", nodes, err)
}
