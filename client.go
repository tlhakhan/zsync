package main

import (
	"flag"
	"fmt"
	"log"
	pb "zsync/service"
	"zsync/workers/zfs"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	server := flag.String("s", "127.0.0.1", "server to connect")
	port := flag.Int("p", 8089, "sever port to connect")
	flag.Parse()

	// start up a zfs daemon
	zfsD := zfs.NewDaemon(*zpool)

	for !zfsD.ready {
		// block
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", *server, *port), grpc.WithInsecure())

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
