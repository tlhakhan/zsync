package zfs

import (
	"log"
	"strings"

	"github.com/tlhakhan/golib/cmd"

	pb "zsync/service"
)

const (
	FILESYSTEM = iota
	SNAPSHOT   = iota
)

type Daemon struct {
	Pool string `json:pool`
	// all zfs strings
	Filesystems pb.DatasetList
	Snapshots   pb.DatasetList
	ready       bool
}

func NewDaemon(pool string) *Daemon {
	d := &Daemon{Pool: pool, ready: false}
	go d.run()
	return d
}

func (d *Daemon) run() {

	// zfs list -Hro name,origin -t filesystem clusters
	fsWorker := cmd.NewWorker([]string{"zfs", "list", "-Hro", "name", "-t", "filesystem", d.Pool}, 1)
	snapWorker := cmd.NewWorker([]string{"zfs", "list", "-Hro", "name", "-t", "snapshot", d.Pool}, 10)

	// listens for new output sent on worker channels
	for {
		select {
		case fsOut := <-fsWorker:
			d.processOutput(fsOut, FILESYSTEM)
		case snapOut := <-snapWorker:
			d.processOutput(snapOut, SNAPSHOT)
		default:
		}
	}
}

func (d *Daemon) processOutput(work string, fsType int) {

	// parse
	dsl := pb.DatasetList{}
	lines := strings.Split(work, "\n")
	for _, value := range lines {
		dsl.List = append(dsl.List, &pb.Dataset{Name: value})
	}

	switch fsType {
	case FILESYSTEM:
		log.Println("Adding zfs Filesystems to Daemon struct.")
		d.Filesystems = dsl
	case SNAPSHOT:
		log.Println("Adding zfs snapshots to Daemon struct.")
		d.Snapshots = dsl
	}
	d.ready = true
}

func (d *Daemon) ListFilesystems() *pb.DatasetList {
	return &d.Filesystems
}

func (d *Daemon) FindFilesystem(dataset *pb.Dataset) *pb.Dataset {

	found := false
	for _, item := range d.Filesystems.List {
		if item.Name == dataset.Name {
			found = true
			break
		}
	}
	if found == true {
		return dataset
	}
	return &pb.Dataset{}

}

func (d *Daemon) ListSnapshots(dataset *pb.Dataset) *pb.DatasetList {

	matched := pb.DatasetList{}
	if len(dataset.Name) > 0 {
		for _, item := range d.Snapshots.List {
			if strings.Split(item.Name, "@")[0] == dataset.Name {
				matched.List = append(matched.List, item)
			}
		}
		return &matched
	} else {
		return &d.Snapshots
	}
}
