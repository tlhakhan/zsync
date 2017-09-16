package main

import (
	"fmt"
	"log"
	pb "zsync/service"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8089", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewZsyncClient(conn)
	response, err := c.GetFilesystems(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}
