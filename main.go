package main

import (
	"fmt"
	"net"
	"os"

	"github.com/keithalucas/jsonrpc/pkg/jsonrpc"
	"github.com/keithalucas/jsonrpc/pkg/spdk"
)

func main() {

	conn, err := net.Dial("unix", "/var/tmp/spdk.sock")

	if err != nil {
		fmt.Printf("Error opening socket: %v", err)
		return
	}

	client := jsonrpc.NewClient(conn)

	errChan := client.Init()

	aio := spdk.NewAioCreate("sata1", "/dev/sda", 4096)
	client.SendMsg(aio.GetMethod(), aio)

	lvs := spdk.NewBdevLvolCreateLvstore("sata1", "longhorn")
	client.SendMsg(lvs.GetMethod(), lvs)

	lrc := spdk.NewLonghornCreateReplica("demo", 4*1024*1024*1024, "longhorn", "", 0)
	client.SendMsg(lrc.GetMethod(), lrc)
	r1 := spdk.NewLonghornCreateReplica("test", 4<<30, "longhorn", "", 0)
	client.SendCommand(r1)

	replicas := []spdk.LonghornVolumeReplica{
		spdk.LonghornVolumeReplica{
			Lvs: "longhorn",
		},
	}

	for _, arg := range os.Args[1:] {
		replicas = append(replicas, spdk.LonghornVolumeReplica{
			Address:  arg,
			NvmfPort: 4420,
			CommPort: 4421,
			Lvs:      "longhorn",
		})
	}

	longhornCreate := spdk.NewLonghornVolumeCreateWithReplicas(
		"test", replicas)

	client.SendCommand(longhornCreate)
	<-errChan
}
