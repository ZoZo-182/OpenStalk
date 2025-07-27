package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"os"
)

// Lipgloss styles with rainbow colors
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.NewStyle().Background(lipgloss.Color("#FF6B6B")). // Coral red
												PaddingLeft(1).PaddingRight(1).GetBackground()).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#4ECDC4")). // Turquoise
			Padding(0, 1).
			MarginBottom(1)

	bodyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#95E1D3")). // Light mint
			MarginLeft(2)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFE66D")). // Bright yellow
			Background(lipgloss.Color("#FF8B94")). // Pink
			Padding(0, 1)

	detailStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#A8E6CF")). // Light green
			Padding(1, 2).
			MarginTop(1).
			Background(lipgloss.Color("#1A1A2E")) // Dark blue background

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA07A")). // Light salmon
			MarginTop(1)

	// Additional colorful styles for variety
	languageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98D8C8")). // Mint green
			Bold(true)

	countStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F7DC6F")). // Light yellow
			Bold(true)

	urlStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AED6F1")). // Light blue
			Underline(true)
)

// Custom delegate for better styling
type customDelegate struct{}

func (d customDelegate) Height() int                               { return 2 }
func (d customDelegate) Spacing() int                              { return 0 }
func (d customDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d customDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(repoItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Title())
	desc := i.Description()

	if index == m.Index() {
		str = selectedStyle.Render("🌈 " + str)
		desc = "  " + countStyle.Render(fmt.Sprintf("%d issues", i.r.Count)) +
			" • " + languageStyle.Render(func() string {
			if i.r.Language == "" {
				return "N/A"
			}
			return i.r.Language
		}())
	} else {
		// Use different colors for unselected items
		titleColor := getColorForIndex(index)
		str = lipgloss.NewStyle().Foreground(lipgloss.Color(titleColor)).Render("  " + str)
		desc = bodyStyle.Render("  " + desc)
	}

	fmt.Fprint(w, str+"\n"+desc)
}

// Helper function to get different colors for different items
func getColorForIndex(index int) string {
	colors := []string{
		"#FF6B6B", // Red
		"#4ECDC4", // Turquoise
		"#45B7D1", // Blue
		"#96CEB4", // Green
		"#FFEAA7", // Yellow
		"#DDA0DD", // Plum
		"#98D8C8", // Mint
		"#F7DC6F", // Light yellow
		"#AED6F1", // Light blue
		"#D7BDE2", // Light purple
		"#A9DFBF", // Light green
		"#F8C471", // Orange
	}
	return colors[index%len(colors)]
}

type viewState int

const (
	listView viewState = iota
	detailView
	prListView
)

// Wrap your RepoInfo as a list.Item
type repoItem struct {
	r RepoInfo
}

func (i repoItem) Title() string { return i.r.Name }
func (i repoItem) Description() string {
	lang := i.r.Language
	if lang == "" {
		lang = "N/A"
	}
	return fmt.Sprintf("%d issues • %s", i.r.Count, lang)
}
func (i repoItem) FilterValue() string { return i.r.Name }

type model struct {
	repos    []RepoInfo
	list     list.Model
	detail   bool
	selected RepoInfo
	state    viewState
	width    int
	height   int
	prs      []PullRequest // Store fetched PRs
	prIndex  int
}

func newModel(repos []RepoInfo) model {
	items := make([]list.Item, len(repos))
	for i, r := range repos {
		items[i] = repoItem{r}
	}

	// Initialize the list with the custom delegate
	delegate := customDelegate{}
	l := list.New(items, delegate, 80, 20)
	l.Title = "🌈 Top Active Repositories (Last 7 Days)"
	l.SetShowHelp(true)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	// Style the list itself
	l.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("#FF6B6B")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)

	return model{
		repos:  repos,
		list:   l,
		state:  listView,
		width:  80,
		height: 20,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4) // Leave space for title and help
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if !m.detail {
				idx := m.list.Index()
				if idx >= 0 && idx < len(m.repos) {
					m.selected = m.repos[idx]
					m.detail = true
				}
			}
			return m, nil
		case "esc", "b":
			if m.detail {
				m.detail = false
			}
			return m, nil
		}
	}

	if !m.detail {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.detail {
		r := m.selected

		header := headerStyle.Render(fmt.Sprintf("📊 %s", r.Name))

		details := detailStyle.Render(
			"🔥 Issues in last week: " + countStyle.Render(fmt.Sprintf("%d", r.Count)) + "\n" +
				"💻 Primary Language: " + languageStyle.Render(func() string {
				if r.Language == "" {
					return "Not specified"
				}
				return r.Language
			}()) + "\n" +
				"🔗 Repository URL: " + urlStyle.Render(r.HTMLURL),
		)

		help := helpStyle.Render("Press 'b' or 'esc' to go back • Press 'q' to quit")

		return fmt.Sprintf("\n%s\n\n%s\n\n%s", header, details, help)
	}

	return fmt.Sprintf("\n%s\n", m.list.View())
}

func launchUI(repos []RepoInfo) {
	p := tea.NewProgram(newModel(repos), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error launching UI: %v\n", err)
		os.Exit(1)
	}
}
