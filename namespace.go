package main

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updateNamespaces(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.list)-1 {
				m.cursor++
			}

		case "enter":
			m.state["namespace"] = m.list[m.cursor]
			pods, err := ListPodsByNamespaces(context.Background(), m.clientset, m.state["namespace"])
			if err != nil {
				panic(err)
			}

			m.list = pods
			m.cursor = 0
			m.currentPage = "pods"
		}
	}

	return m, nil
}

func namespaceScreen(m model) string {
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Render("Namespaces")

	footer := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Render("Press q to quit.")

	s := ""
	for i, ns := range m.list {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, ns)
	}

	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-lipgloss.Height(header)-lipgloss.Height(footer)).
		Align(lipgloss.Center, lipgloss.Center).
		Render(s)

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}
