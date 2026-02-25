package main

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
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

	prTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFE66D")). // Bright yellow
			Bold(true)

	prNumberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98D8C8")). // Mint green
			Bold(true)

	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA07A")). // Light salmon
			Bold(true)
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
		str = selectedStyle.Render(str)
		desc = "  "  +
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

type prDelegate struct{}

func (d prDelegate) Height() int                               { return 3 }
func (d prDelegate) Spacing() int                              { return 0 }
func (d prDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d prDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(prItem)
	if !ok {
		return
	}

	title := i.Title()
	desc := i.Description()

	if index == m.Index() {
		title = selectedStyle.Render(title)
		desc = "  " + bodyStyle.Render(desc)
	} else {
		titleColor := getColorForIndex(index)
		title = lipgloss.NewStyle().Foreground(lipgloss.Color(titleColor)).Render("  " + title)
		desc = bodyStyle.Render("  " + desc)
	}

	fmt.Fprint(w, title+"\n"+desc+"\n")
}

// Helper function to get different colors for different items
func getColorForIndex(index int) string {
	colors := []string{
		"#FF6B6B", 
		"#4ECDC4", 
		"#45B7D1", 
		"#96CEB4", 
		"#FFEAA7", 
		"#DDA0DD", 
		"#98D8C8", 
		"#F7DC6F", 
		"#AED6F1", 
		"#D7BDE2", 
		"#A9DFBF", 
		"#F8C471", 
	}
	return colors[index%len(colors)]
}

type repoItem struct {
	r RepoInfo
}

func (i repoItem) Title() string { return i.r.Name }
func (i repoItem) Description() string {
	lang := i.r.Language
	if lang == "" {
		lang = "N/A"
	}
	return fmt.Sprintf("%s", lang)
}
func (i repoItem) FilterValue() string { return i.r.Name }

// Wrap PullRequest as a list.Item
type prItem struct {
	pr PullRequest
}

func (i prItem) Title() string {
	// Truncate long titles for better display
	title := i.pr.Title
	if len(title) > 60 {
		title = title[:57] + "..."
	}
	return fmt.Sprintf("#%d: %s", i.pr.Number, title)
}

func (i prItem) Description() string {
	return fmt.Sprintf("🔗 %s", i.pr.HTMLURL)
}

func (i prItem) FilterValue() string { return i.pr.Title }


// Message types for async operations
type pullRequestsLoadedMsg struct {
	prs []PullRequest
	err error
}


// Command to fetch pull requests
func fetchPullRequestsCmd(pullsURL string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		prs, err := fetchRecentPulls(pullsURL, 10)
		return pullRequestsLoadedMsg{prs: prs, err: err}
	})
}


type viewState int

const (
	repoListView viewState = iota
	repoDetailView
	pullRequestView
)

type model struct {
	repos        []RepoInfo
	list         list.Model
	prList       list.Model
	currentView  viewState
	selected     RepoInfo
	width        int
	height       int
	loading      bool
	loadingText  string
	textInput    string
	cursor       int
}

func newModel(repos []RepoInfo) model {
	items := make([]list.Item, len(repos))
	for i, r := range repos {
		items[i] = repoItem{r}
	}

	delegate := customDelegate{}
	l := list.New(items, delegate, 0, 0)
	l.Title = "Top Active Repositories (Last 7 Days)"
	l.SetShowHelp(true)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	l.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("#FF6B6B")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)

	prDelegate := prDelegate{}
	prList := list.New([]list.Item{}, prDelegate, 0, 0) // Start with 0 dimensions
	prList.Title = "Recent Pull Requests"
	prList.SetShowHelp(true)
	prList.SetShowStatusBar(true)
	prList.SetFilteringEnabled(true)
	prList.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("#4ECDC4")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)

}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4) // Leave space for title and help
		m.prList.SetWidth(msg.Width)
		m.prList.SetHeight(msg.Height - 4)
		return m, nil

	case pullRequestsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			// Handle error - you might want to show an error message
			m.loadingText = fmt.Sprintf("Error loading pull requests: %v", msg.err)
			return m, nil
		}

		// Convert PRs to list items
		items := make([]list.Item, len(msg.prs))
		for i, pr := range msg.prs {
			items[i] = prItem{pr}
		}
		m.prList.SetItems(items)
		m.currentView = pullRequestView
		return m, nil

	case tea.KeyMsg:

		// Handle regular navigation keys for other views
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			switch m.currentView {
			case repoListView:
				idx := m.list.Index()
				if idx >= 0 && idx < len(m.repos) {
					m.selected = m.repos[idx]
					m.currentView = repoDetailView
				}
			}
			return m, nil

		case "b":
			switch m.currentView {
			case repoDetailView, pullRequestView:
				m.currentView = repoListView
			}
			return m, nil

		case "a":
			return m, nil

		case "p":
			if m.currentView == repoDetailView && !m.loading {
				m.loading = true
				m.loadingText = "Loading pull requests..."
				return m, fetchPullRequestsCmd(m.selected.PullsURL)
			}
			return m, nil

		case "esc":
			switch m.currentView {
			case repoDetailView:
				m.currentView = repoListView
			case pullRequestView:
				m.currentView = repoDetailView
			}
			return m, nil
		}
	}

	// Update the appropriate list based on current view
	switch m.currentView {
	case repoListView:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	case pullRequestView:
		var cmd tea.Cmd
		m.prList, cmd = m.prList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	switch m.currentView {
	case repoListView:
		help := helpStyle.Render("Press 'enter' to select • Press 'q' to quit")
		return fmt.Sprintf("\n%s\n\n%s", m.list.View(), help)
	case repoDetailView:
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

		var help string
		if m.loading {
			loadingBar := m.renderLoadingBar()
			help = fmt.Sprintf("%s\n\n%s",
				bodyStyle.Render(m.loadingText),
				loadingBar)
		} else {
			help = helpStyle.Render("Press 'p' for pull requests • Press 'b' or 'esc' to go back • Press 'q' to quit")
		}

		return fmt.Sprintf("\n%s\n\n%s\n\n%s", header, details, help)
	case pullRequestView:
		header := headerStyle.Render(fmt.Sprintf("🔥 Pull Requests for %s", m.selected.Name))
		help := helpStyle.Render("Press 'b' or 'esc' to go back • Press 'q' to quit")

		if len(m.prList.Items()) == 0 {
			noItems := bodyStyle.Render("No recent pull requests found.")
			return fmt.Sprintf("\n%s\n\n%s\n\n%s", header, noItems, help)
		}

		return fmt.Sprintf("\n%s\n%s\n\n%s", header, m.prList.View(), help)
	}

	return ""
}

func (m model) renderLoadingBar() string {
	dots := strings.Repeat(".", (int(time.Now().UnixNano()/1e8)%4)+1)
	return loadingStyle.Render("Loading" + dots)
}

func launchUI(repos []RepoInfo) {
	p := tea.NewProgram(newModel(repos), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error launching UI: %v\n", err)
		os.Exit(1)
	}
}
