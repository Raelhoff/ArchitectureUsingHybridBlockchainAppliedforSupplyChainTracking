package main

import (
	"fmt"
	pb "github.com/karankumarshreds/GoProto/protofiles"
	"google.golang.org/protobuf/proto"
)

func main() {

	// using the profo created struct
	p := &pb.Person{
		Id: 1234,
		Name: "John Doe",
		Email: "test@test.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "555-444", Type: pb.Person_HOME},
		},
	}

	// Serializing the struct and assigning it to body
	body, _ := proto.Marshal(p)

	// De-serializing the body and saving it to p1 for testing
	p1 := &pb.Person{}
	_ = proto.Unmarshal(body, p1)

	fmt.Println("Original struct loaded from proto file:", p)
	fmt.Println("Marshalled proto data: ", body)
	fmt.Println("Unmarshalled struct: ", p1)

}
