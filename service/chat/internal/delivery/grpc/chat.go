package grpc

import (
	"context"
	"log"
	"sync"

	chatv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const maxMessageInChannel = 100

type ChatImplementation struct {
	chatv1.UnimplementedChatServiceServer

	chatService ChatService
	userClient  UserClient
	mu          sync.RWMutex
	clients     map[string]map[string]chan *chatv1.Message
}

func NewGrpcChatImplementation(chatService ChatService, userClient UserClient) *ChatImplementation {
	return &ChatImplementation{
		chatService: chatService,
		userClient:  userClient,
		clients:     make(map[string]map[string]chan *chatv1.Message),
	}
}

func (impl *ChatImplementation) CreateChat(ctx context.Context, request *chatv1.CreateChatRequest) (*chatv1.CreateChatResponse, error) {
	id, err := impl.chatService.Create(ctx, request.GetEmails())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &chatv1.CreateChatResponse{Id: id}, nil
}

func (impl *ChatImplementation) ConnectChat(request *chatv1.ConnectChatRequest, stream chatv1.ChatService_ConnectChatServer) error {
	ch := make(chan *chatv1.Message, maxMessageInChannel)

	impl.mu.Lock()

	if impl.clients[request.GetChatId()] == nil {
		impl.clients[request.GetChatId()] = make(map[string]chan *chatv1.Message)
	}

	impl.clients[request.GetChatId()][request.GetUserId()] = ch

	impl.mu.Unlock()

	log.Printf("User %s connected to chat %s", request.GetUserId(), request.GetChatId())

	defer func() {
		impl.mu.Lock()
		delete(impl.clients[request.GetChatId()], request.GetUserId())

		if len(impl.clients[request.GetChatId()]) == 0 {
			delete(impl.clients, request.GetChatId())
		}

		impl.mu.Unlock()

		log.Printf("User %s disconnected from chat %s", request.GetUserId(), request.GetChatId())
		close(ch)
	}()

	for {
		select {
		case msg := <-ch:
			if err := stream.Send(msg); err != nil {
				log.Printf("Error sending message to client: %v", err)
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (impl *ChatImplementation) SendMessage(ctx context.Context, request *chatv1.SendMessageRequest) (*emptypb.Empty, error) {
	c, err := impl.chatService.OneByID(ctx, request.GetChatId())

	if c == nil || err != nil {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	u, err := impl.userClient.GetUserInfo(ctx)

	if err != nil {
		return nil, err
	}

	m, err := impl.chatService.CreateMessage(ctx, request.GetText(), request.GetChatId(), *u)

	if err != nil {
		return nil, err
	}

	msgProto := &chatv1.Message{
		Id:        m.ID,
		Sender:    &chatv1.User{Id: m.Sender.ID, Name: m.Sender.Name},
		Text:      m.Text,
		CreatedAt: timestamppb.New(m.CreatedAt),
	}

	impl.mu.RLock()
	streams, ok := impl.clients[request.GetChatId()]
	impl.mu.RUnlock()

	if ok {
		for _, ch := range streams {
			select {
			case ch <- msgProto:
			default:
				log.Printf("Warning: message channel for user is full, dropping message")
			}
		}
	}

	return &emptypb.Empty{}, nil
}
