package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Options struct {
	Continue bool
	Install  bool
}

func Run(options *Options) error {
	bubble := newBubble()

	switch {
	case options.Install:
		bubble.newState(scrapersInstallState)
	case options.Continue:
		_, err := bubble.loadHistory()
		if err != nil {
			return err
		}

		bubble.newState(historyState)
	default:
		bubble.newState(sourcesState)
	}

	program := tea.NewProgram(bubble, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}
