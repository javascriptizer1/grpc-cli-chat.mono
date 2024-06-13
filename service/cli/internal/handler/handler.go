package handler

type Handler struct {
	authClient   AuthClient
	chatClient   ChatClient
	tokenManager TokenManager
}

func New(authClient AuthClient, chatClient ChatClient, tokenManager TokenManager) *Handler {
	return &Handler{
		authClient:   authClient,
		chatClient:   chatClient,
		tokenManager: tokenManager,
	}
}
