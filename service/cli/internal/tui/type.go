package tui

import tea "github.com/charmbracelet/bubbletea"

type errMsg error

type teaModelMsg struct {
	model tea.Model
	cmd   tea.Cmd
}
