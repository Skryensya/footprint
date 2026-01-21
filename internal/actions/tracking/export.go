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
	exportRepo := deps.GetExportRepo()

	if err := ensureExportRepo(exportRepo); err != nil {
		return 0, false, fmt.Errorf("could not initialize export repo: %w", err)
	}

	// Sync with remote before writing (offline mode: continue if pull fails)
	if deps.HasRemote(exportRepo) {
		if err := deps.PullExportRepo(exportRepo); err != nil {
			log.Warn("export: could not sync with remote, continuing offline: %v", err)
		}
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
	if deps.HasRemote(exportRepo) {
		if err := deps.PushExportRepo(exportRepo); err != nil {
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
// Uses map-based deduplication: new records replace existing ones with same repo:commit.
// Returns the IDs of exported events and the files that were modified.
func exportAllEvents(exportRepo string, events []store.RepoEvent, deps Deps) ([]int64, []string, error) {
	// Build a map of repo paths for metadata enrichment
	repoPaths := make(map[string]string)
	for _, e := range events {
		if e.RepoPath != "" {
			repoPaths[e.RepoID] = e.RepoPath
		}
	}

	currentYear := deps.Now().Year()

	// Group events by target CSV file
	eventsByFile := make(map[string][]store.RepoEvent)
	for _, e := range events {
		csvPath := getCSVPath(exportRepo, e.Timestamp, currentYear)
		eventsByFile[csvPath] = append(eventsByFile[csvPath], e)
	}

	var exportedIDs []int64
	var modifiedFiles []string

	// Process each CSV file
	for csvPath, fileEvents := range eventsByFile {
		// Load existing records into map (repo:commit -> record)
		records := loadCSVRecords(csvPath)

		// Add/replace with new events
		for _, e := range fileEvents {
			var meta git.CommitMetadata
			if repoPath, ok := repoPaths[e.RepoID]; ok {
				meta = git.GetCommitMetadata(repoPath, e.Commit)
			}

			record := buildRecord(e, meta)
			key := e.RepoID + ":" + e.Commit
			records[key] = record

			exportedIDs = append(exportedIDs, e.ID)
		}

		// Write all records sorted by authored_at
		if err := writeCSVSorted(csvPath, records); err != nil {
			return nil, nil, fmt.Errorf("could not write %s: %w", csvPath, err)
		}

		relPath, _ := filepath.Rel(exportRepo, csvPath)
		modifiedFiles = append(modifiedFiles, relPath)
	}

	return exportedIDs, modifiedFiles, nil
}

// getCSVPath returns the path to the CSV file for an event based on its year.
func getCSVPath(exportRepo string, eventTime time.Time, currentYear int) string {
	eventYear := eventTime.Year()
	var csvName string
	if eventYear == currentYear {
		csvName = activeCSVName
	} else {
		csvName = fmt.Sprintf("commits-%d.csv", eventYear)
	}
	return filepath.Join(exportRepo, csvName)
}

// loadCSVRecords loads existing CSV into a map keyed by repo:commit.
func loadCSVRecords(csvPath string) map[string][]string {
	records := make(map[string][]string)

	file, err := os.Open(csvPath)
	if err != nil {
		return records // File doesn't exist yet, return empty map
	}
	defer file.Close()

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		return records
	}

	// Skip header, parse records
	for i, line := range lines {
		if i == 0 || len(line) < 4 {
			continue
		}
		key := line[1] + ":" + line[3] // repo:commit
		records[key] = line
	}

	return records
}

// buildRecord creates a CSV record from an event and its metadata.
func buildRecord(e store.RepoEvent, meta git.CommitMetadata) []string {
	subject := strings.ReplaceAll(meta.Subject, "\n", " ")
	subject = strings.ReplaceAll(subject, "\r", "")
	subject = strings.TrimSpace(subject)

	// Use event timestamp as fallback if git metadata not available
	authoredAt := meta.AuthoredAt
	if authoredAt == "" {
		authoredAt = e.Timestamp.UTC().Format(time.RFC3339)
	}

	return []string{
		authoredAt,
		e.RepoID,
		e.Branch,
		e.Commit,
		subject,
		meta.AuthorName,
		meta.AuthorEmail,
		strconv.Itoa(meta.FilesChanged),
		strconv.Itoa(meta.Insertions),
		strconv.Itoa(meta.Deletions),
		meta.ParentCommits,
		meta.CommitterName,
		meta.CommitterEmail,
		e.Timestamp.UTC().Format(time.RFC3339),
		e.Source.String(),
		getHostname(),
	}
}

// writeCSVSorted writes all records to CSV, sorted by authored_at (column 0).
func writeCSVSorted(csvPath string, records map[string][]string) error {
	// Collect and sort records by authored_at
	var lines [][]string
	for _, record := range records {
		lines = append(lines, record)
	}
	sort.Slice(lines, func(i, j int) bool {
		return lines[i][0] < lines[j][0] // authored_at is RFC3339, sorts correctly
	})

	// Write file
	file, err := os.OpenFile(csvPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	if err := w.Write(csvHeader); err != nil {
		return err
	}
	for _, line := range lines {
		if err := w.Write(line); err != nil {
			return err
		}
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

// pullExportRepo pulls changes from remote before writing.
// Handles three scenarios:
// 1. Remote is empty (first push) - skip pull
// 2. Normal divergence - rebase local commits on top of remote
// 3. Unrelated histories - merge with automatic CSV conflict resolution
func pullExportRepo(exportRepo string) error {
	// First, fetch to see what's on the remote
	fetchCmd := exec.Command("git", "fetch", "origin")
	fetchCmd.Dir = exportRepo
	if err := fetchCmd.Run(); err != nil {
		return fmt.Errorf("git fetch failed: %w", err)
	}

	// Check if remote has any branches (empty remote = first push)
	checkCmd := exec.Command("git", "branch", "-r")
	checkCmd.Dir = exportRepo
	output, _ := checkCmd.Output()
	if len(strings.TrimSpace(string(output))) == 0 {
		log.Debug("export: remote is empty, skipping pull")
		return nil
	}

	// Try normal pull with rebase first
	pullCmd := exec.Command("git", "pull", "--rebase", "origin", "HEAD")
	pullCmd.Dir = exportRepo
	pullOutput, err := pullCmd.CombinedOutput()
	if err == nil {
		return nil
	}

	pullOutputStr := string(pullOutput)

	// Check if it's an unrelated histories error
	if strings.Contains(pullOutputStr, "unrelated histories") {
		log.Info("export: detected unrelated histories, attempting merge")
		return mergeUnrelatedHistories(exportRepo)
	}

	return fmt.Errorf("git pull --rebase failed: %s", strings.TrimSpace(pullOutputStr))
}

// mergeUnrelatedHistories handles the case where local and remote repos
// started independently. Merges them and auto-resolves CSV conflicts.
func mergeUnrelatedHistories(exportRepo string) error {
	// Get the remote branch name
	remoteBranch := "origin/HEAD"

	// Try to merge with allow-unrelated-histories
	mergeCmd := exec.Command("git", "merge", remoteBranch, "--allow-unrelated-histories", "--no-edit")
	mergeCmd.Dir = exportRepo
	mergeOutput, err := mergeCmd.CombinedOutput()

	if err == nil {
		log.Info("export: successfully merged unrelated histories")
		return nil
	}

	// Check if there are conflicts
	if !strings.Contains(string(mergeOutput), "CONFLICT") {
		return fmt.Errorf("merge failed: %s", strings.TrimSpace(string(mergeOutput)))
	}

	// There are conflicts - try to auto-resolve CSV files
	log.Info("export: resolving CSV conflicts automatically")
	if err := resolveCSVConflicts(exportRepo); err != nil {
		// Abort the merge if we can't resolve
		abortCmd := exec.Command("git", "merge", "--abort")
		abortCmd.Dir = exportRepo
		_ = abortCmd.Run()
		return fmt.Errorf("could not resolve conflicts: %w", err)
	}

	// Commit the resolved merge
	commitCmd := exec.Command("git", "commit", "--no-edit")
	commitCmd.Dir = exportRepo
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("could not commit merge: %w", err)
	}

	log.Info("export: successfully consolidated local and remote histories")
	return nil
}

// resolveCSVConflicts auto-resolves conflicts in CSV files by combining
// both versions and sorting by date. Only works for append-only CSVs.
func resolveCSVConflicts(exportRepo string) error {
	// Get list of conflicted files
	statusCmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	statusCmd.Dir = exportRepo
	output, err := statusCmd.Output()
	if err != nil {
		return fmt.Errorf("could not get conflicted files: %w", err)
	}

	conflictedFiles := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, file := range conflictedFiles {
		if file == "" {
			continue
		}
		if !strings.HasSuffix(file, ".csv") {
			return fmt.Errorf("non-CSV conflict in %s, manual resolution required", file)
		}

		filePath := filepath.Join(exportRepo, file)
		if err := resolveCSVFile(exportRepo, filePath); err != nil {
			return fmt.Errorf("could not resolve %s: %w", file, err)
		}

		// Stage the resolved file
		addCmd := exec.Command("git", "add", file)
		addCmd.Dir = exportRepo
		if err := addCmd.Run(); err != nil {
			return fmt.Errorf("could not stage %s: %w", file, err)
		}
	}

	return nil
}

// resolveCSVFile resolves a conflicted CSV by combining both versions.
// "theirs" (incoming remote) replaces "ours" (local) for duplicates.
// Result is sorted by authored_at.
func resolveCSVFile(exportRepo, filePath string) error {
	// Get "ours" version (local) and "theirs" version (remote)
	oursCmd := exec.Command("git", "show", ":2:"+filepath.Base(filePath))
	oursCmd.Dir = exportRepo
	oursOutput, _ := oursCmd.Output()

	theirsCmd := exec.Command("git", "show", ":3:"+filepath.Base(filePath))
	theirsCmd.Dir = exportRepo
	theirsOutput, _ := theirsCmd.Output()

	// Parse both CSVs into maps
	records := make(map[string][]string) // repo:commit -> record

	// Load ours first
	parseCSVIntoMap(string(oursOutput), records)
	// Load theirs second - will replace duplicates (incoming wins)
	parseCSVIntoMap(string(theirsOutput), records)

	// Write using shared function
	return writeCSVSorted(filePath, records)
}

// parseCSVIntoMap parses CSV content and adds records to the map.
// Later calls overwrite earlier entries (last write wins).
func parseCSVIntoMap(content string, records map[string][]string) {
	r := csv.NewReader(strings.NewReader(content))
	lines, err := r.ReadAll()
	if err != nil {
		return
	}

	for i, line := range lines {
		if i == 0 || len(line) < 4 {
			continue // skip header
		}
		key := line[1] + ":" + line[3] // repo:commit
		records[key] = line
	}
}

// pushExportRepo pushes the export repository to its remote.
func pushExportRepo(exportRepo string) error {
	// Push to origin, set upstream if needed
	cmd := exec.Command("git", "push", "-u", "origin", "HEAD")
	cmd.Dir = exportRepo
	return cmd.Run()
}
