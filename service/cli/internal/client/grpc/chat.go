package client

import (
	"context"

	chatv1 "github.com/javascriptizer1/grpc-cli-chat.mono/pkg/grpc/chat_v1"
	"github.com/javascriptizer1/grpc-cli-chat.mono/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/cli/internal/domain"
)

type ChatClient struct {
	client chatv1.ChatServiceClient
}

func NewChatClient(client chatv1.ChatServiceClient) *ChatClient {
	return &ChatClient{client: client}
}

func (c *ChatClient) CreateChat(ctx context.Context, name string, userIDs []string) (id string, err error) {
	res, err := c.client.CreateChat(ctx, &chatv1.CreateChatRequest{
		Name:    name,
		UserIDs: userIDs,
	})

	if err != nil {
		return id, err
	}

	return res.GetId(), nil
}

func (c *ChatClient) ConnectChat(ctx context.Context, chatID string) (cha chatv1.ChatService_ConnectChatClient, err error) {
	return c.client.ConnectChat(ctx, &chatv1.ConnectChatRequest{
		ChatId: chatID,
	})
}

func (c *ChatClient) SendMessage(ctx context.Context, text string, chatID string) error {
	_, err := c.client.SendMessage(ctx, &chatv1.SendMessageRequest{ChatId: chatID, Text: text})

	return err
}

func (c *ChatClient) GetChatList(ctx context.Context, p *pagination.Pagination) ([]*domain.ChatListInfo, uint32, error) {
	res, err := c.client.GetChatList(ctx, &chatv1.GetChatListRequest{Limit: uint32(p.Limit()), Page: uint32(p.Page())})

	if err != nil {
		return []*domain.ChatListInfo{}, 0, err
	}

	var chats []*domain.ChatListInfo

	for _, v := range res.Chats {
		chats = append(chats, &domain.ChatListInfo{ID: v.Id, Name: v.Name})
	}

	return chats, res.GetTotal(), nil

}

func (c *ChatClient) GetChat(ctx context.Context, id string) (*domain.ChatInfo, error) {
	res, err := c.client.GetChat(ctx, &chatv1.GetChatRequest{Id: id})

	if err != nil {
		return nil, err
	}

	var users []*domain.ChatUser

	for _, v := range res.GetUsers() {
		users = append(users, &domain.ChatUser{ID: v.Id, Name: v.Name})
	}

	chat := &domain.ChatInfo{
		ID:    res.GetId(),
		Name:  res.GetName(),
		Users: users,
	}

	return chat, nil
}
