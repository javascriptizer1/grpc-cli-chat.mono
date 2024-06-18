package cmd

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/bubble/internal/app"
)

type tab int

const (
	chatsTab tab = iota
	peopleTab
)

type tabsModel struct {
	ctx        context.Context
	sp         *app.ServiceProvider
	currentTab tab
	chatList   chatListModel
	peopleList userListModel
	list       list.Model
}

func InitialTabsModel(ctx context.Context, sp *app.ServiceProvider) tabsModel {
	items := []list.Item{
		listItem{title: "Chats", desc: "List of chats"},
		listItem{title: "People", desc: "List of people"},
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Menu"

	return tabsModel{
		ctx:        ctx,
		sp:         sp,
		currentTab: chatsTab,
		chatList:   chatListModel{},
		peopleList: userListModel{},
		list:       l,
	}
}

func (m tabsModel) Init() tea.Cmd {
	return nil
}

func (m tabsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.currentTab == chatsTab {
				m.currentTab = peopleTab
			} else {
				m.currentTab = chatsTab
			}
		case "enter":
			switch m.list.SelectedItem().(listItem).title {
			case "Chats":
				chatListModel := initialChatListModel(m.ctx, m.sp)
				return chatListModel, chatListModel.Init()
			case "People":
				userListModel := initialUserListModel(m.ctx, m.sp)
				return userListModel, userListModel.Init()
			case "Create Chat":
				return initialCreateChatModel(m.ctx, m.sp), nil
			}
		}
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m tabsModel) View() string {
	return m.list.View()
}

type listItem struct {
	title string
	desc  string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title }
