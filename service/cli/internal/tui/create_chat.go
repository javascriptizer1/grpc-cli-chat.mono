package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/domain"
)

var (
	focusedButtonCreateChat = focusedStyle.Render("[ Create ]")
	blurredButtonCreateChat = fmt.Sprintf("[ %s ]", blurredStyle.Render("Create"))
)

type createChatModel struct {
	ctx             context.Context
	sp              *app.ServiceProvider
	name            textinput.Model
	table           table.Model
	selectedUserIDs map[string]string
	cursorMode      cursor.Mode
	focusedElement  int
	width           int
	height          int
	err             error
}

func InitialCreateChatModel(ctx context.Context, sp *app.ServiceProvider, width int, height int) createChatModel {
	name := initTextInput()
	columns := initTableColumns()
	users, err := fetchUsers(ctx, sp)
	rows := createTableRows(users)

	t := initTable(columns, rows)

	return createChatModel{
		ctx:             ctx,
		sp:              sp,
		name:            name,
		table:           t,
		selectedUserIDs: make(map[string]string),
		focusedElement:  0,
		width:           width,
		height:          height,
		err:             err,
	}
}

func (m createChatModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, func() tea.Msg {
		return tea.WindowSizeMsg{
			Width:  m.width,
			Height: m.height,
		}
	})
}

func (m createChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlR:
			m.cursorMode++

			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}

			return m, m.name.Cursor.SetMode(m.cursorMode)

		case tea.KeyShiftTab:
			chatListModel := InitialChatListModel(m.ctx, m.sp, m.width, m.height)
			return chatListModel, chatListModel.Init()
		case tea.KeyTab:
			m.toggleFocus()
		case tea.KeyEnter:
			if m.focusedElement == 1 {
				m.toggleSelectedUser()
			} else if m.focusedElement == 2 {
				selectedUserIDs := m.getSelectedUserIDs()
				_, err := m.sp.HandlerService(m.ctx).CreateChat(m.ctx, m.name.Value(), selectedUserIDs)

				if err != nil {
					m.err = err
					return m, nil
				}

				chatListModel := InitialChatListModel(m.ctx, m.sp, m.width, m.height)
				return chatListModel, chatListModel.Init()
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	if m.focusedElement == 0 {
		m.name, cmd = m.name.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.focusedElement == 1 {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m createChatModel) View() string {
	var view strings.Builder

	view.WriteString("Create a new chat\n\n")
	view.WriteString(m.name.View())
	view.WriteRune('\n')

	if len(m.selectedUserIDs) > 0 {
		view.WriteString("\nSelected users: ")
		view.WriteString(m.getSelectedUsernames())
		view.WriteRune('\n')
	} else {
		view.WriteString("\n\n")
	}

	view.WriteRune('\n')
	view.WriteString(m.table.View())

	button := &blurredButtonCreateChat

	if m.focusedElement == 2 {
		button = &focusedButtonCreateChat
	}

	fmt.Fprintf(&view, "\n%s\n\n", *button)

	if m.err != nil {
		view.WriteString(fmt.Sprintf("Error: %v\n", m.err))
	}

	view.WriteRune('\n')
	view.WriteString(helpStyle.Render("cursor mode is "))
	view.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	view.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	view.WriteString("\n\nPress Tab to switch focus, Shift+Tab to return chat list or Esc/Ctrl+C to exit.")

	content := addMargin(view.String())

	return content
}

func (m *createChatModel) toggleFocus() {
	m.focusedElement++

	if m.focusedElement > 2 {
		m.focusedElement = 0
		m.name.Focus()
		m.name.PromptStyle = focusedStyle
		m.name.TextStyle = focusedStyle
		m.table.Blur()
	} else if m.focusedElement < 0 {
		m.focusedElement = 2
		m.name.Blur()
		m.name.PromptStyle = noStyle
		m.name.TextStyle = noStyle
		m.table.Blur()
	} else if m.focusedElement == 1 {
		m.name.Blur()
		m.name.PromptStyle = noStyle
		m.name.TextStyle = noStyle
		m.table.Focus()
	}
}

func (m *createChatModel) toggleSelectedUser() {
	selectedRow := m.table.SelectedRow()
	userID := selectedRow[0]
	username := selectedRow[1]

	if _, exists := m.selectedUserIDs[userID]; exists {
		delete(m.selectedUserIDs, userID)
	} else {
		m.selectedUserIDs[userID] = username
	}
}

func (m *createChatModel) getSelectedUserIDs() []string {
	var userIDs []string

	for userID := range m.selectedUserIDs {
		userIDs = append(userIDs, userID)
	}

	return userIDs
}

func (m *createChatModel) getSelectedUsernames() string {
	var usernames []string

	for _, username := range m.selectedUserIDs {
		usernames = append(usernames, username)
	}

	return strings.Join(usernames, ", ")
}

func initTextInput() textinput.Model {
	name := textinput.New()
	name.Placeholder = "Chat Name"
	name.Cursor.Style = cursorStyle
	name.PromptStyle = focusedStyle
	name.TextStyle = focusedStyle
	name.Focus()

	return name
}

func initTableColumns() []table.Column {
	return []table.Column{
		{Title: "IDX", Width: 10},
		{Title: "Username", Width: 20},
		{Title: "Email", Width: 31},
	}
}

func fetchUsers(ctx context.Context, sp *app.ServiceProvider) ([]*domain.UserInfo, error) {
	users, _, err := sp.HandlerService(ctx).GetUserList(ctx, &domain.UserListOption{
		Pagination: *pagination.New(10, 1),
		UserIDs:    []string{},
	})

	return users, err
}

func createTableRows(users []*domain.UserInfo) []table.Row {
	rows := make([]table.Row, len(users))

	for i, user := range users {
		rows[i] = []string{strconv.Itoa(i + 1), user.Name, user.Email}
	}

	return rows
}

func initTable(columns []table.Column, rows []table.Row) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}
