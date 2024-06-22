package grpc

import (
	"context"
	"fmt"
	"sync"

	chatv1 "github.com/javascriptizer1/grpc-cli-chat.mono/pkg/grpc/chat_v1"
	"github.com/javascriptizer1/grpc-cli-chat.mono/pkg/helper/array"
	"github.com/javascriptizer1/grpc-cli-chat.mono/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/chat/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/chat/internal/logger"
	"go.uber.org/zap"
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
	u, err := impl.userClient.GetUserInfo(ctx)

	if err != nil {
		return nil, err
	}

	userIDs := array.Unique(append(request.GetUserIDs(), u.ID))

	id, err := impl.chatService.Create(ctx, request.GetName(), userIDs)

	if err != nil {
		return nil, err
	}

	return &chatv1.CreateChatResponse{Id: id}, nil
}

func (impl *ChatImplementation) ConnectChat(request *chatv1.ConnectChatRequest, stream chatv1.ChatService_ConnectChatServer) error {
	u, err := impl.userClient.GetUserInfo(stream.Context())

	if err != nil {
		return err
	}

	ch := make(chan *chatv1.Message, maxMessageInChannel)

	impl.mu.Lock()

	if impl.clients[request.GetChatId()] == nil {
		impl.clients[request.GetChatId()] = make(map[string]chan *chatv1.Message)
	}

	impl.clients[request.GetChatId()][u.ID] = ch

	impl.mu.Unlock()

	logger.Info(fmt.Sprintf("User %s connected to chat %s", u.ID, request.GetChatId()))

	defer func() {
		impl.mu.Lock()
		delete(impl.clients[request.GetChatId()], u.ID)

		if len(impl.clients[request.GetChatId()]) == 0 {
			delete(impl.clients, request.GetChatId())
		}

		impl.mu.Unlock()

		logger.Info(fmt.Sprintf("User %s disconnected from chat %s", u.ID, request.GetChatId()))
		close(ch)
	}()

	for {
		select {
		case msg := <-ch:
			if err := stream.Send(msg); err != nil {
				logger.Error("Error sending message to client", zap.Error(err))
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
				logger.Warn("Warning: message channel for user is full, dropping message")
			}
		}
	}

	return &emptypb.Empty{}, nil
}

func (impl *ChatImplementation) GetChatList(ctx context.Context, request *chatv1.GetChatListRequest) (*chatv1.GetChatListResponse, error) {
	u, err := impl.userClient.GetUserInfo(ctx)

	if err != nil {
		return nil, err
	}

	chats, total, err := impl.chatService.List(ctx, u.ID, *pagination.New(request.Limit, request.Page))

	if err != nil {
		return nil, err
	}

	var protoChats = make([]*chatv1.Chat, len(chats))

	for i, v := range chats {
		c := &chatv1.Chat{Id: v.ID, Name: v.Name}
		protoChats[i] = c
	}

	return &chatv1.GetChatListResponse{
			Chats: protoChats,
			Total: uint32(total),
		},
		nil
}
func (impl *ChatImplementation) GetChat(ctx context.Context, request *chatv1.GetChatRequest) (*chatv1.GetChatResponse, error) {
	chat, err := impl.chatService.OneByID(ctx, request.GetId())

	if err != nil {
		return nil, err
	}

	var userIDs []string
	var chatUsers []*chatv1.User

	for _, v := range chat.Users {
		userIDs = append(userIDs, v.ID)
	}

	if len(userIDs) != 0 {
		userInfos, _, err := impl.userClient.GetUserList(
			ctx,
			&domain.UserInfoListFilter{
				Limit:   uint32(len(chat.Users)),
				Page:    1,
				UserIDs: userIDs,
			},
		)

		if err != nil {
			return nil, err
		}

		for _, v := range userInfos {
			chatUsers = append(chatUsers, &chatv1.User{
				Id:   v.ID,
				Name: v.Name,
			})
		}
	}

	return &chatv1.GetChatResponse{
		Id:    chat.ID,
		Name:  chat.Name,
		Users: chatUsers,
	}, nil
}
