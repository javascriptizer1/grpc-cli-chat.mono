package tui

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/cli/internal/app"
)

var (
	focusedLoginButton = focusedStyle.Render("[ Login ]")
	blurredLoginButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Login"))
)

type authModel struct {
	ctx        context.Context
	sp         *app.ServiceProvider
	spinner    spinner.Model
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	width      int
	height     int
	loading    bool
	err        error
}

func InitialAuthModel(ctx context.Context, sp *app.ServiceProvider) authModel {
	m := authModel{
		ctx:     ctx,
		sp:      sp,
		spinner: initSpinner(),
		inputs:  make([]textinput.Model, 2),
		loading: false,
		err:     nil,
	}

	m.inputs[0] = newTextInput("Email", true)
	m.inputs[1] = newPasswordInput("Password")

	return m
}

func (m authModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func (m authModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlR:
			return m.handleCursorChange()

		case tea.KeyTab, tea.KeyShiftTab:
			return m.switchToRegister()

		case tea.KeyEnter:
			if m.focusIndex == len(m.inputs) {
				return m.handleLogin()
			}
			return m.handleFocusChange(msg)

		case tea.KeyUp, tea.KeyDown:
			return m.handleFocusChange(msg)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case errMsg:
		m.err = msg
		m.loading = false
		return m, nil

	case teaModelMsg:
		m.loading = false
		return msg.model, msg.cmd
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m authModel) View() string {
	var b strings.Builder

	for _, input := range m.inputs {
		b.WriteString(input.View())
		b.WriteRune('\n')
	}

	button := &blurredLoginButton

	if m.focusIndex == len(m.inputs) {
		button = &focusedLoginButton
	}

	fmt.Fprintf(&b, "\n%s\n\n", *button)

	if m.err != nil {
		b.WriteString(fmt.Sprintf("Error: %v\n", m.err))
	} else if m.loading {
		b.WriteString(fmt.Sprintf("%s Loading...\n\n", m.spinner.View()))
	} else {
		b.WriteRune('\n')
	}

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("press Tab to switch to register."))

	content := addMargin(b.String())

	return content
}

func (m authModel) handleCursorChange() (tea.Model, tea.Cmd) {
	m.cursorMode++

	if m.cursorMode > cursor.CursorHide {
		m.cursorMode = cursor.CursorBlink
	}

	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
	}

	return m, tea.Batch(cmds...)
}

func (m authModel) switchToRegister() (tea.Model, tea.Cmd) {
	regModel := InitialRegisterModel(m.ctx, m.sp)
	return regModel, regModel.Init()
}

func (m authModel) handleLogin() (tea.Model, tea.Cmd) {
	login := m.inputs[0].Value()
	password := m.inputs[1].Value()

	if login == "" || password == "" {
		m.err = errors.New("empty fields")
		return m, nil
	}

	m.err = nil
	m.loading = true

	return m, tea.Batch(m.spinner.Tick, m.doLogin(login, password))
}

func (m authModel) doLogin(login, password string) tea.Cmd {
	return func() tea.Msg {
		_, err := m.sp.HandlerService(m.ctx).Login(m.ctx, login, password)

		if err != nil {
			return errMsg(err)
		}

		chatListModel := InitialChatListModel(m.ctx, m.sp, m.width, m.height)

		return teaModelMsg{model: chatListModel, cmd: chatListModel.Init()}
	}
}

func (m authModel) handleFocusChange(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyUp {
		m.focusIndex--
	} else {
		m.focusIndex++
	}

	if m.focusIndex > len(m.inputs) {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs)
	}

	return m, m.updateFocus()
}

func (m *authModel) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := 0; i < len(m.inputs); i++ {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].TextStyle = noStyle
		}
	}

	return tea.Batch(cmds...)
}

func (m *authModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func newTextInput(placeholder string, focused bool) textinput.Model {
	t := textinput.New()
	t.Cursor.Style = cursorStyle
	t.CharLimit = 32
	t.Width = 20
	t.Placeholder = placeholder

	if focused {
		t.Focus()
		t.PromptStyle = focusedStyle
		t.TextStyle = focusedStyle
	}

	return t
}

func newPasswordInput(placeholder string) textinput.Model {
	t := newTextInput(placeholder, false)
	t.EchoMode = textinput.EchoPassword
	t.EchoCharacter = '•'
	return t
}

func initSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return s
}
