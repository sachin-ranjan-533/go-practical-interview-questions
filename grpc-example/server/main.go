package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	pb "grpc-example/gen/proto/userpb"

	"google.golang.org/grpc"
)

// userServer implements the UserService gRPC server.
type userServer struct {
	pb.UnimplementedUserServiceServer
}

// GetUser handles unary RPC requests.
func (s *userServer) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log.Printf("[Unary] GetUser called with ID=%d", req.Id)
	return &pb.UserResponse{
		Id:    req.Id,
		Name:  "John Doe",
		Email: "john@example.com",
	}, nil
}

// ListUsers handles server-side streaming RPC.
func (s *userServer) ListUsers(req *pb.UserRequest, stream pb.UserService_ListUsersServer) error {
	users := []pb.UserResponse{
		{Id: 1, Name: "Alice", Email: "alice@example.com"},
		{Id: 2, Name: "Bob", Email: "bob@example.com"},
		{Id: 3, Name: "Charlie", Email: "charlie@example.com"},
	}

	for _, user := range users {
		if err := stream.Send(&user); err != nil {
			return err
		}
		log.Printf("[Server Streaming] Sent user: ID=%d Name=%s", user.Id, user.Name)
		time.Sleep(time.Second) // simulate delay
	}
	return nil
}

// UploadActions handles client-side streaming RPC.
func (s *userServer) UploadActions(stream pb.UserService_UploadActionsServer) error {
	var count int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("[Client Streaming] All actions received. Total: %d", count)
			return stream.SendAndClose(&pb.ActionSummary{Count: count})
		}
		if err != nil {
			return err
		}
		count++
		log.Printf("[Client Streaming] Received action: %s", req.Action)
	}
}

// Chat handles bidirectional streaming RPC.
func (s *userServer) Chat(stream pb.UserService_ChatServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("[Bidirectional] Client closed the stream")
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("[Bidirectional] Client: %s", msg.Text)

		if err := stream.Send(&pb.ChatMessage{Text: "Ack: " + msg.Text}); err != nil {
			return err
		}
		log.Printf("[Bidirectional] Sent acknowledgment for: %s", msg.Text)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &userServer{})

	log.Println("User service listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
