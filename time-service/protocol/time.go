package protocol

import (
	"context"
	"fmt"
	"io"
	"time"

	messages "deshev.com/libp2p-time-service/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	ggio "github.com/gogo/protobuf/io"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

const (
	timeRequest  = "/deshev.com/time-request/1.0.0"
	timeResponse = "/deshev.com/time-response/1.0.0"
)

type TimeProtocol struct {
	node host.Host
}

func NewProtocol(host host.Host) *TimeProtocol {
	proto := &TimeProtocol{node: host}
	host.SetStreamHandler(timeRequest, proto.timeRequestHandler)
	host.SetStreamHandler(timeResponse, proto.timeResponseHandler)
	return proto
}

func (p *TimeProtocol) StartRequest(ctx context.Context, peer peerstore.ID) error {
	s, err := p.node.NewStream(ctx, peer, timeRequest)
	if err != nil {
		return fmt.Errorf("request_stream: %w", err)
	}
	defer s.Close()

	msg := &messages.TimeRequest{
		Greeting: "What time is it?",
	}
	writer := ggio.NewFullWriter(s)
	err = writer.WriteMsg(msg)
	if err != nil {
		s.Reset()
		return fmt.Errorf("request_write: %w", err)
	}

	return nil
}

func (p *TimeProtocol) timeRequestHandler(stream network.Stream) {
	remotePeer := stream.Conn().RemotePeer()
	data, err := io.ReadAll(stream)
	if err != nil {
		stream.Reset()
		fmt.Printf("read_error: %v", err)
	}
	stream.Close()

	msg := messages.TimeRequest{}
	err = proto.Unmarshal(data, &msg)
	if err != nil {
		fmt.Printf("unmarshal_error: %v", err)
	}

	fmt.Printf("time request(%s): %s\n", remotePeer, msg.GetGreeting())

	s, err := p.node.NewStream(context.Background(), remotePeer, timeResponse)
	defer s.Close()

	resp := &messages.TimeResponse{
		ServerTime: timestamppb.New(time.Now().UTC()),
	}
	writer := ggio.NewFullWriter(s)
	err = writer.WriteMsg(resp)
	if err != nil {
		s.Reset()
		fmt.Printf("response_write: %v\n", err)
	}
}

func (p *TimeProtocol) timeResponseHandler(stream network.Stream) {
	remotePeer := stream.Conn().RemotePeer()
	data, err := io.ReadAll(stream)
	if err != nil {
		stream.Reset()
		fmt.Printf("read_error: %v", err)
	}
	stream.Close()

	msg := messages.TimeResponse{}
	err = proto.Unmarshal(data, &msg)
	if err != nil {
		fmt.Printf("unmarshal_error: %v", err)
	}

	fmt.Printf("time reponse(%s): %s\n", remotePeer, msg.GetServerTime().AsTime())
}
