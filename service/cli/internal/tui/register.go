package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/cli/internal/app"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/cli/internal/client/grpc/dto"
)

var (
	focusedButtonRegister = focusedStyle.Render("[ Register ]")
	blurredButtonRegister = fmt.Sprintf("[ %s ]", blurredStyle.Render("Register"))
)

type registerModel struct {
	ctx        context.Context
	sp         *app.ServiceProvider
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	err        error
}

func InitialRegisterModel(ctx context.Context, sp *app.ServiceProvider) registerModel {
	m := registerModel{
		ctx:    ctx,
		sp:     sp,
		inputs: make([]textinput.Model, 4),
		err:    nil,
	}

	m.inputs[0] = newTextInput("Username", true)
	m.inputs[1] = newTextInput("Email", false)
	m.inputs[1].CharLimit = 64
	m.inputs[2] = newPasswordInput("Password")
	m.inputs[3] = newPasswordInput("Confirm Password")

	return m
}

func (m registerModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m registerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlR:
			return m.handleCursorChange()

		case tea.KeyTab, tea.KeyShiftTab:
			loginModel := InitialAuthModel(m.ctx, m.sp)
			return loginModel, loginModel.Init()

		case tea.KeyEnter:
			if m.focusIndex == len(m.inputs) {
				return m.handleRegister()
			}
			return m.handleFocusChange(msg)

		case tea.KeyUp, tea.KeyDown:
			return m.handleFocusChange(msg)
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m registerModel) View() string {
	var b strings.Builder

	for _, input := range m.inputs {
		b.WriteString(input.View())
		b.WriteRune('\n')
	}

	button := &blurredButtonRegister
	if m.focusIndex == len(m.inputs) {
		button = &focusedButtonRegister
	}

	fmt.Fprintf(&b, "\n%s\n\n", *button)

	if m.err != nil {
		b.WriteString(fmt.Sprintf("Error: %v\n", m.err))
	}

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("press Tab to switch to login."))

	content := addMargin(b.String())

	return content

}

func (m registerModel) handleCursorChange() (tea.Model, tea.Cmd) {
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

func (m registerModel) handleRegister() (tea.Model, tea.Cmd) {
	name := m.inputs[0].Value()
	email := m.inputs[1].Value()
	password := m.inputs[2].Value()
	passwordConfirm := m.inputs[3].Value()

	if name == "" || email == "" || password == "" || passwordConfirm == "" {
		m.err = fmt.Errorf("empty fields")
		return m, nil
	}

	if password != passwordConfirm {
		m.err = fmt.Errorf("passwords do not match")
		return m, nil
	}

	_, err := m.sp.HandlerService(m.ctx).Register(m.ctx, dto.RegisterInputDto{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
		Role:            1, // user
	})

	if err != nil {
		m.err = err
		return m, nil
	}

	loginModel := InitialAuthModel(m.ctx, m.sp)
	return loginModel, loginModel.Init()
}

func (m registerModel) handleFocusChange(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

func (m *registerModel) updateFocus() tea.Cmd {
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

func (m *registerModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}
