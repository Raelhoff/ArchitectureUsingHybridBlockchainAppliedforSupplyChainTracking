package initializers

import (
	"fmt"

	"fiber-mongo-api/packge"
	pb "github.com/Raelhoff/gRPC_GO/protofiles"
	"google.golang.org/grpc"
)

// grpc server address
// const address = "127.0.0.1:8000"
const address = "192.168.0.192:8017"

func ConnectGRPC() *grpc.ClientConn {
	client, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Error while making connection")
	}

	fmt.Println("Connected to GRPC")
	return client
}

// Client instance
var conn = ConnectGRPC()

func ReturnClientGRPC() *grpc.ClientConn {
	return conn
}

// getting database collections
func GetProtofileLoraTransaction() pb.LoraTransactionClient {
	c := pb.NewLoraTransactionClient(conn)
	if c == nil {
		fmt.Println("Error GetProtofileLoraTransaction")
	}
	return c
}

// getting database collections
func GetProtofileAlert() packge.AlertClient {
	c := packge.NewAlertClient(conn)
	if c == nil {
		fmt.Println("Error GetProtofileLoraTransaction")
	}
	return c
}
