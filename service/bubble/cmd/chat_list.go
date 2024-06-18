package cmd

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/bubble/internal/app"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type chatListModel struct {
	ctx  context.Context
	sp   *app.ServiceProvider
	list list.Model
	err  error
}

func initialChatListModel(ctx context.Context, sp *app.ServiceProvider) chatListModel {
	chats, _, err := sp.HandlerService(ctx).GetChatList(ctx, pagination.New(10, 1))

	items := make([]list.Item, len(chats))

	for i, chat := range chats {
		items[i] = listItem{title: chat.Name, desc: chat.ID}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Chats"

	return chatListModel{
		ctx:  ctx,
		sp:   sp,
		list: l,
		err:  err,
	}
}

func (m chatListModel) Init() tea.Cmd {
	return nil
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
			chatModel := initialChatModel(m.ctx, m.sp, selectedChatID)
			return chatModel, chatModel.Init()
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m chatListModel) View() string {
	return docStyle.Render(m.list.View())
}
