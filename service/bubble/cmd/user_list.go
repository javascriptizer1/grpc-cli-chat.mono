package cmd

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/bubble/internal/app"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/bubble/internal/domain"
)

type userListModel struct {
	ctx  context.Context
	sp   *app.ServiceProvider
	list list.Model
	err  error
}

func initialUserListModel(ctx context.Context, sp *app.ServiceProvider) userListModel {
	people, _, err := sp.HandlerService(ctx).GetUserList(ctx, &domain.UserListOption{
		Pagination: *pagination.New(10, 1),
		UserIDs:    []string{},
	})

	items := make([]list.Item, len(people))
	for i, person := range people {
		items[i] = listItem{title: person.Name, desc: person.Email}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "People"

	return userListModel{
		ctx:  ctx,
		sp:   sp,
		list: l,
		err:  err,
	}
}

func (m userListModel) Init() tea.Cmd {
	return nil
}

func (m userListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m userListModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	return m.list.View()
}
