package client

import (
	"context"

	chatv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/chat_v1"
)

type ChatClient struct {
	client chatv1.ChatServiceClient
}

func NewChatClient(client chatv1.ChatServiceClient) *ChatClient {
	return &ChatClient{client: client}
}

func (c *ChatClient) CreateChat(ctx context.Context, emails []string) (id string, err error) {
	res, err := c.client.CreateChat(ctx, &chatv1.CreateChatRequest{
		Emails: emails,
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
