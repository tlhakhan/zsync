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
	zpool := flag.String("Z", "zones", "zpool to sync")
	flag.Parse()

	// start up a zfs daemon
	zfsD := zfs.NewDaemon(*zpool)

	for !zfsD.Ready {
		// block
	}

	// start up rpc connection
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *server, *port), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// rpc client
	c := pb.NewZsyncClient(conn)

	// ask existence

	for _, filesystem := range zfsD.Filesystems {
		log.Println("Asking for existence of zfs filesystem: ", filesystem.Name)
		// ask questions
		response, err := c.Exists(context.Background(), &pb.Empty{})
		if err != nil {
			log.Fatal(err)
		}
		if response.Name == nil {
			log.Println("Server did not find your zfs filesystem.")
		} else {
			log.Println("Server found your zfs filesystem.")
		}

	}
	/*


	 */
}
