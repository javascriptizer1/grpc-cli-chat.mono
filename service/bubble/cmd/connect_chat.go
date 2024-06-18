package cmd

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	chatv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/chat_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/bubble/internal/app"
)

type errMsg error

type chatModel struct {
	ctx context.Context
	sp  *app.ServiceProvider

	viewport    viewport.Model
	textarea    textarea.Model
	senderStyle lipgloss.Style

	chatID   string
	messages []*chatv1.Message
	err      error
	sub      chan *chatv1.Message
}

func initialChatModel(ctx context.Context, sp *app.ServiceProvider, chatID string) chatModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(60)
	ta.SetHeight(6)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(60, 10)
	vp.SetContent("Welcome to the chat room!\nType a message and press Enter to send.")

	ta.KeyMap.InsertNewline.SetEnabled(false)

	m := chatModel{
		ctx:         ctx,
		sp:          sp,
		textarea:    ta,
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		chatID:      chatID,
		sub:         make(chan *chatv1.Message),
		err:         nil,
	}

	return m
}

func (m chatModel) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.connectChat, m.waitForMessage())
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.Type {

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			message := strings.Trim(m.textarea.Value(), "")

			if message != "" {
				err := m.sp.HandlerService(m.ctx).SendMessage(m.ctx, message, m.chatID)

				if err != nil {
					m.err = err
					return m, nil
				}

				m.textarea.Reset()
				m.viewport.GotoBottom()
			}
		}

	case *chatv1.Message:
		m.messages = append(m.messages, msg)

		m.viewport.SetContent(strings.Join(m.formatMessages(), "\n"))
		m.viewport.GotoBottom()

		return m, m.waitForMessage()

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 5

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m chatModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func (m chatModel) connectChat() tea.Msg {
	stream, err := m.sp.HandlerService(m.ctx).ConnectChat(m.ctx, m.chatID)

	if err != nil {
		m.err = err
	}

	go m.receiveMessages(stream)

	return nil
}

func (m *chatModel) receiveMessages(stream chatv1.ChatService_ConnectChatClient) {
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			m.err = err
			return
		}

		m.sub <- msg
	}
}

func (m chatModel) waitForMessage() tea.Cmd {
	return func() tea.Msg {
		select {
		case msg := <-m.sub:
			return msg
		}
	}
}

func (m chatModel) formatMessages() []string {
	formatted := make([]string, len(m.messages))

	for i, msg := range m.messages {
		formatted[i] = fmt.Sprintf("[%s] %s: %s",
			msg.CreatedAt.AsTime().Format("15:04:05"),
			m.senderStyle.Render(msg.Sender.Name),
			msg.Text,
		)
	}

	return formatted
}
