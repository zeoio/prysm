package eventsv1

import (
	"time"

	gwpb "github.com/grpc-ecosystem/grpc-gateway/v2/proto/gateway"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

func (s *Server) StreamEvents(
	_ *ethpb.streame, stream pb.Events_StreamEventsServer,
) error {
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			data, err := anypb.New(req)
			if err != nil {
				return err
			}
			if err := stream.Send(&gwpb.EventSource{
				Event: "pong",
				Data:  data,
			}); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return errors.New("context canceled")
		}
	}
}
