package dto

type RegisterInputDto struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            uint16
}
