package tui

import (
	"context"
	"fmt"
	"hash/fnv"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	chatv1 "github.com/javascriptizer1/grpc-cli-chat.backend/pkg/grpc/chat_v1"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/internal/app"
)

type errMsg error

type chatModel struct {
	ctx      context.Context
	sp       *app.ServiceProvider
	viewport viewport.Model
	textarea textarea.Model
	chatID   string
	messages []*chatv1.Message
	err      error
	sub      chan *chatv1.Message
}

func InitialConnectChatModel(ctx context.Context, sp *app.ServiceProvider, chatID string) chatModel {
	ta := initializeTextarea()
	vp := initializeViewport()

	return chatModel{
		ctx:      ctx,
		sp:       sp,
		textarea: ta,
		viewport: vp,
		chatID:   chatID,
		sub:      make(chan *chatv1.Message),
		err:      nil,
	}
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
	view := fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"

	content := addMargin(view)

	return content
}

func (m chatModel) connectChat() tea.Msg {
	stream, err := m.sp.HandlerService(m.ctx).ConnectChat(m.ctx, m.chatID)

	if err != nil {
		m.err = err
		return nil
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
		return <-m.sub
	}
}

func (m chatModel) formatMessages() []string {
	formatted := make([]string, len(m.messages))

	for i, msg := range m.messages {
		color := pickColorForSender(msg.Sender.Name)
		senderStyle := lipgloss.NewStyle().Foreground(color)

		formatted[i] = fmt.Sprintf("[%s] %s: %s",
			msg.CreatedAt.AsTime().Format("15:04:05"),
			senderStyle.Render(msg.Sender.Name),
			msg.Text,
		)
	}

	return formatted
}

func initializeTextarea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 280
	ta.SetWidth(60)
	ta.SetHeight(6)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return ta
}

func initializeViewport() viewport.Model {
	vp := viewport.New(60, 10)
	vp.SetContent("Welcome to the chat room!\nType a message and press Enter to send.")

	return vp
}

func pickColorForSender(name string) lipgloss.Color {
	hash := fnv.New32a()
	hash.Write([]byte(name))
	colorIndex := hash.Sum32() % 256
	return lipgloss.Color(fmt.Sprintf("%d", colorIndex))
}
