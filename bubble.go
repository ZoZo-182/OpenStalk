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
		title = selectedStyle.Render("🔥 " + title)
		desc = "  " + bodyStyle.Render(desc)
	} else {
		titleColor := getColorForIndex(index)
		title = lipgloss.NewStyle().Foreground(lipgloss.Color(titleColor)).Render("  " + title)
		desc = bodyStyle.Render("  " + desc)
	}

	fmt.Fprint(w, title+"\n"+desc+"\n")
}

type commitDelegate struct{}

func (d commitDelegate) Height() int                               { return 3 }
func (d commitDelegate) Spacing() int                              { return 0 }
func (d commitDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d commitDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(commitItem)
	if !ok {
		return
	}

	title := i.Title()
	desc := i.Description()

	if index == m.Index() {
		title = selectedStyle.Render("💾 " + title)
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

type commitItem struct {
	commit Commit
}

func (i commitItem) Title() string {
	// Get first line of commit message (title)
	lines := strings.Split(i.commit.Message, "\n")
	title := lines[0]
	if len(title) > 60 {
		title = title[:57] + "..."
	}
	return fmt.Sprintf("%.7s: %s", i.commit.SHA, title)
}

func (i commitItem) Description() string {
	// Show author and full message if it's multiline
	lines := strings.Split(i.commit.Message, "\n")
	if len(lines) > 1 && strings.TrimSpace(lines[1]) != "" {
		desc := strings.TrimSpace(lines[1])
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		return fmt.Sprintf("👤 %s • %s", i.commit.Author, desc)
	}
	return fmt.Sprintf("👤 %s", i.commit.Author)
}

func (i commitItem) FilterValue() string { return i.commit.Message }

// Message types for async operations
type pullRequestsLoadedMsg struct {
	prs []PullRequest
	err error
}

type commitsLoadedMsg struct {
	commits []Commit
	err     error
}

type bookmarkRepoLoadedMsg struct {
	repo RepoInfo
	err  error
}

// Command to fetch pull requests
func fetchPullRequestsCmd(pullsURL string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		prs, err := fetchRecentPulls(pullsURL, 10)
		return pullRequestsLoadedMsg{prs: prs, err: err}
	})
}

// Command to fetch commits
func fetchCommitsCmd(repoAPIURL string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		commits, err := fetchRecentCommits(repoAPIURL, 10)
		return commitsLoadedMsg{commits: commits, err: err}
	})
}

// Command to fetch bookmark repo details
func fetchBookmarkRepoCmd(repoURL string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		repo, err := fetchBookmarkRepoDetails(repoURL)
		return bookmarkRepoLoadedMsg{repo: repo, err: err}
	})
}

type healthCheckLoadedMsg struct {
	output string
	err    error
}

func runHealthCheckCmd(repoURL string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		exec.Command("chmod", "+x", "./health_check.sh").Run()

		cmd := exec.Command("./health_check.sh", repoURL)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			return healthCheckLoadedMsg{
				output: fmt.Sprintf("Error running health check: %v\nStderr: %s", err, stderr.String()),
				err:    err,
			}
		}

		return healthCheckLoadedMsg{
			output: stdout.String(),
			err:    nil,
		}
	})
}

type viewState int

const (
	repoListView viewState = iota
	repoDetailView
	pullRequestView
	commitView
	bookmarksView
	addBookmarkView
	healthCheckView
)

type model struct {
	repos        []RepoInfo
	bookmarks    []RepoInfo
	list         list.Model
	bookmarkList list.Model
	prList       list.Model
	commitList   list.Model
	currentView  viewState
	selected     RepoInfo
	width        int
	height       int
	loading      bool
	loadingText  string
	textInput    string
	cursor       int
	healthOutput string
}

func newModel(repos []RepoInfo) model {
	items := make([]list.Item, len(repos))
	for i, r := range repos {
		items[i] = repoItem{r}
	}

	delegate := customDelegate{}
	l := list.New(items, delegate, 0, 0)
	l.Title = "🌈 Top Active Repositories (Last 7 Days)"
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
	prList.Title = "🔥 Recent Pull Requests"
	prList.SetShowHelp(true)
	prList.SetShowStatusBar(true)
	prList.SetFilteringEnabled(true)
	prList.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("#4ECDC4")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)

	commitDelegate := commitDelegate{}
	commitList := list.New([]list.Item{}, commitDelegate, 0, 0) // Start with 0 dimensions
	commitList.Title = "💾 Recent Commits"
	commitList.SetShowHelp(true)
	commitList.SetShowStatusBar(true)
	commitList.SetFilteringEnabled(true)
	commitList.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("#9B59B6")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)

	bookmarkDelegate := customDelegate{} // Reuse the same delegate as repo list
	bookmarkList := list.New([]list.Item{}, bookmarkDelegate, 0, 0)
	bookmarkList.Title = "🔖 Bookmarked Repositories"
	bookmarkList.SetShowHelp(true)
	bookmarkList.SetShowStatusBar(true)
	bookmarkList.SetFilteringEnabled(true)
	bookmarkList.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("#E67E22")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)

	var bookmarks []RepoInfo
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Bookmark loading failed, continuing with empty bookmarks: %v\n", r)
				bookmarks = []RepoInfo{}
			}
		}()
		bookmarks = loadBookmarks()
	}()

	bookmarkItems := make([]list.Item, len(bookmarks))
	for i, bookmark := range bookmarks {
		bookmarkItems[i] = repoItem{bookmark}
	}
	bookmarkList.SetItems(bookmarkItems)

	return model{
		repos:        repos,
		bookmarks:    bookmarks,
		list:         l,
		bookmarkList: bookmarkList,
		prList:       prList,
		commitList:   commitList,
		currentView:  repoListView,
		width:        0,
		height:       0,
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
		m.prList.SetWidth(msg.Width)
		m.prList.SetHeight(msg.Height - 4)
		m.commitList.SetWidth(msg.Width)
		m.commitList.SetHeight(msg.Height - 4)
		m.bookmarkList.SetWidth(msg.Width)
		m.bookmarkList.SetHeight(msg.Height - 4)
		return m, nil

	case commitsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.loadingText = fmt.Sprintf("Error loading commits: %v", msg.err)
			return m, nil
		}

		// Convert commits to list items
		items := make([]list.Item, len(msg.commits))
		for i, commit := range msg.commits {
			items[i] = commitItem{commit}
		}
		m.commitList.SetItems(items)
		m.currentView = commitView
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

	case bookmarkRepoLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.loadingText = fmt.Sprintf("Error loading repository: %v", msg.err)
			return m, nil
		}

		// Save the bookmark
		saveBookmark(msg.repo)

		// Update the bookmarks list
		m.bookmarks = loadBookmarks()
		bookmarkItems := make([]list.Item, len(m.bookmarks))
		for i, bookmark := range m.bookmarks {
			bookmarkItems[i] = repoItem{bookmark}
		}
		m.bookmarkList.SetItems(bookmarkItems)

		// Clear text input and go back to bookmarks view
		m.textInput = ""
		m.cursor = 0
		m.currentView = bookmarksView
		m.loadingText = fmt.Sprintf("✅ Successfully bookmarked %s!", msg.repo.Name)
		return m, nil

	case healthCheckLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.healthOutput = fmt.Sprintf("❌ Health check failed: %v", msg.err)
		} else {
			m.healthOutput = msg.output
		}
		m.currentView = healthCheckView
		return m, nil

	case tea.KeyMsg:
		// Handle text input first when in addBookmarkView
		if m.currentView == addBookmarkView && !m.loading {
			switch msg.String() {
			case "enter":
				if strings.TrimSpace(m.textInput) != "" {
					m.loading = true
					m.loadingText = "Loading repository details..."
					return m, fetchBookmarkRepoCmd(strings.TrimSpace(m.textInput))
				}
				return m, nil
			case "esc":
				m.currentView = bookmarksView
				m.textInput = ""
				m.cursor = 0
				return m, nil
			case "backspace":
				if m.cursor > 0 {
					m.textInput = m.textInput[:m.cursor-1] + m.textInput[m.cursor:]
					m.cursor--
				}
				return m, nil
			case "left":
				if m.cursor > 0 {
					m.cursor--
				}
				return m, nil
			case "right":
				if m.cursor < len(m.textInput) {
					m.cursor++
				}
				return m, nil
			case "home":
				m.cursor = 0
				return m, nil
			case "end":
				m.cursor = len(m.textInput)
				return m, nil
			case "ctrl+c", "q":
				return m, tea.Quit
			default:
				// Handle regular character input (including a, d, p, etc.)
				if len(msg.String()) == 1 {
					char := msg.String()
					m.textInput = m.textInput[:m.cursor] + char + m.textInput[m.cursor:]
					m.cursor++
				}
				return m, nil
			}
		}

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
			case bookmarksView:
				idx := m.bookmarkList.Index()
				if idx >= 0 && idx < len(m.bookmarks) {
					m.selected = m.bookmarks[idx]
					m.currentView = repoDetailView
				}
			}
			return m, nil

		case "b":
			switch m.currentView {
			case repoListView:
				m.currentView = bookmarksView
			case repoDetailView, pullRequestView, commitView:
				m.currentView = repoListView
			}
			return m, nil

		case "a":
			if m.currentView == bookmarksView {
				m.currentView = addBookmarkView
				m.textInput = ""
				m.cursor = 0
			}
			return m, nil

		case "d", "delete":
			if m.currentView == bookmarksView && len(m.bookmarkList.Items()) > 0 {
				idx := m.bookmarkList.Index()
				if idx >= 0 && idx < len(m.bookmarks) {
					// Remove bookmark
					removeBookmark(m.bookmarks[idx].HTMLURL)

					// Update the bookmarks list
					m.bookmarks = loadBookmarks()
					bookmarkItems := make([]list.Item, len(m.bookmarks))
					for i, bookmark := range m.bookmarks {
						bookmarkItems[i] = repoItem{bookmark}
					}
					m.bookmarkList.SetItems(bookmarkItems)
				}
			}
			return m, nil

		case "p":
			if m.currentView == repoDetailView && !m.loading {
				m.loading = true
				m.loadingText = "Loading pull requests..."
				return m, fetchPullRequestsCmd(m.selected.PullsURL)
			}
			return m, nil

		case "h":
			if m.currentView == repoDetailView && !m.loading {
				m.loading = true
				m.loadingText = "Running repository health check..."
				return m, runHealthCheckCmd(m.selected.HTMLURL)
			}
			return m, nil

		case "c":
			if m.currentView == repoDetailView && !m.loading {
				m.loading = true
				m.loadingText = "Loading commits..."
				return m, fetchCommitsCmd(m.selected.APIURL)
			}
			return m, nil

		case "esc":
			switch m.currentView {
			case repoDetailView:
				m.currentView = repoListView
			case pullRequestView, commitView, healthCheckView:
				m.currentView = repoDetailView
			case bookmarksView:
				m.currentView = repoListView
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
	case bookmarksView:
		var cmd tea.Cmd
		m.bookmarkList, cmd = m.bookmarkList.Update(msg)
		return m, cmd
	case pullRequestView:
		var cmd tea.Cmd
		m.prList, cmd = m.prList.Update(msg)
		return m, cmd
	case commitView:
		var cmd tea.Cmd
		m.commitList, cmd = m.commitList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	switch m.currentView {
	case repoListView:
		help := helpStyle.Render("Press 'b' for bookmarks • Press 'enter' to select • Press 'q' to quit")
		return fmt.Sprintf("\n%s\n\n%s", m.list.View(), help)

	case bookmarksView:
		help := helpStyle.Render("Press 'a' to add bookmark • Press 'd' to delete • Press 'enter' to select • Press 'esc' to go back • Press 'q' to quit")
		if len(m.bookmarkList.Items()) == 0 {
			header := headerStyle.Render("🔖 No bookmarks yet")
			noItems := bodyStyle.Render("Press 'a' to add your first bookmark!")
			return fmt.Sprintf("\n%s\n\n%s\n\n%s", header, noItems, help)
		}
		return fmt.Sprintf("\n%s\n\n%s", m.bookmarkList.View(), help)

	case addBookmarkView:
		header := headerStyle.Render("🔖 Add Repository Bookmark")

		var content string
		if m.loading {
			loadingBar := m.renderLoadingBar()
			content = fmt.Sprintf("%s\n\n%s",
				bodyStyle.Render(m.loadingText),
				loadingBar)
		} else {
			prompt := bodyStyle.Render("Enter GitHub repository URL (e.g., https://github.com/owner/repo):")

			// Create text input display with cursor
			input := m.textInput
			if m.cursor < len(input) {
				input = input[:m.cursor] + "│" + input[m.cursor:]
			} else {
				input = input + "│"
			}

			inputDisplay := detailStyle.Render(input)
			content = fmt.Sprintf("%s\n\n%s", prompt, inputDisplay)
		}

		help := helpStyle.Render("Press 'enter' to add • Press 'esc' to cancel • Press 'q' to quit")
		return fmt.Sprintf("\n%s\n\n%s\n\n%s", header, content, help)

	case healthCheckView:
		header := headerStyle.Render(fmt.Sprintf("🏥 Health Check Results for %s", m.selected.Name))

		// Style the health output
		styledOutput := detailStyle.Render(m.healthOutput)

		help := helpStyle.Render("Press 'esc' to go back • Press 'q' to quit")
		return fmt.Sprintf("\n%s\n\n%s\n\n%s", header, styledOutput, help)

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
			help = helpStyle.Render("Press 'p' for pull requests • Press 'c' for commits • Press 'h' for health check • Press 'b' or 'esc' to go back • Press 'q' to quit")
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

	case commitView:
		header := headerStyle.Render(fmt.Sprintf("💾 Commits for %s", m.selected.Name))
		help := helpStyle.Render("Press 'b' or 'esc' to go back • Press 'q' to quit")

		if len(m.commitList.Items()) == 0 {
			noItems := bodyStyle.Render("No recent commits found.")
			return fmt.Sprintf("\n%s\n\n%s\n\n%s", header, noItems, help)
		}

		return fmt.Sprintf("\n%s\n%s\n\n%s", header, m.commitList.View(), help)
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
