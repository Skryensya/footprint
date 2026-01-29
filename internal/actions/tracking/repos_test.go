package tracking

import (
	"database/sql"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"github.com/footprint-tools/cli/internal/dispatchers"
	"github.com/footprint-tools/cli/internal/store"
	"github.com/footprint-tools/cli/internal/store/migrations"
)

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = db.Close()
	})

	err = migrations.Run(db)
	require.NoError(t, err)

	return store.NewWithDB(db)
}

func TestReposList_Empty(t *testing.T) {
	s := newTestStore(t)
	var printedLines []string

	deps := DefaultDeps()
	deps.DBPath = func() string { return ":memory:" }
	deps.OpenStore = func(_ string) (*store.Store, error) {
		return s, nil
	}
	deps.Println = func(a ...any) (int, error) {
		if len(a) > 0 {
			if s, ok := a[0].(string); ok {
				printedLines = append(printedLines, s)
			}
		}
		return 0, nil
	}

	flags := &dispatchers.ParsedFlags{}
	err := reposList(nil, flags, deps)
	require.NoError(t, err)
	require.Len(t, printedLines, 2)
	require.Contains(t, printedLines[0], "no tracked repositories")
	require.Contains(t, printedLines[1], "fp setup")
}

func TestReposList_WithRepos(t *testing.T) {
	s := newTestStore(t)

	// Add some repos
	require.NoError(t, s.AddRepo("/path/to/repo1"))
	require.NoError(t, s.AddRepo("/path/to/repo2"))

	var printedLines []string

	deps := DefaultDeps()
	deps.DBPath = func() string { return ":memory:" }
	deps.OpenStore = func(_ string) (*store.Store, error) {
		return s, nil
	}
	deps.Println = func(a ...any) (int, error) {
		if len(a) > 0 {
			if s, ok := a[0].(string); ok {
				printedLines = append(printedLines, s)
			}
		}
		return 0, nil
	}

	flags := &dispatchers.ParsedFlags{}
	err := reposList(nil, flags, deps)
	require.NoError(t, err)
	require.Len(t, printedLines, 2)
	require.Equal(t, "/path/to/repo1", printedLines[0])
	require.Equal(t, "/path/to/repo2", printedLines[1])
}

func TestReposList_OpenStoreError(t *testing.T) {
	deps := DefaultDeps()
	deps.DBPath = func() string { return "/invalid/path" }
	deps.OpenStore = func(_ string) (*store.Store, error) {
		return nil, errors.New("failed to open store")
	}
	deps.Println = func(a ...any) (int, error) {
		return 0, nil
	}

	flags := &dispatchers.ParsedFlags{}
	err := reposList(nil, flags, deps)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to open store")
}

func TestReposList_JSON_Empty(t *testing.T) {
	s := newTestStore(t)
	var printedOutput string

	deps := DefaultDeps()
	deps.DBPath = func() string { return ":memory:" }
	deps.OpenStore = func(_ string) (*store.Store, error) {
		return s, nil
	}
	deps.Println = func(a ...any) (int, error) {
		if len(a) > 0 {
			printedOutput, _ = a[0].(string)
		}
		return 0, nil
	}

	flags := dispatchers.NewParsedFlags([]string{"--json"})
	err := reposList(nil, flags, deps)
	require.NoError(t, err)
	require.Equal(t, "[]", printedOutput)
}

func TestReposList_JSON_WithRepos(t *testing.T) {
	s := newTestStore(t)

	// Add some repos
	require.NoError(t, s.AddRepo("/path/to/repo1"))
	require.NoError(t, s.AddRepo("/path/to/repo2"))

	var printedOutput string

	deps := DefaultDeps()
	deps.DBPath = func() string { return ":memory:" }
	deps.OpenStore = func(_ string) (*store.Store, error) {
		return s, nil
	}
	deps.Println = func(a ...any) (int, error) {
		if len(a) > 0 {
			printedOutput, _ = a[0].(string)
		}
		return 0, nil
	}

	flags := dispatchers.NewParsedFlags([]string{"--json"})
	err := reposList(nil, flags, deps)
	require.NoError(t, err)
	require.Contains(t, printedOutput, `"path": "/path/to/repo1"`)
	require.Contains(t, printedOutput, `"path": "/path/to/repo2"`)
}
