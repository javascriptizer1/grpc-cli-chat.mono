package interceptor

type TokenManager interface {
	SetTokens(accessToken, refreshToken string) error
	AccessToken() string
	RefreshToken() string
}
