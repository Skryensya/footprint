package actions

import (
	"fmt"
	"sort"

	"github.com/Skryensya/footprint/internal/repo"
)

func RepoList(args []string, flags []string) error {
	repos, err := repo.ListTracked()
	if err != nil {
		return err
	}

	if len(repos) == 0 {
		fmt.Println("no tracked repositories")
		return nil
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i] < repos[j]
	})

	for _, r := range repos {
		fmt.Println(r)
	}

	return nil
}
