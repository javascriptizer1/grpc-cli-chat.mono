package domain

type ChatUser struct {
	ID   string
	Name string
}

type ChatInfo struct {
	ID    string
	Name  string
	Users []*ChatUser
}

type ChatListInfo struct {
	ID   string
	Name string
}
