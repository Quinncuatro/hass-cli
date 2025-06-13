package tui

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quinncuatro/hass-cli/internal/client"
	"github.com/quinncuatro/hass-cli/internal/config"
)

type App struct {
	config *config.Config
	client *client.HomeAssistantClient
}

type model struct {
	config     *config.Config
	client     *client.HomeAssistantClient
	entities   []client.EntityState
	cursor     int
	selected   map[int]struct{}
	width      int
	height     int
	filter     string
	loading    bool
	err        error
	statusMsg  string
}

type entitiesLoadedMsg []client.EntityState
type errorMsg error
type statusMsg string

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
			Padding(0, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#EE6FF8")).
				Bold(true).
				Padding(0, 2)

	paginationStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)
)

func NewApp(cfg *config.Config, client *client.HomeAssistantClient) *App {
	return &App{
		config: cfg,
		client: client,
	}
}

func (a *App) Run(ctx context.Context) error {
	m := model{
		config:   a.config,
		client:   a.client,
		selected: make(map[int]struct{}),
		loading:  true,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	
	// Load entities in background
	go func() {
		entities, err := a.client.GetStates(ctx)
		if err != nil {
			p.Send(errorMsg(err))
			return
		}
		p.Send(entitiesLoadedMsg(entities))
	}()

	_, err := p.Run()
	if err != nil && strings.Contains(err.Error(), "TTY") {
		return fmt.Errorf("TUI requires a proper terminal environment. Try running directly from your terminal")
	}
	return err
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case entitiesLoadedMsg:
		m.entities = []client.EntityState(msg)
		m.loading = false
		m.filterAndSortEntities()

	case errorMsg:
		m.err = error(msg)
		m.loading = false

	case statusMsg:
		m.statusMsg = string(msg)
		return m, tea.Tick(time.Second*3, func(t time.Time) tea.Msg {
			return statusMsg("")
		})

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.entities)-1 {
				m.cursor++
			}

		case "enter", "space":
			return m, m.toggleEntity()

		case "r":
			m.loading = true
			return m, m.refreshEntities()

		case "/":
			// TODO: Implement search/filter
			return m, nil

		case "?":
			// TODO: Show help
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.loading {
		return "\n  Loading entities...\n"
	}

	if m.err != nil {
		return fmt.Sprintf("\n  %s\n", errorStyle.Render("Error: "+m.err.Error()))
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("Home Assistant - Entity Control"))
	b.WriteString("\n\n")

	// Status message
	if m.statusMsg != "" {
		b.WriteString(statusStyle.Render(m.statusMsg))
		b.WriteString("\n\n")
	}

	// Entity list
	start, end := m.paginate()
	for i := start; i < end; i++ {
		entity := m.entities[i]
		
		// Format entity display
		friendlyName := entity.EntityID
		if name, ok := entity.Attributes["friendly_name"].(string); ok {
			friendlyName = name
		}

		// State indicator
		stateIcon := "○"
		stateColor := lipgloss.Color("#626262")
		switch entity.State {
		case "on":
			stateIcon = "●"
			stateColor = lipgloss.Color("#04B575")
		case "off":
			stateIcon = "○"
			stateColor = lipgloss.Color("#626262")
		case "unavailable":
			stateIcon = "✗"
			stateColor = lipgloss.Color("#FF0000")
		default:
			stateIcon = "◐"
			stateColor = lipgloss.Color("#FFB86C")
		}

		stateStyle := lipgloss.NewStyle().Foreground(stateColor)
		prefix := fmt.Sprintf("%s %s", stateStyle.Render(stateIcon), friendlyName)
		
		// Add domain info
		domain := strings.Split(entity.EntityID, ".")[0]
		suffix := fmt.Sprintf(" [%s]", domain)

		line := prefix + suffix

		if i == m.cursor {
			b.WriteString(selectedItemStyle.Render("> " + line))
		} else {
			b.WriteString(itemStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}

	// Pagination info
	if len(m.entities) > 0 {
		b.WriteString("\n")
		b.WriteString(paginationStyle.Render(fmt.Sprintf("Showing %d-%d of %d entities", start+1, end, len(m.entities))))
	}

	// Help
	b.WriteString("\n\n")
	help := "↑/k: up • ↓/j: down • enter/space: toggle • r: refresh • q: quit"
	b.WriteString(helpStyle.Render(help))

	return b.String()
}

func (m *model) filterAndSortEntities() {
	// Filter to controllable entities (lights, switches, fans, etc.)
	controllable := []string{"light", "switch", "fan", "climate", "cover"}
	var filtered []client.EntityState
	
	for _, entity := range m.entities {
		domain := strings.Split(entity.EntityID, ".")[0]
		for _, c := range controllable {
			if domain == c {
				filtered = append(filtered, entity)
				break
			}
		}
	}

	// Sort by domain, then by friendly name
	sort.Slice(filtered, func(i, j int) bool {
		domainI := strings.Split(filtered[i].EntityID, ".")[0]
		domainJ := strings.Split(filtered[j].EntityID, ".")[0]
		
		if domainI != domainJ {
			return domainI < domainJ
		}
		
		nameI := filtered[i].EntityID
		if fn, ok := filtered[i].Attributes["friendly_name"].(string); ok {
			nameI = fn
		}
		nameJ := filtered[j].EntityID
		if fn, ok := filtered[j].Attributes["friendly_name"].(string); ok {
			nameJ = fn
		}
		
		return nameI < nameJ
	})

	m.entities = filtered
}

func (m model) paginate() (int, int) {
	maxItems := m.height - 8 // Leave room for title, help, etc.
	if maxItems <= 0 {
		maxItems = 10
	}

	start := 0
	end := len(m.entities)

	// Simple pagination around cursor
	if len(m.entities) > maxItems {
		half := maxItems / 2
		start = m.cursor - half
		if start < 0 {
			start = 0
		}
		end = start + maxItems
		if end > len(m.entities) {
			end = len(m.entities)
			start = end - maxItems
			if start < 0 {
				start = 0
			}
		}
	}

	return start, end
}

func (m model) toggleEntity() tea.Cmd {
	if m.cursor >= len(m.entities) {
		return nil
	}

	entity := m.entities[m.cursor]
	return tea.Cmd(func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		var action string

		switch entity.State {
		case "on":
			err = m.client.TurnOffEntity(ctx, entity.EntityID)
			action = "turned off"
		case "off":
			err = m.client.TurnOnEntity(ctx, entity.EntityID)
			action = "turned on"
		default:
			err = m.client.ToggleEntity(ctx, entity.EntityID)
			action = "toggled"
		}

		if err != nil {
			return errorMsg(fmt.Errorf("failed to control %s: %w", entity.EntityID, err))
		}

		friendlyName := entity.EntityID
		if name, ok := entity.Attributes["friendly_name"].(string); ok {
			friendlyName = name
		}

		return statusMsg(fmt.Sprintf("✓ %s %s", friendlyName, action))
	})
}

func (m model) refreshEntities() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		entities, err := m.client.GetStates(ctx)
		if err != nil {
			return errorMsg(err)
		}
		return entitiesLoadedMsg(entities)
	})
}