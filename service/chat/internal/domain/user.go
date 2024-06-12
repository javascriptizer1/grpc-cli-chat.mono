package domain

type UserInfo struct {
	ID    string
	Name  string
	Email string
	Role  uint16
}

func NewUserInfo(ID string, name string, email string, role uint16) *UserInfo {
	return &UserInfo{
		ID:    ID,
		Name:  name,
		Email: email,
		Role:  role,
	}
}
