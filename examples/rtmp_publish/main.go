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
	// url := "rtsp://admin:123456@192.168.51.113:8555/stream/profile/0"
	url := "rtsp://admin:123456@192.168.51.110:8555/stream/profile/0"
	file, err := avutil.Open(ctx, url)
	if err != nil {
		panic(err)
	}

	// conn, err := rtmp.DialTimeout("rtmp://localhost:1935/app/publish", 120*time.Second)
	conn, err := rtmp.Dial("rtmp://localhost:1935/app/publish")
	// // conn, _ := avutil.Create("rtmp://localhost:1936/app/publish")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %v \n", err)
		return
	}

	demuxer := &pktque.FilterDemuxer{Demuxer: file, Filter: &pktque.Walltime{}}
	// demuxer.Streams()

	avutil.CopyFile(conn, demuxer)

	fmt.Println("Wait for finish")
	time.Sleep(100 * time.Second)

	file.Close()
	conn.Close()
}
