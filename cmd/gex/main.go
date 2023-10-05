package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hizkifw/gex/doc"
	"github.com/hizkifw/gex/internal/display"
)

func main() {
	m := display.NewModel()

	if len(os.Args) < 2 {
		fmt.Println("Usage: gex <filename>")
		fmt.Println("")
		fmt.Println("For help, run:")
		fmt.Println("  gex -h")
		os.Exit(1)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fname := "help"
		if len(os.Args) > 2 {
			fname = os.Args[2]
		}
		fname += ".md"
		b, err := doc.Docs.ReadFile(fname)
		if err != nil {
			fmt.Printf("Error loading help file: %v", err)
			os.Exit(1)
		}
		os.Stdout.Write(b)
		os.Exit(0)
	}

	if os.Args[1] == "--list-help" {
		files, err := doc.Docs.ReadDir(".")
		if err != nil {
			fmt.Printf("Error loading help files: %v", err)
			os.Exit(1)
		}
		fmt.Println("Availble help files:")
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".md") {
				fmt.Printf("- %s\n", f.Name()[:len(f.Name())-3])
			}
		}
		os.Exit(0)
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
