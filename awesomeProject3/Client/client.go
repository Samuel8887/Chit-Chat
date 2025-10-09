package main

import (
	proto "ITUserver/gRPC"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // for security
)

func main() {
	conn, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	client := proto.NewITUdatabaseClient(conn)
	students, err := client.GetStudents(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatalf("could not get students: %v", err)
	}

	for _, student := range students.Students {
		log.Printf("Student: %v", student)
	}
}
