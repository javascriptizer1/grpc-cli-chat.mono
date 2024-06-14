package handler

type Handler struct {
	authClient   AuthClient
	userClient   UserClient
	chatClient   ChatClient
	tokenManager TokenManager
}

func New(authClient AuthClient, userClient UserClient, chatClient ChatClient, tokenManager TokenManager) *Handler {
	return &Handler{
		authClient:   authClient,
		userClient:   userClient,
		chatClient:   chatClient,
		tokenManager: tokenManager,
	}
}
