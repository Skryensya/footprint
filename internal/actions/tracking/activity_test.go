package tracking

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/Skryensya/footprint/internal/dispatchers"
	"github.com/Skryensya/footprint/internal/store"
	_ "github.com/mattn/go-sqlite3"
)

func TestActivity_LimitValidation(t *testing.T) {
	tests := []struct {
		name      string
		flags     []string
		wantError bool
		errorText string
	}{
		{
			name:      "no limit specified",
			flags:     []string{},
			wantError: false,
		},
		{
			name:      "valid limit --limit=5",
			flags:     []string{"--limit=5"},
			wantError: false,
		},
		{
			name:      "valid limit --limit=1",
			flags:     []string{"--limit=1"},
			wantError: false,
		},
		{
			name:      "valid limit --limit=100",
			flags:     []string{"--limit=100"},
			wantError: false,
		},
		{
			name:      "invalid limit --limit=0",
			flags:     []string{"--limit=0"},
			wantError: true,
			errorText: "invalid limit value 0: must be greater than 0",
		},
		{
			name:      "invalid limit --limit=-5",
			flags:     []string{"--limit=-5"},
			wantError: true,
			errorText: "invalid limit value -5: must be greater than 0",
		},
		{
			name:      "invalid limit --limit=abc",
			flags:     []string{"--limit=abc"},
			wantError: true,
			errorText: "invalid limit value 'abc': must be a positive integer",
		},
		{
			name:      "invalid limit --limit=1.5",
			flags:     []string{"--limit=1.5"},
			wantError: true,
			errorText: "invalid limit value '1.5': must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := dispatchers.NewParsedFlags(tt.flags)
			deps := mockDeps()

			err := activity([]string{}, flags, deps)

			if tt.wantError {
				if err == nil {
					t.Errorf("activity() expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorText) {
					t.Errorf("activity() error = %v, want error containing %q", err, tt.errorText)
				}
			} else {
				if err != nil {
					t.Errorf("activity() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestActivity_LimitApplication(t *testing.T) {
	tests := []struct {
		name          string
		flags         []string
		wantLimit     int
		totalEvents   int
	}{
		{
			name:        "no limit - returns all events",
			flags:       []string{},
			wantLimit:   0,
			totalEvents: 10,
		},
		{
			name:        "limit=5 - applies limit",
			flags:       []string{"--limit=5"},
			wantLimit:   5,
			totalEvents: 10,
		},
		{
			name:        "limit=1 - applies limit",
			flags:       []string{"--limit=1"},
			wantLimit:   1,
			totalEvents: 10,
		},
		{
			name:        "limit=100 - applies limit (more than available)",
			flags:       []string{"--limit=100"},
			wantLimit:   100,
			totalEvents: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := dispatchers.NewParsedFlags(tt.flags)

			var capturedFilter store.EventFilter
			deps := mockDepsWithFilterCapture(&capturedFilter, tt.totalEvents)

			err := activity([]string{}, flags, deps)
			if err != nil {
				t.Fatalf("activity() unexpected error = %v", err)
			}

			if capturedFilter.Limit != tt.wantLimit {
				t.Errorf("activity() limit = %d, want %d", capturedFilter.Limit, tt.wantLimit)
			}
		})
	}
}

// mockDeps creates minimal mock dependencies for testing validation
func mockDeps() Deps {
	return Deps{
		DBPath: func() string { return ":memory:" },
		OpenDB: func(path string) (*sql.DB, error) {
			db, err := sql.Open("sqlite3", path)
			if err != nil {
				return nil, err
			}
			return db, nil
		},
		ListEvents: func(*sql.DB, store.EventFilter) ([]store.RepoEvent, error) {
			return []store.RepoEvent{}, nil
		},
		Pager: func(string) {},
	}
}

// mockDepsWithFilterCapture creates mock dependencies that capture the filter for testing
func mockDepsWithFilterCapture(capturedFilter *store.EventFilter, totalEvents int) Deps {
	return Deps{
		DBPath: func() string { return ":memory:" },
		OpenDB: func(path string) (*sql.DB, error) {
			db, err := sql.Open("sqlite3", path)
			if err != nil {
				return nil, err
			}
			return db, nil
		},
		ListEvents: func(_ *sql.DB, filter store.EventFilter) ([]store.RepoEvent, error) {
			*capturedFilter = filter

			// Generate mock events
			events := make([]store.RepoEvent, totalEvents)
			for i := 0; i < totalEvents; i++ {
				events[i] = store.RepoEvent{
					ID:       int64(i + 1),
					RepoID:   "test-repo",
					RepoPath: "/test/path",
					Commit:   "abc123",
					Branch:   "main",
				}
			}

			// Simulate LIMIT in SQL
			if filter.Limit > 0 && filter.Limit < len(events) {
				return events[:filter.Limit], nil
			}

			return events, nil
		},
		Pager: func(string) {},
	}
}
