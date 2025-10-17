package main

import (
	proto "Chit_Chat/gRPC"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

type client struct {
	id     string
	stream proto.Chit_Chat_JoinServer
}

type Server struct {
	grpcServer *grpc.Server
	proto.UnimplementedChit_ChatServer

	clients   []client
	clientsMu sync.Mutex
}

func main() {
	server := &Server{}
	server.StartServer()
}

func (s *Server) StartServer() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("[%s] Server starting on :8000", time.Now().Format(time.RFC3339))
	s.grpcServer = grpc.NewServer()
	proto.RegisterChit_ChatServer(s.grpcServer, s)

	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) StopServer() {
	if s.grpcServer != nil {
		log.Printf("[%s] Server stopping on :8000", time.Now().Format(time.RFC3339))
		s.grpcServer.Stop()
	}
}

func (s *Server) Join(req *proto.JoinRequest, stream proto.Chit_Chat_JoinServer) error {
	log.Printf("[%s] User %s joined", time.Now().Format(time.RFC3339), req.ClientId)

	s.clientsMu.Lock()
	s.clients = append(s.clients, client{id: req.ClientId, stream: stream})
	s.clientsMu.Unlock()

	stream.Send(&proto.ChatMessage{
		From:    "server",
		Content: "Welcome " + req.ClientId + "!",
	})

	select {}
}

func (s *Server) Publish(ctx context.Context, req *proto.PublishRequest) (*proto.Ack, error) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	msg := &proto.ChatMessage{
		From:    req.ClientId,
		Content: req.Content,
	}

	for i := 0; i < len(s.clients); i++ {
		if err := s.clients[i].stream.Send(msg); err != nil {
			log.Printf("failed to send to %s: %v", s.clients[i].id, err)
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			i--
		}
	}

	return &proto.Ack{Success: true}, nil
}

func (s *Server) Leave(ctx context.Context, req *proto.LeaveRequest) (*proto.Ack, error) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	for i, c := range s.clients {
		if c.id == req.ClientId {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			log.Printf("[%s] User %s left", time.Now().Format(time.RFC3339), req.ClientId)
			break
		}
	}

	return &proto.Ack{Success: true}, nil
}
