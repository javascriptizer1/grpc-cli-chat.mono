package handler

type Handler struct {
	authClient AuthClient
	chatClient ChatClient
}

func New(authClient AuthClient, chatClient ChatClient) *Handler {
	return &Handler{authClient: authClient, chatClient: chatClient}
}
