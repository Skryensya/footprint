package tracking

import (
	"sort"
)

func Repos(args []string, flags []string) error {
	return repos(args, flags, DefaultDeps())
}

func repos(_ []string, _ []string, deps Deps) error {
	trackedRepos, err := deps.ListTracked()
	if err != nil {
		return err
	}

	if len(trackedRepos) == 0 {
		deps.Println("no tracked repositories")
		return nil
	}

	sort.Slice(trackedRepos, func(i, j int) bool {
		return trackedRepos[i] < trackedRepos[j]
	})

	for _, repoID := range trackedRepos {
		deps.Println(repoID)
	}

	return nil
}
