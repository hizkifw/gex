package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hizkifw/gex/internal/display"
)

func main() {
	m := display.NewModel()

	if len(os.Args) < 2 {
		fmt.Println("Usage: gex <filename>")
		os.Exit(1)
	}

	if err := m.LoadFile(os.Args[1]); err != nil {
		fmt.Printf("Error loading file: %v", err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
