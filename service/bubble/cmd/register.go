package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/bubble/internal/app"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/bubble/internal/client/grpc/dto"
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

func initialRegisterModel(ctx context.Context, sp *app.ServiceProvider) registerModel {
	m := registerModel{
		ctx:    ctx,
		sp:     sp,
		inputs: make([]textinput.Model, 4),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		t.Width = 20

		switch i {
		case 0:
			t.Placeholder = "Username"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Email"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'

		case 3:
			t.Placeholder = "Confirm Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'

		}

		m.inputs[i] = t
	}

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
			m.cursorMode++

			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}

			cmds := make([]tea.Cmd, len(m.inputs))

			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}

			return m, tea.Batch(cmds...)

		case tea.KeyTab, tea.KeyShiftTab:
			return initialAuthModel(m.ctx, m.sp), nil

		case tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			s := msg.Type

			if s == tea.KeyEnter {
				if m.focusIndex == len(m.inputs) {
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
					})

					if err != nil {
						m.err = err
						return m, nil
					}

					loginModel := initialAuthModel(m.ctx, m.sp)
					return loginModel, loginModel.Init()
				}
			}

			if s == tea.KeyUp {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))

			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
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

	return b.String()
}

func (m *registerModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
