package main

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
)

func updateLogs(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	header := header(m, "Logs")
	footer := footer(m, "Press q to quit")
	content := content(m, header, footer)
	headerHeight := lipgloss.Height(header)
	footerHeight := lipgloss.Height(footer)
	contentHeight := lipgloss.Height(content)

	if _, ok := msg.(tea.WindowSizeMsg); ok {
		m.viewport = viewport.New(m.width, contentHeight-headerHeight-footerHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.SetContent(content)
	} else {
		m.viewport.Height = m.height - headerHeight - footerHeight
		m.viewport.Width = m.width
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func logsScreen(m model) string {

	header := header(m, "Logs")
	footer := footer(m, "Press q to quit")
	content := content(m, header, footer)
	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}

func content(m model, header, footer string) string {
	logs := wrap.String(m.content, lipgloss.Width(m.content))
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - lipgloss.Height(header) - lipgloss.Height(footer)).
		Padding(10).
		AlignVertical(lipgloss.Center).
		Render(logs)

}

func header(m model, title string) string {
	return lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Render(title)
}

func footer(m model, footer string) string {
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render(footer)
}
