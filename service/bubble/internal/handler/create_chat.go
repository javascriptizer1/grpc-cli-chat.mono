package handler

import "context"

func (h *Handler) CreateChat(ctx context.Context, name string, emails []string) (string, error) {
	id, err := h.chatClient.CreateChat(ctx, name, emails)

	if err != nil {
		return "", err
	}

	return id, nil
}
