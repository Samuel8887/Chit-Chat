package main

import (
	proto "Chit_Chat/gRPC"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // for security
)

func main() {
	conn, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	client := proto.NewChit_ChatClient(conn)
	reader := bufio.NewReader(os.Stdin)
	logicalTime := 0

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("input went wrong %v", err)
		}
		input = strings.TrimSpace(input)
		if input == "joinP1" {
			joinRequest, _ := client.Join(context.Background(), &proto.JoinRequest{
				ClientId: "P1", LogicalTime: int64(logicalTime + 1),
			})
			log.Printf("joining: %v", joinRequest)
		}

		if input == "publishP1" {
			publishRequest, _ := client.Publish(context.Background(), &proto.PublishRequest{
				ClientId: "P1", LogicalTime: int64(logicalTime + 1), Content: input,
			})
			log.Printf("Publishing: %v", publishRequest)
		}
		if input == "leaveP1" {
			leaveRequest, _ := client.Leave(context.Background(), &proto.LeaveRequest{
				ClientId: "P1", LogicalTime: int64(logicalTime + 1),
			})
			log.Printf("Leaving: %v", leaveRequest)
		}
	}
	/*students, err := client.GetStudents(context.Background(), &proto.Empty{})

	for _, student := range students.Students {
		log.Printf("Student: %v", student)
	}*/
}
