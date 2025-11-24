package main

import (
	"context"
	"flag"
	"path/filepath"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type model struct {
	list          []string
	cursor        int
	clientset     *kubernetes.Clientset
	width, height int
	currentPage   string
	state         map[string]string
	content       string
	viewport      viewport.Model
}

func (m model) Init() tea.Cmd {

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	}

	switch m.currentPage {
	case "pods":
		return updatePods(msg, m)
	case "logs":
		return updateLogs(msg, m)
	}

	return updateNamespaces(msg, m)
}

func (m model) View() string {
	s := "Initializing..."

	switch m.currentPage {
	case "pods":
		s = podsScreen(m)
	case "logs":
		s = logsScreen(m)
	default:
		s = namespaceScreen(m)
	}

	return s
}

func main() {
	ctx := context.Background()
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the config")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the config")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	namespaces, err := ListNamespaces(ctx, clientset)
	if err != nil {
		panic(err)
	}

	m := model{
		list:        namespaces,
		clientset:   clientset,
		currentPage: "namespaces",
		state:       make(map[string]string),
		viewport:    viewport.Model{},
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}

}
