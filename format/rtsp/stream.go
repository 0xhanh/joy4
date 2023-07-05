package rtsp

import (
	"time"

	"github.com/kerberos-io/joy4/av"
	"github.com/kerberos-io/joy4/format/rtsp/sdp"
)

type Stream struct {
	av.CodecData
	Sdp    sdp.Media
	client *Client

	// h264
	fuStarted  bool
	fuBuffer   []byte
	sps        []byte
	pps        []byte
	spsChanged bool
	ppsChanged bool

	gotpkt         bool
	pkt            av.Packet
	timestamp      uint32
	firsttimestamp uint32

	lasttime time.Duration
}

// hvd
func (s *Stream) IsStreamAvail() bool {
	return (s.Sdp.Config != nil) || (s.sps != nil && s.pps != nil)
}
