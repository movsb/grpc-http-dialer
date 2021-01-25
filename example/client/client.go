package main

import (
	"context"
	"fmt"

	"github.com/movsb/grpc-http-dialer"
	ping "github.com/movsb/grpc-http-dialer/example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	for i := 0; i < 64; i++ {
		conn, err := grpc.Dial(
			`localhost:43210`,
			grpc.WithInsecure(),
			grpc.WithContextDialer(grpchttpdialer.Dialer(`localhost:8080`)),
		)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		client := ping.NewPingServiceClient(conn)
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(`aaaaaaaaaaaaaaaaaaaaaaaaaaa`, `b`))
		fmt.Println(client.Ping(ctx, &ping.PingRequest{Ping: `test`}))
	}
}
