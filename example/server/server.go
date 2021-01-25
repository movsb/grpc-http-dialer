package main

import (
	"context"
	"net"
	"net/http"

	"github.com/movsb/grpc-http-dialer"
	ping "github.com/movsb/grpc-http-dialer/example/proto"
	"google.golang.org/grpc"
)

// Service ...
type Service struct {
}

// Ping ...
func (s *Service) Ping(ctx context.Context, in *ping.PingRequest) (*ping.PingResponse, error) {
	return &ping.PingResponse{
		Pong: in.Ping + ` & pong`,
	}, nil
}

func main() {
	s := grpc.NewServer()
	l, err := net.Listen(`tcp4`, `localhost:43210`)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	go func() {
		http.Handle(grpchttpdialer.ProxyPath, grpchttpdialer.Handler())
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()

	ping.RegisterPingServiceServer(s, &Service{})
	if err := s.Serve(l); err != nil {
		panic(err)
	}
}
