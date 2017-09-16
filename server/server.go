package main

import (
	"flag"
	"fmt"
	"net"
	"zsync/server/workers/zfs"
	pb "zsync/service"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct {
	daemon *zfs.Daemon
}

func (s *server) GetFilesystems(ctx context.Context, req *pb.Empty) (*pb.DatasetList, error) {
	return s.daemon.ListFilesystems(), nil
}

func (s *server) GetSnapshotsFor(ctx context.Context, req *pb.Dataset) (*pb.DatasetList, error) {
	return s.daemon.ListSnapshots(req), nil
}

func (s *server) Exists(ctx context.Context, req *pb.Dataset) (*pb.Dataset, error) {
	return s.daemon.FindFilesystem(req), nil
}

func main() {

	// get a port number to start on
	var port = flag.Int("p", 8089, "api listen port")
	var zpool = flag.String("Z", "zones", "zpool for zsync-api access")
	flag.Parse()

	// create the rpc server and start listening
	s := grpc.NewServer()
	zfsD := zfs.NewDaemon(*zpool)
	pb.RegisterZsyncServer(s, &server{daemon: zfsD})
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	s.Serve(ln)
}
