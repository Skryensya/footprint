package tracking

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Skryensya/footprint/internal/git"
	"github.com/Skryensya/footprint/internal/store"
	"github.com/stretchr/testify/require"
)

func TestGetCSVForEvent_CurrentYear(t *testing.T) {
	dir := t.TempDir()
	currentYear := 2025
	eventTime := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)

	path, err := getCSVForEvent(dir, eventTime, currentYear)

	require.NoError(t, err)
	require.Equal(t, filepath.Join(dir, "commits.csv"), path)

	// Verify file was created with header
	_, err = os.Stat(path)
	require.NoError(t, err)
}

func TestGetCSVForEvent_PastYear(t *testing.T) {
	dir := t.TempDir()
	currentYear := 2025
	eventTime := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)

	path, err := getCSVForEvent(dir, eventTime, currentYear)

	require.NoError(t, err)
	require.Equal(t, filepath.Join(dir, "commits-2024.csv"), path)

	// Verify file was created with header
	_, err = os.Stat(path)
	require.NoError(t, err)
}

func TestGetCSVForEvent_OlderYear(t *testing.T) {
	dir := t.TempDir()
	currentYear := 2025
	eventTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	path, err := getCSVForEvent(dir, eventTime, currentYear)

	require.NoError(t, err)
	require.Equal(t, filepath.Join(dir, "commits-2023.csv"), path)
}

func TestGetCSVForEvent_ReusesExistingFile(t *testing.T) {
	dir := t.TempDir()
	currentYear := 2025
	eventTime := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)

	// Create file first time
	path1, err := getCSVForEvent(dir, eventTime, currentYear)
	require.NoError(t, err)

	// Get info before second call
	info1, _ := os.Stat(path1)

	// Call again - should return same file without recreating
	path2, err := getCSVForEvent(dir, eventTime, currentYear)
	require.NoError(t, err)
	require.Equal(t, path1, path2)

	// File should not have been truncated/recreated
	info2, _ := os.Stat(path2)
	require.Equal(t, info1.ModTime(), info2.ModTime())
}

func TestWriteCSVHeader_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.csv")

	err := writeCSVHeader(path)

	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(path)
	require.NoError(t, err)
}

func TestWriteCSVHeader_ContainsHeader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.csv")

	err := writeCSVHeader(path)
	require.NoError(t, err)

	// Read file and verify header
	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()

	reader := csv.NewReader(file)
	record, err := reader.Read()
	require.NoError(t, err)

	// Check that header contains expected columns
	require.Contains(t, record, "authored_at")
	require.Contains(t, record, "repo")
	require.Contains(t, record, "commit")
	require.Contains(t, record, "branch")
}

func TestWriteCSVHeader_HasRestrictivePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.csv")

	err := writeCSVHeader(path)
	require.NoError(t, err)

	info, err := os.Stat(path)
	require.NoError(t, err)

	// Check permissions (0600 = rw-------)
	perm := info.Mode().Perm()
	require.Equal(t, os.FileMode(0600), perm)
}

func TestAppendRecord_AppendsData(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.csv")

	// First create file with header
	err := writeCSVHeader(path)
	require.NoError(t, err)

	// Create test event
	event := store.RepoEvent{
		ID:        1,
		RepoID:    "github.com/user/repo",
		Commit:    "abc123def456",
		Branch:    "main",
		Timestamp: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		Status:    store.StatusPending,
		Source:    store.SourcePostCommit,
	}

	meta := git.CommitMetadata{
		AuthoredAt:     "2024-01-15T10:30:00Z",
		ParentCommits:  "parent1",
		AuthorName:     "John Doe",
		AuthorEmail:    "john@example.com",
		CommitterName:  "John Doe",
		CommitterEmail: "john@example.com",
		FilesChanged:   3,
		Insertions:     10,
		Deletions:      5,
		Subject:        "Fix bug",
	}

	// Append record
	err = appendRecord(path, event, meta)
	require.NoError(t, err)

	// Read and verify
	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	require.Len(t, records, 2, "Should have header + 1 data row")

	dataRow := records[1]
	// New column order: authored_at, repo, branch, commit, subject, ...
	require.Contains(t, dataRow[0], "2024-01-15", "Should have authored_at")
	require.Equal(t, "github.com/user/repo", dataRow[1], "Should have repo")
	require.Equal(t, "main", dataRow[2], "Should have branch")
	require.Equal(t, "abc123def456", dataRow[3], "Should have commit")
	require.Equal(t, "Fix bug", dataRow[4], "Should have subject")
}

func TestAppendRecord_SanitizesNewlines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.csv")

	err := writeCSVHeader(path)
	require.NoError(t, err)

	event := store.RepoEvent{
		Timestamp: time.Now().UTC(),
	}

	meta := git.CommitMetadata{
		Subject: "Line 1\nLine 2\rLine 3",
	}

	err = appendRecord(path, event, meta)
	require.NoError(t, err)

	// Read and verify no newlines in message field
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	// The CSV should only have 2 lines: header + data
	lines := 0
	for _, b := range content {
		if b == '\n' {
			lines++
		}
	}
	require.Equal(t, 2, lines, "Should have exactly 2 lines (header + data)")
}

func TestShouldExport_ReturnsNoError(t *testing.T) {
	deps := Deps{
		Now: func() time.Time {
			return time.Now()
		},
	}

	// Test that function doesn't error, regardless of config state
	_, err := shouldExport(deps)

	require.NoError(t, err)
}

func TestEnsureExportRepo_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	exportDir := filepath.Join(dir, "export")

	err := ensureExportRepo(exportDir)

	require.NoError(t, err)

	// Check directory exists
	info, err := os.Stat(exportDir)
	require.NoError(t, err)
	require.True(t, info.IsDir())

	// Check git was initialized
	gitDir := filepath.Join(exportDir, ".git")
	_, err = os.Stat(gitDir)
	require.NoError(t, err)
}

func TestCommitExportChanges_EmptyFiles(t *testing.T) {
	dir := t.TempDir()

	// Should not error with empty file list
	err := commitExportChanges(dir, nil)

	require.NoError(t, err)
}

func TestExportAllEvents_SortsAndGroupsByYear(t *testing.T) {
	dir := t.TempDir()
	exportDir := filepath.Join(dir, "export")
	err := ensureExportRepo(exportDir)
	require.NoError(t, err)

	// Create events from different years (out of order)
	events := []store.RepoEvent{
		{
			ID:        1,
			RepoID:    "github.com/user/repo1",
			RepoPath:  "",
			Commit:    "abc123",
			Branch:    "main",
			Timestamp: time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC),
			Source:    store.SourcePostCommit,
		},
		{
			ID:        2,
			RepoID:    "github.com/user/repo2",
			RepoPath:  "",
			Commit:    "def456",
			Branch:    "main",
			Timestamp: time.Date(2024, 3, 10, 10, 0, 0, 0, time.UTC),
			Source:    store.SourcePostCommit,
		},
		{
			ID:        3,
			RepoID:    "github.com/user/repo1",
			RepoPath:  "",
			Commit:    "ghi789",
			Branch:    "develop",
			Timestamp: time.Date(2025, 1, 5, 10, 0, 0, 0, time.UTC),
			Source:    store.SourcePostCommit,
		},
	}

	deps := Deps{
		Now: func() time.Time {
			return time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
		},
	}

	ids, files, err := exportAllEvents(exportDir, events, deps)

	require.NoError(t, err)
	require.Len(t, ids, 3)
	require.Len(t, files, 2, "Should have 2 files: commits.csv (2025) and commits-2024.csv")

	// Check that both files exist
	_, err = os.Stat(filepath.Join(exportDir, "commits.csv"))
	require.NoError(t, err, "commits.csv should exist for 2025 events")

	_, err = os.Stat(filepath.Join(exportDir, "commits-2024.csv"))
	require.NoError(t, err, "commits-2024.csv should exist for 2024 events")
}

func TestExportAllEvents_MultipleReposInSameFile(t *testing.T) {
	dir := t.TempDir()
	exportDir := filepath.Join(dir, "export")
	err := ensureExportRepo(exportDir)
	require.NoError(t, err)

	// Create events from different repos but same year
	events := []store.RepoEvent{
		{
			ID:        1,
			RepoID:    "github.com/user/repo1",
			Commit:    "abc123",
			Branch:    "main",
			Timestamp: time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC),
			Source:    store.SourcePostCommit,
		},
		{
			ID:        2,
			RepoID:    "github.com/user/repo2",
			Commit:    "def456",
			Branch:    "main",
			Timestamp: time.Date(2025, 6, 16, 10, 0, 0, 0, time.UTC),
			Source:    store.SourcePostCommit,
		},
	}

	deps := Deps{
		Now: func() time.Time {
			return time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
		},
	}

	ids, files, err := exportAllEvents(exportDir, events, deps)

	require.NoError(t, err)
	require.Len(t, ids, 2)
	require.Len(t, files, 1, "Should have 1 file: all events in same year")

	// Read and verify both repos are in the same file
	file, err := os.Open(filepath.Join(exportDir, "commits.csv"))
	require.NoError(t, err)
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	require.Len(t, records, 3, "Should have header + 2 data rows")

	// Check both repos are present
	repos := make(map[string]bool)
	for _, row := range records[1:] {
		repos[row[1]] = true // repo is column 1
	}
	require.True(t, repos["github.com/user/repo1"])
	require.True(t, repos["github.com/user/repo2"])
}

func TestExportAllEvents_EventsAreSortedByTimestamp(t *testing.T) {
	dir := t.TempDir()
	exportDir := filepath.Join(dir, "export")
	err := ensureExportRepo(exportDir)
	require.NoError(t, err)

	// Create events out of order
	events := []store.RepoEvent{
		{
			ID:        3,
			RepoID:    "github.com/user/repo",
			Commit:    "third",
			Branch:    "main",
			Timestamp: time.Date(2025, 6, 20, 10, 0, 0, 0, time.UTC), // Latest
			Source:    store.SourcePostCommit,
		},
		{
			ID:        1,
			RepoID:    "github.com/user/repo",
			Commit:    "first",
			Branch:    "main",
			Timestamp: time.Date(2025, 6, 10, 10, 0, 0, 0, time.UTC), // Earliest
			Source:    store.SourcePostCommit,
		},
		{
			ID:        2,
			RepoID:    "github.com/user/repo",
			Commit:    "second",
			Branch:    "main",
			Timestamp: time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC), // Middle
			Source:    store.SourcePostCommit,
		},
	}

	deps := Deps{
		Now: func() time.Time {
			return time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
		},
	}

	_, _, err = exportAllEvents(exportDir, events, deps)
	require.NoError(t, err)

	// Read CSV and verify order
	file, err := os.Open(filepath.Join(exportDir, "commits.csv"))
	require.NoError(t, err)
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	require.Len(t, records, 4, "Should have header + 3 data rows")

	// Verify commits are in chronological order (oldest first)
	require.Equal(t, "first", records[1][3])  // commit is column 3
	require.Equal(t, "second", records[2][3])
	require.Equal(t, "third", records[3][3])
}

func TestExportAllEvents_EmptyEvents(t *testing.T) {
	dir := t.TempDir()
	exportDir := filepath.Join(dir, "export")
	err := ensureExportRepo(exportDir)
	require.NoError(t, err)

	deps := Deps{
		Now: func() time.Time {
			return time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
		},
	}

	ids, files, err := exportAllEvents(exportDir, []store.RepoEvent{}, deps)

	require.NoError(t, err)
	require.Empty(t, ids)
	require.Empty(t, files)
}

func TestExportAllEvents_YearBoundary(t *testing.T) {
	dir := t.TempDir()
	exportDir := filepath.Join(dir, "export")
	err := ensureExportRepo(exportDir)
	require.NoError(t, err)

	// Events at year boundary
	events := []store.RepoEvent{
		{
			ID:        1,
			RepoID:    "github.com/user/repo",
			Commit:    "last_of_2024",
			Branch:    "main",
			Timestamp: time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			Source:    store.SourcePostCommit,
		},
		{
			ID:        2,
			RepoID:    "github.com/user/repo",
			Commit:    "first_of_2025",
			Branch:    "main",
			Timestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Source:    store.SourcePostCommit,
		},
	}

	deps := Deps{
		Now: func() time.Time {
			return time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
		},
	}

	ids, files, err := exportAllEvents(exportDir, events, deps)

	require.NoError(t, err)
	require.Len(t, ids, 2)
	require.Len(t, files, 2, "Should have 2 files: one for each year")

	// Verify 2024 event is in commits-2024.csv
	file2024, err := os.Open(filepath.Join(exportDir, "commits-2024.csv"))
	require.NoError(t, err)
	defer file2024.Close()

	reader2024 := csv.NewReader(file2024)
	records2024, err := reader2024.ReadAll()
	require.NoError(t, err)
	require.Len(t, records2024, 2)
	require.Equal(t, "last_of_2024", records2024[1][3]) // commit is column 3

	// Verify 2025 event is in commits.csv
	file2025, err := os.Open(filepath.Join(exportDir, "commits.csv"))
	require.NoError(t, err)
	defer file2025.Close()

	reader2025 := csv.NewReader(file2025)
	records2025, err := reader2025.ReadAll()
	require.NoError(t, err)
	require.Len(t, records2025, 2)
	require.Equal(t, "first_of_2025", records2025[1][3]) // commit is column 3
}

func TestExportAllEvents_AppendsToExistingFile(t *testing.T) {
	dir := t.TempDir()
	exportDir := filepath.Join(dir, "export")
	err := ensureExportRepo(exportDir)
	require.NoError(t, err)

	deps := Deps{
		Now: func() time.Time {
			return time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
		},
	}

	// First batch of events
	events1 := []store.RepoEvent{
		{
			ID:        1,
			RepoID:    "github.com/user/repo",
			Commit:    "commit1",
			Branch:    "main",
			Timestamp: time.Date(2025, 6, 10, 10, 0, 0, 0, time.UTC),
			Source:    store.SourcePostCommit,
		},
	}

	_, _, err = exportAllEvents(exportDir, events1, deps)
	require.NoError(t, err)

	// Second batch of events
	events2 := []store.RepoEvent{
		{
			ID:        2,
			RepoID:    "github.com/user/repo",
			Commit:    "commit2",
			Branch:    "main",
			Timestamp: time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC),
			Source:    store.SourcePostCommit,
		},
	}

	_, _, err = exportAllEvents(exportDir, events2, deps)
	require.NoError(t, err)

	// Verify both commits are in the file
	file, err := os.Open(filepath.Join(exportDir, "commits.csv"))
	require.NoError(t, err)
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	require.Len(t, records, 3, "Should have header + 2 data rows from both batches")

	commits := make(map[string]bool)
	for _, row := range records[1:] {
		commits[row[3]] = true // commit is column 3
	}
	require.True(t, commits["commit1"], "First batch commit should be present")
	require.True(t, commits["commit2"], "Second batch commit should be present")
}

func TestGetCSVForEvent_CreatesFileWithHeader(t *testing.T) {
	dir := t.TempDir()
	currentYear := 2025
	eventTime := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)

	path, err := getCSVForEvent(dir, eventTime, currentYear)
	require.NoError(t, err)

	// Verify header was written
	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	require.NoError(t, err)

	require.Equal(t, "authored_at", header[0])
	require.Equal(t, "repo", header[1])
	require.Equal(t, "branch", header[2])
	require.Equal(t, "commit", header[3])
}
