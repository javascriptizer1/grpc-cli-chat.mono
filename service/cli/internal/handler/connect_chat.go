package handler

import (
	"bufio"
	"context"

	"io"
	"os"

	chatv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/chat_v1"
	colog "github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/util"
)

func (h *Handler) ConnectChat(ctx context.Context, chatID string) {
	stream, err := h.chatClient.ConnectChat(ctx, chatID)

	if err != nil {
		colog.Error("could not connect to chat: %v", err)
	}

	colog.Success("happy chatting!")

	go h.receiveMessages(stream)

	h.sendMessages(ctx, chatID)
}

func (h *Handler) receiveMessages(stream chatv1.ChatService_ConnectChatClient) {
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			colog.Fatal("Error receiving message: %v", err)
		}

		colog.Info(
			"[%s] %s: %s",
			msg.CreatedAt.AsTime().Format("15:04:05"),
			msg.Sender.Name,
			msg.Text,
		)
	}
}

func (h *Handler) sendMessages(ctx context.Context, chatID string) {
	for {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			line := scanner.Text()

			if len(line) == 0 {
				break
			}

			line = line + "\n"

			if err := h.chatClient.SendMessage(ctx, line, chatID); err != nil {
				colog.Warn("failed to send message: %v", err)
			}
		}

		if err := scanner.Err(); err != nil {
			colog.Warn("failed to scan message: %v", err)
		}
	}
}
