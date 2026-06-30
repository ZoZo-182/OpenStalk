package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"strings"
	"os"
)

var CLI struct {
	NoBanner bool       `help:"Hide the OpenStalk banner."`
	Search   SearchCmd  `cmd:"" help:"Search Github for recent open pull requests."`
	Version  VersionCmd `cmd:"" help:"Print version info"`
}

type SearchCmd struct {
	Language string `help:"Programming Language to filter PRs by." short:"l"`
	Days     int    `help:"Number of days to look back." short:"d" default:"1"`
	Limit    int    `help:"Maximum number of result to show." short:"n" default:"10"`
	MinStars int    `help:"Minimum repo stars." default:"100"`
	MaxStars int    `help:"Maximum repo stars." default:"500"`
	Format   string `help:"Output format." enum:"text,json" default:"text"`
}

func (cmd *SearchCmd) Run() error {
	printBanner()

	// get slice of recent prsmain
	prList, err := fetchRecentPulls(SearchOptions{
		Days:     cmd.Days,
		Language: strings.ToLower(cmd.Language),
		Limit:    cmd.Limit,
		MinStars: cmd.MinStars,
		MaxStars: cmd.MaxStars,
	})
	if err != nil {
		return err
	}

	if len(prList) == 0 {
		fmt.Printf("No recent PRs to %s based repositories within the last %d day(s).\n", cmd.Language, cmd.Days)
		return nil
	}

	if cmd.Format == "json" {
		return printJSON(os.Stdout, prList)
	}

	printPullRequests(prList)

	return nil
}

// no subs for now
type VersionCmd struct{}

func (cmd *VersionCmd) Run() error {
	fmt.Println("openstalk some build time version")
	return nil
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("openstalk"),
		kong.Description("Find active open source projects through GitHub pull request activity."),
	)

	if !CLI.NoBanner {
		printBanner()
	}

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
