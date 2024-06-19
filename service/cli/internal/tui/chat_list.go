package tui

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/javascriptizer1/grpc-cli-chat.mono/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/cli/internal/app"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/cli/internal/domain"
)

var docStyle = lipgloss.NewStyle().Margin(3, 2)

type chatListModel struct {
	ctx    context.Context
	sp     *app.ServiceProvider
	list   list.Model
	width  int
	height int
	err    error
}

func InitialChatListModel(ctx context.Context, sp *app.ServiceProvider, width int, height int) chatListModel {
	chats, _, err := sp.HandlerService(ctx).GetChatList(ctx, pagination.New(10, 1))
	items := createListItems(chats)

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Chats"

	return chatListModel{
		ctx:    ctx,
		sp:     sp,
		list:   l,
		width:  width,
		height: height,
		err:    err,
	}
}

func (m chatListModel) Init() tea.Cmd {
	return func() tea.Msg {
		return tea.WindowSizeMsg{
			Width:  m.width,
			Height: m.height,
		}
	}
}

func (m chatListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			selectedChatID := m.list.SelectedItem().(listItem).Description()
			if selectedChatID != "" {
				chatModel := InitialConnectChatModel(m.ctx, m.sp, selectedChatID)
				return chatModel, chatModel.Init()
			}

		case tea.KeyTab:
			createChatModel := InitialCreateChatModel(m.ctx, m.sp, m.width, m.height)
			return createChatModel, createChatModel.Init()
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m chatListModel) View() string {
	var b strings.Builder

	b.WriteString(docStyle.Render(m.list.View()))

	if len(m.list.Items()) == 0 {
		b.WriteString("\nYou don't have any chats. Create new")
	}

	b.WriteString(helpStyle.Render("\nPress Tab to switch to create new chat"))

	return b.String()
}

func createListItems(chats []*domain.ChatListInfo) []list.Item {
	items := make([]list.Item, len(chats))

	for i, chat := range chats {
		items[i] = listItem{title: chat.Name, desc: chat.ID}
	}

	return items
}
