package handler

import "context"

func (h *Handler) CreateChat(ctx context.Context, emails []string) (string, error) {
	id, err := h.chatClient.CreateChat(ctx, emails)

	if err != nil {
		return "", err
	}

	return id, nil
}
