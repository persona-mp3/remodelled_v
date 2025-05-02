package main

import (
	"fmt"
	"os"
	
	tea "github.com/charmbracelet/bubbletea"
)


type model struct {
	choices []string
	cursor int
	selected map[int]struct{}
}



func initModel() model{
	return model {
		choices : []string{"task1", "task2", "task3a"},
		selected: make(map[int]struct{}),
	}
}


func (m model) Init() tea.Cmd {
	return nil
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// is the message a press
	case tea.KeyMsg:
		// okay, lets find out what key was actually pressed
		switch msg.String() {
			
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k": 
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices) - 1 {
				m.cursor++
			}
		case "enter", " ": 
			_, valid := m.selected[m.cursor]
			if valid {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil			
}



func (m model)View() string {
	s := "What task are you currently on: \n\n"

	for i, choice := range m.choices {

		cursor := "  "

		if m.cursor == i {
			cursor = ">"
		}

		// checks if the choice has been selected
		checked := " " // this has not been selected
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)

	}

	s += "\nEnter q to quit\n"

	//send the ui for rendering
	return s
}


func main(){
	p := tea.NewProgram(initModel())
	
	if _, err := p.Run(); err != nil {
		fmt.Println("skill issues noob~")
		fmt.Println(err)
		os.Exit(1)
	}
}
