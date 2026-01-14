package main

import (
	"cooler/internal/actions"
	"cooler/internal/diagnostics"
	"cooler/internal/gemini"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle         = lipgloss.NewStyle().Margin(1, 2)
	titleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
	statusStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	menuStyle        = lipgloss.NewStyle().Margin(1, 0, 0, 2)
	errorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	geminiResultStyle = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))
)

const asciiArt = `
  ____                __
 / __ \ ____ _ _____ / /___ _____
/ /_/ // __ '// ___// // _ \/ ___/
/ _, _// /_/ // /__ / //  __/ /
/_/ |_|\__,_/ \___//_/ \___/_/
`

type model struct {
	maxTemp        float64
	topProcess     *diagnostics.ProcessInfo
	cursor         int
	choices        []string
	status         string
	geminiAnalysis string
	loading        bool
	width          int // New field for terminal width
	err            error
}

type dataMsg struct {
	maxTemp    float64
	topProcess *diagnostics.ProcessInfo
	err        error
}

type geminiAnalysisMsg string

func fetchData() tea.Msg {
	// Fetch temperature and process info in parallel for better performance
	var maxTemp float64
	var topProcess *diagnostics.ProcessInfo
	var err1, err2 error

	done := make(chan struct{}, 2)

	go func() {
		maxTemp, err1 = diagnostics.GetMaxCpuTemperature()
		done <- struct{}{}
	}()

	go func() {
		topProcess, err2 = diagnostics.GetTopProcessInfo()
		done <- struct{}{}
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	if err1 != nil || err2 != nil {
		return dataMsg{err: fmt.Errorf("failed to fetch system data")}
	}
	return dataMsg{maxTemp: maxTemp, topProcess: topProcess}
}

func (m model) askGeminiCmd() tea.Msg {
	analysis, err := gemini.AskGemini(m.maxTemp, m.topProcess)
	if err != nil {
		return dataMsg{err: err} // Re-use dataMsg for errors
	}
	return geminiAnalysisMsg(analysis)
}

func initialModel() model {
	return model{
		choices: []string{"Kill Process", "Powersave Mode", "Default Mode", "Ask Gemini", "Refresh Data", "Quit"},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(fetchData, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg: // Handle window resizing
		m.width = msg.Width
		return m, nil

	case dataMsg:
		m.loading = false
		m.maxTemp = msg.maxTemp
		m.topProcess = msg.topProcess
		m.err = msg.err
		return m, nil

	case geminiAnalysisMsg:
		m.loading = false
		m.geminiAnalysis = string(msg)
		m.status = "Gemini analysis complete."
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter":
			// Clear previous status/analysis on new action
			m.status = ""
			m.geminiAnalysis = ""

			choice := m.choices[m.cursor]
			switch choice {
			case "Kill Process":
				if m.topProcess != nil {
					if err := actions.KillProcess(m.topProcess.PID); err != nil {
						m.status = fmt.Sprintf("Failed to kill %s", m.topProcess.Name)
					} else {
						m.status = fmt.Sprintf("Killed process %s", m.topProcess.Name)
					}
					return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return fetchData() }) // Refresh after action
				}
			case "Powersave Mode":
				actions.SetCpuGovernor("powersave")
				m.status = "Switched to Powersave mode."
			case "Default Mode":
				actions.SetCpuGovernor("schedutil")
				m.status = "Switched to Default (schedutil) mode."
			case "Ask Gemini":
				m.loading = true
				m.status = "Asking Gemini..."
				return m, m.askGeminiCmd
			case "Refresh Data":
				m.loading = true
				m.status = "Refreshing data..."
				return m, fetchData
			case "Quit":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	// Use strings.Builder for efficient string concatenation
	var sb strings.Builder
	sb.Grow(512) // Pre-allocate reasonable capacity to reduce allocations

	sb.WriteString(titleStyle.Render(asciiArt))
	sb.WriteString("\n--- Cooler (Go/Bubble Tea Edition) ---\n")

	if m.loading {
		sb.WriteString(docStyle.Render("Loading..."))
	} else if m.topProcess != nil {
		cpuFloat, _ := strconv.ParseFloat(m.topProcess.CPU, 64)
		diag := fmt.Sprintf("🌡️ Max CPU Temp: %.1f°C\n🔥 Top Process: '%s' (PID: %s) @ %.1f%% CPU",
			m.maxTemp, m.topProcess.Name, m.topProcess.PID, cpuFloat)
		sb.WriteString(docStyle.Render(diag))
	} else {
		sb.WriteString(docStyle.Render("✅ No significant CPU usage detected."))
	}

	// Build menu with strings.Builder
	var menuBuilder strings.Builder
	menuBuilder.Grow(128)
	for i, choice := range m.choices {
		if m.cursor == i {
			menuBuilder.WriteString("> ")
		} else {
			menuBuilder.WriteString("  ")
		}
		menuBuilder.WriteString(choice)
		menuBuilder.WriteByte('\n')
	}
	sb.WriteString(menuStyle.Render(menuBuilder.String()))

	if m.geminiAnalysis != "" {
		// Calculate width for the box, leaving some margin
		boxWidth := m.width - docStyle.GetHorizontalFrameSize() - geminiResultStyle.GetHorizontalFrameSize()
		if boxWidth < 10 {
			boxWidth = 10
		}
		analysisStyle := geminiResultStyle.Width(boxWidth)
		sb.WriteByte('\n')
		sb.WriteString(analysisStyle.Render(m.geminiAnalysis))
	}

	if m.status != "" {
		sb.WriteString(statusStyle.Render("\nStatus: " + m.status))
	}

	sb.WriteString("\n(q to quit)")
	return sb.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
