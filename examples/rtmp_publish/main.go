package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kerberos-io/joy4/av/avutil"
	"github.com/kerberos-io/joy4/av/pktque"
	"github.com/kerberos-io/joy4/format"
	"github.com/kerberos-io/joy4/format/rtmp"
)

func init() {
	format.RegisterAll()
}

// as same as: ffmpeg -re -i projectindex.flv -c copy -f flv rtmp://localhost:1936/app/publish

func main() {
	ctx := context.Background()
	// file, _ := avutil.Open(ctx, "/home/hanh/elcom/vms/joy4/sample.flv")
	url := "rtsp://admin:123456@192.168.51.113:8555/stream/profile/0"
	file, err := avutil.Open(ctx, url)
	if err != nil {
		panic(err)
		return
	}

	conn, err := rtmp.DialTimeout("rtmp://localhost:1935/app/publish", 30*time.Second)
	// conn, _ := avutil.Create("rtmp://localhost:1936/app/publish")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %v \n", err)
		return
	}

	demuxer := &pktque.FilterDemuxer{Demuxer: file, Filter: &pktque.Walltime{}}
	avutil.CopyFile(conn, demuxer)

	time.Sleep(100 * time.Second)

	file.Close()
	conn.Close()
}
