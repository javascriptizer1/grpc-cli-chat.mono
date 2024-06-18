package handler

import "context"

func (h *Handler) Login(ctx context.Context, login string, password string) (string, error) {
	refreshToken, err := h.authClient.Login(ctx, login, password)

	if err != nil {
		return "", err
	}

	accessToken, err := h.authClient.GetAccessToken(ctx, refreshToken)

	if err != nil {
		return "", err
	}

	err = h.tokenManager.SetTokens(accessToken, refreshToken)

	return refreshToken, err

}
