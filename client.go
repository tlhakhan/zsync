package main

import (
	"sync"
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
	var wg sync.WaitGroup
	wg.Add(len(zfsD.Filesystems.List))
	fmt.Println(len(zfsD.Filesystems.List))
	for _, filesystem := range zfsD.Filesystems.List {
		go func(filesystem *pb.Dataset, wg *sync.WaitGroup) {
			defer wg.Done()

			// ask questions
			response, err := c.Exists(context.Background(), filesystem)
			if err != nil {
				log.Fatal(err)
			}
			if response.Name == "" {
				log.Printf("%s didn't find %s.", *server, filesystem.Name)
			} else {
				log.Printf("%s found %s.", *server, filesystem.Name)
			}
		}(filesystem, &wg)
	}
	wg.Wait();
	/*


	 */
}
