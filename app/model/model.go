package model

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maitaken/monitor/app/executor"
)

var (
	successStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#008000"))
	failedStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#AA383E"))
	executingStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F6C18B"))
)

type Model struct {
	c       <-chan executor.CommandState
	spinner spinner.Model
	state   executor.CommandState
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.checkUpdate,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case executor.CommandState:
		m.state = msg
		return m, m.checkUpdate
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	switch m.state.State {
	case executor.Done:
		return fmt.Sprintf("%s:\n%s", successStyle.Render("Done"), string(m.state.Output))
	case executor.Error:
		return fmt.Sprintf("%s:\n%s", failedStyle.Render("Command Failed"), string(m.state.Output))
	case executor.Executing:
		return fmt.Sprintf("%s: %s", executingStyle.Render("Executing"), m.spinner.View())
	default:
		return "invalid state"
	}
}

func (m Model) Run() error {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func NewModel(c <-chan executor.CommandState) Model {
	return Model{
		c:       c,
		spinner: spinner.New(),
	}
}

func (m Model) checkUpdate() tea.Msg {
	r := <-m.c
	return r
}
