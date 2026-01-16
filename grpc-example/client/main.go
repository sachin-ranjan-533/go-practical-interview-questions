package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "grpc-example/gen/proto/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := doUnaryRPC(ctx, client, 1); err != nil {
		log.Fatalf("Unary RPC failed: %v", err)
	}

	if err := doServerStreamingRPC(ctx, client); err != nil {
		log.Fatalf("Server streaming RPC failed: %v", err)
	}

	if err := doClientStreamingRPC(ctx, client); err != nil {
		log.Fatalf("Client streaming RPC failed: %v", err)
	}

	if err := doBidirectionalStreamingRPC(ctx, client); err != nil {
		log.Fatalf("Bidirectional streaming RPC failed: %v", err)
	}
}

// Unary RPC
func doUnaryRPC(ctx context.Context, client pb.UserServiceClient, userID int64) error {
	userResp, err := client.GetUser(ctx, &pb.UserRequest{Id: userID})
	if err != nil {
		return err
	}

	log.Printf("[Unary] User received: ID=%d Name=%s Email=%s",
		userResp.Id, userResp.Name, userResp.Email)
	return nil
}

// Server streaming RPC
func doServerStreamingRPC(ctx context.Context, client pb.UserServiceClient) error {
	stream, err := client.ListUsers(ctx, &pb.UserRequest{})
	if err != nil {
		return err
	}

	for {
		user, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("[Server Streaming] User: ID=%d Name=%s Email=%s", user.Id, user.Name, user.Email)
	}
	return nil
}

// Client streaming RPC
func doClientStreamingRPC(ctx context.Context, client pb.UserServiceClient) error {
	stream, err := client.UploadActions(ctx)
	if err != nil {
		return err
	}

	actions := []string{"login", "view_page", "logout"}
	for _, action := range actions {
		if err := stream.Send(&pb.ActionRequest{Action: action}); err != nil {
			return err
		}
		log.Printf("[Client Streaming] Sent action: %s", action)
	}

	summary, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("[Client Streaming] Action summary received: Count=%d", summary.Count)
	return nil
}

// Bidirectional streaming RPC
func doBidirectionalStreamingRPC(ctx context.Context, client pb.UserServiceClient) error {
	stream, err := client.Chat(ctx)
	if err != nil {
		return err
	}

	// Sending messages in a goroutine
	go func() {
		messages := []string{"Hello Server!", "How are you?"}
		for _, msg := range messages {
			if err := stream.Send(&pb.ChatMessage{Text: msg}); err != nil {
				log.Printf("Failed to send message: %v", err)
				return
			}
			log.Printf("[Bidirectional] Sent: %s", msg)
		}
		stream.CloseSend()
	}()

	// Receiving messages
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("[Bidirectional] Server: %s", msg.Text)
	}

	return nil
}
