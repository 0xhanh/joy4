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

func (s *Stream) IsStreamAvail() bool {
	switch s.Sdp.Type.String() {
	case "H264":
		return (s.Sdp.Config != nil) || (s.sps != nil && s.pps != nil)
	default:
		return true
	}
}

func (s *Stream) IsSpsAndPpsReady() bool {
	if s.sps == nil || s.pps == nil {
		return false
	}
	return true
}
