package helper

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/v2fly/v2ray-core/v5/app/stats/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func connect() command.StatsServiceClient {
	conn, err := grpc.Dial("127.0.0.1:5555", grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	return command.NewStatsServiceClient(conn)
}

func GetUserStats(name string) int64 {
	resp, err := connect().GetStats(context.Background(), &command.GetStatsRequest{
		Name:   fmt.Sprintf("user>>>%s>>>traffic>>>downlink", name),
		Reset_: true,
	})
	if err != nil {
		fmt.Println(err)
	}

	return resp.Stat.Value
}
