package tracking

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Skryensya/footprint/internal/config"
	"github.com/Skryensya/footprint/internal/dispatchers"
	"github.com/Skryensya/footprint/internal/git"
	"github.com/Skryensya/footprint/internal/log"
	"github.com/Skryensya/footprint/internal/store"
)

const (
	// CSV filename constants
	activeCSVName = "commits.csv"
)

// CSV header for enriched export (ordered by importance)
var csvHeader = []string{
	"authored_at",
	"repo",
	"branch",
	"commit",
	"subject",
	"author",
	"author_email",
	"files",
	"additions",
	"deletions",
	"parents",
	"committer",
	"committer_email",
	"committed_at",
	"source",
	"machine",
}

// Export handles the manual `fp export` command.
func Export(args []string, flags *dispatchers.ParsedFlags) error {
	return export(args, flags, DefaultDeps())
}

func export(_ []string, flags *dispatchers.ParsedFlags, deps Deps) error {
	force := flags.Has("--force")
	dryRun := flags.Has("--dry-run")
	openDir := flags.Has("--open")

	exportRepo := getExportRepo()

	// Handle --open flag
	if openDir {
		return openInFileManager(exportRepo)
	}

	db, err := deps.OpenDB(deps.DBPath())
	if err != nil {
		return fmt.Errorf("could not open database: %w", err)
	}
	defer db.Close()

	_ = deps.InitDB(db)

	events, err := store.GetPendingEvents(db)
	if err != nil {
		return fmt.Errorf("could not get pending events: %w", err)
	}

	if len(events) == 0 {
		deps.Println("No pending events to export")
		return nil
	}

	if dryRun {
		deps.Printf("Would export %d events:\n", len(events))
		for _, e := range events {
			deps.Printf("  %.7s %s (%s)\n", e.Commit, e.Branch, e.RepoID)
		}
		return nil
	}

	if !force {
		shouldExp, err := shouldExport(deps)
		if err != nil {
			return err
		}
		if !shouldExp {
			deps.Println("Export interval not reached. Use --force to export anyway.")
			return nil
		}
	}

	count, pushed, err := doExportWork(db, events, deps)
	if err != nil {
		return err
	}

	if count == 0 {
		deps.Println("No events were exported")
		return nil
	}

	deps.Printf("Exported %d events\n", count)
	if pushed {
		deps.Println("Pushed to remote")
	}

	return nil
}

// doExportWork performs the core export workflow: export events to CSV, commit, update DB.
// Returns the count of exported events, whether push succeeded, and any error.
func doExportWork(db *sql.DB, events []store.RepoEvent, deps Deps) (int, bool, error) {
	exportRepo := getExportRepo()

	if err := ensureExportRepo(exportRepo); err != nil {
		return 0, false, fmt.Errorf("could not initialize export repo: %w", err)
	}

	exportedIDs, exportedFiles, err := exportAllEvents(exportRepo, events, deps)
	if err != nil {
		return 0, false, fmt.Errorf("could not export events: %w", err)
	}

	if len(exportedFiles) == 0 {
		return 0, false, nil
	}

	if err := commitExportChanges(exportRepo, exportedFiles); err != nil {
		return 0, false, fmt.Errorf("could not commit export: %w", err)
	}

	if err := store.UpdateEventStatuses(db, exportedIDs, store.StatusExported); err != nil {
		return 0, false, fmt.Errorf("could not update event statuses: %w", err)
	}

	_ = saveExportLast(deps.Now().Unix())

	pushed := false
	if hasRemote(exportRepo) {
		if err := pushExportRepo(exportRepo); err != nil {
			log.Warn("export: failed to push to remote: %v", err)
		} else {
			pushed = true
		}
	}

	return len(exportedIDs), pushed, nil
}

// MaybeExport checks if it's time to export and does so if needed.
func MaybeExport(db *sql.DB, deps Deps) {
	shouldExp, err := shouldExport(deps)
	if err != nil {
		log.Debug("export: shouldExport error: %v", err)
		return
	}
	if !shouldExp {
		log.Debug("export: interval not reached, skipping auto-export")
		return
	}

	events, err := store.GetPendingEvents(db)
	if err != nil {
		log.Error("export: failed to get pending events: %v", err)
		return
	}
	if len(events) == 0 {
		log.Debug("export: no pending events")
		return
	}

	log.Debug("export: auto-exporting %d pending events", len(events))

	count, _, err := doExportWork(db, events, deps)
	if err != nil {
		log.Error("export: %v", err)
		return
	}

	if count == 0 {
		log.Debug("export: no files were exported")
		return
	}

	log.Info("export: auto-exported %d events", count)
}

// exportAllEvents exports all events to a flat CSV structure with year-based rotation.
// Returns the IDs of exported events and the files that were modified.
func exportAllEvents(exportRepo string, events []store.RepoEvent, deps Deps) ([]int64, []string, error) {
	// Sort events by timestamp (oldest first for append-only)
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

	// Build a map of repo paths for metadata enrichment
	repoPaths := make(map[string]string)
	for _, e := range events {
		if e.RepoPath != "" {
			repoPaths[e.RepoID] = e.RepoPath
		}
	}

	var exportedIDs []int64
	var modifiedFiles []string
	modifiedFilesSet := make(map[string]bool)

	currentYear := deps.Now().Year()

	for _, e := range events {
		// Enrich event with Git metadata
		var meta git.CommitMetadata
		if repoPath, ok := repoPaths[e.RepoID]; ok {
			meta = git.GetCommitMetadata(repoPath, e.Commit)
		}

		// Determine target CSV file based on event year
		csvPath, err := getCSVForEvent(exportRepo, e.Timestamp, currentYear)
		if err != nil {
			return nil, nil, fmt.Errorf("could not get CSV for event: %w", err)
		}

		// Append the enriched record
		if err := appendRecord(csvPath, e, meta); err != nil {
			return nil, nil, fmt.Errorf("could not append record: %w", err)
		}

		exportedIDs = append(exportedIDs, e.ID)

		// Track modified files (relative to export repo)
		relPath, _ := filepath.Rel(exportRepo, csvPath)
		if !modifiedFilesSet[relPath] {
			modifiedFilesSet[relPath] = true
			modifiedFiles = append(modifiedFiles, relPath)
		}
	}

	return exportedIDs, modifiedFiles, nil
}

// getCSVForEvent returns the path to the CSV file for an event based on its year.
// Events from the current year go to commits.csv, older events go to commits-{year}.csv
func getCSVForEvent(exportRepo string, eventTime time.Time, currentYear int) (string, error) {
	eventYear := eventTime.Year()

	var csvName string
	if eventYear == currentYear {
		csvName = activeCSVName // commits.csv
	} else {
		csvName = fmt.Sprintf("commits-%d.csv", eventYear)
	}

	csvPath := filepath.Join(exportRepo, csvName)

	// Check if file exists, create with header if not
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		if err := writeCSVHeader(csvPath); err != nil {
			return "", err
		}
	}

	return csvPath, nil
}

// writeCSVHeader writes a new CSV file with just the header row.
func writeCSVHeader(path string) error {
	// Create CSV file with restrictive permissions
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	if err := w.Write(csvHeader); err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}

// appendRecord appends a single enriched record to a CSV file.
func appendRecord(path string, e store.RepoEvent, meta git.CommitMetadata) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)

	// Sanitize commit message: single line, no newlines
	subject := strings.ReplaceAll(meta.Subject, "\n", " ")
	subject = strings.ReplaceAll(subject, "\r", "")
	subject = strings.TrimSpace(subject)

	// Build the record (ordered by importance)
	record := []string{
		meta.AuthoredAt,                        // authored_at
		e.RepoID,                               // repo
		e.Branch,                               // branch
		e.Commit,                               // commit (full hash)
		subject,                                // subject
		meta.AuthorName,                        // author
		meta.AuthorEmail,                       // author_email
		strconv.Itoa(meta.FilesChanged),        // files
		strconv.Itoa(meta.Insertions),          // additions
		strconv.Itoa(meta.Deletions),           // deletions
		meta.ParentCommits,                     // parents (space-separated)
		meta.CommitterName,                     // committer
		meta.CommitterEmail,                    // committer_email
		e.Timestamp.UTC().Format(time.RFC3339), // committed_at
		e.Source.String(),                      // source
		getHostname(),                          // machine
	}

	if err := w.Write(record); err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}

// getHostname returns the machine hostname or empty string if unavailable.
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

func shouldExport(deps Deps) (bool, error) {
	intervalStr, _ := config.Get("export_interval")
	lastExportStr, _ := config.Get("export_last")

	interval, _ := strconv.Atoi(intervalStr)
	lastExport, _ := strconv.ParseInt(lastExportStr, 10, 64)

	now := deps.Now().Unix()
	return (now - lastExport) >= int64(interval), nil
}

func getExportRepo() string {
	value, _ := config.Get("export_repo")
	return value
}

func ensureExportRepo(path string) error {
	// Create export directory with restrictive permissions
	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}

	gitDir := filepath.Join(path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		if err := runGitInDir(path, "init"); err != nil {
			return fmt.Errorf("git init failed: %w", err)
		}
	}

	return nil
}

// commitExportChanges commits all modified files to the export repo.
func commitExportChanges(exportRepo string, files []string) error {
	if len(files) == 0 {
		return nil
	}

	// Add all modified files
	for _, file := range files {
		if err := runGitInDir(exportRepo, "add", file); err != nil {
			return fmt.Errorf("git add failed for %s: %w", file, err)
		}
	}

	// Check if there are changes to commit
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = exportRepo
	if err := cmd.Run(); err == nil {
		// No changes staged, nothing to commit
		return nil
	}

	// Commit with a descriptive message
	msg := fmt.Sprintf("Export %d files", len(files))
	if err := runGitInDir(exportRepo, "commit", "-m", msg); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	return nil
}

func saveExportLast(timestamp int64) error {
	lines, err := config.ReadLines()
	if err != nil {
		return err
	}

	lines, _ = config.Set(lines, "export_last", strconv.FormatInt(timestamp, 10))
	return config.WriteLines(lines)
}

func runGitInDir(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	return cmd.Run()
}

// openInFileManager opens a directory in the system's file manager.
func openInFileManager(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("explorer", path)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", path)
	}

	return cmd.Run()
}

// SetupExportRemote configures the remote URL for the export repository.
// It ensures the export repo exists and sets up the git remote.
// This is called by 'fp config set export_remote <url>'.
func SetupExportRemote(remoteURL string) error {
	exportRepo := getExportRepo()
	if err := ensureExportRepo(exportRepo); err != nil {
		return fmt.Errorf("could not initialize export repo: %w", err)
	}

	// Check if origin already exists
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = exportRepo
	if err := cmd.Run(); err == nil {
		// Origin exists, update it
		return runGitInDir(exportRepo, "remote", "set-url", "origin", remoteURL)
	}
	// Origin doesn't exist, add it
	return runGitInDir(exportRepo, "remote", "add", "origin", remoteURL)
}

// hasRemote checks if the export repository has a remote configured.
func hasRemote(exportRepo string) bool {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = exportRepo
	return cmd.Run() == nil
}

// pushExportRepo pushes the export repository to its remote.
func pushExportRepo(exportRepo string) error {
	// Push to origin, set upstream if needed
	cmd := exec.Command("git", "push", "-u", "origin", "HEAD")
	cmd.Dir = exportRepo
	return cmd.Run()
}
