package main

import (
	proto "Chit_Chat/gRPC"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type ITU_databaseServer struct {
	proto.UnimplementedITUdatabaseServer
	students []string
}

func (s *ITU_databaseServer) GetStudents(ctx context.Context, in *proto.Empty) (*proto.Students, error) {
	return &proto.Students{Students: s.students}, nil
} // the receiver

func main() {
	server := &ITU_databaseServer{students: []string{}}
	server.students = append(server.students, "Alice")
	server.students = append(server.students, "Bob")
	server.students = append(server.students, "Charlie")
	server.students = append(server.students, "David")

	server.start_server()
}

func (s *ITU_databaseServer) start_server() {
	grpcserver := grpc.NewServer()
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	proto.RegisterITUdatabaseServer(grpcserver, s)

	err = grpcserver.Serve(listener)

	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
} // server
