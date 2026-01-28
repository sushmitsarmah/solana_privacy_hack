package cli

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type inputForm struct {
	title       string
	inputs      []textinput.Model
	focusIndex  int
	submitFunc  func([]string) tea.Cmd
	cancelFunc  func() tea.Cmd
}

func newInputForm(title string, fields []string, submitFunc func([]string) tea.Cmd) inputForm {
	inputs := make([]textinput.Model, len(fields))
	for i, field := range fields {
		ti := textinput.New()
		ti.Placeholder = field
		ti.CharLimit = 156

		if i == 0 {
			ti.Focus()
			ti.PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)
			ti.TextStyle = lipgloss.NewStyle().Foreground(textColor)
		}

		inputs[i] = ti
	}

	return inputForm{
		title:      title,
		inputs:     inputs,
		focusIndex: 0,
		submitFunc: submitFunc,
	}
}

func (f *inputForm) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Submit on enter from last field
			if s == "enter" && f.focusIndex == len(f.inputs)-1 {
				values := make([]string, len(f.inputs))
				for i, input := range f.inputs {
					values[i] = input.Value()
				}
				return f.submitFunc(values)
			}

			// Cycle through inputs
			if s == "up" || s == "shift+tab" {
				f.focusIndex--
			} else {
				f.focusIndex++
			}

			if f.focusIndex > len(f.inputs)-1 {
				f.focusIndex = 0
			} else if f.focusIndex < 0 {
				f.focusIndex = len(f.inputs) - 1
			}

			cmds := make([]tea.Cmd, len(f.inputs))
			for i := 0; i <= len(f.inputs)-1; i++ {
				if i == f.focusIndex {
					cmds[i] = f.inputs[i].Focus()
					f.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(primaryColor)
					f.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(textColor)
					continue
				}
				f.inputs[i].Blur()
				f.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(subtleColor)
				f.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(subtleColor)
			}

			return tea.Batch(cmds...)
		}
	}

	// Handle character input
	cmd := f.updateInputs(msg)
	return cmd
}

func (f *inputForm) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.inputs))
	for i := range f.inputs {
		f.inputs[i], cmds[i] = f.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (f inputForm) View(width, height int) string {
	title := titleStyle.Render(f.title)

	var inputsView string
	for i, input := range f.inputs {
		inputsView += input.View() + "\n"
		if i < len(f.inputs)-1 {
			inputsView += "\n"
		}
	}

	help := helpStyle.Render("tab/shift+tab: navigate • enter: submit • esc: cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		inputsView,
		"",
		help,
	)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
