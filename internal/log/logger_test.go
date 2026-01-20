package log

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogger_BasicLogging(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	logger, err := New(logPath, LevelDebug)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Write messages at different levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warning message")
	logger.Error("error message")

	// Close to flush
	logger.Close()

	// Read log file
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify all messages are present
	if !strings.Contains(logContent, "DEBUG: debug message") {
		t.Error("Debug message not found in log")
	}
	if !strings.Contains(logContent, "INFO: info message") {
		t.Error("Info message not found in log")
	}
	if !strings.Contains(logContent, "WARN: warning message") {
		t.Error("Warning message not found in log")
	}
	if !strings.Contains(logContent, "ERROR: error message") {
		t.Error("Error message not found in log")
	}
}

func TestLogger_LevelFiltering(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	// Create logger with Warn level (should filter out Debug and Info)
	logger, err := New(logPath, LevelWarn)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warning message")
	logger.Error("error message")

	logger.Close()

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Debug and Info should NOT be present
	if strings.Contains(logContent, "DEBUG") {
		t.Error("Debug message should have been filtered")
	}
	if strings.Contains(logContent, "INFO") {
		t.Error("Info message should have been filtered")
	}

	// Warn and Error SHOULD be present
	if !strings.Contains(logContent, "WARN: warning message") {
		t.Error("Warning message should be present")
	}
	if !strings.Contains(logContent, "ERROR: error message") {
		t.Error("Error message should be present")
	}
}

func TestLogger_FilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	logger, err := New(logPath, LevelInfo)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("test message")
	logger.Close()

	// Check file permissions
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}

	mode := info.Mode()
	// File should be readable and writable by owner only (0600)
	expected := os.FileMode(0600)
	if mode.Perm() != expected {
		t.Errorf("Log file permissions = %o, want %o", mode.Perm(), expected)
	}
}

func TestLogger_DirectoryPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	logDir := filepath.Join(tmpDir, "logs")
	logPath := filepath.Join(logDir, "test.log")

	logger, err := New(logPath, LevelInfo)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Check directory permissions
	info, err := os.Stat(logDir)
	if err != nil {
		t.Fatalf("Failed to stat log directory: %v", err)
	}

	mode := info.Mode()
	// Directory should be rwx for owner only (0700)
	expected := os.FileMode(0700) | os.ModeDir
	if mode != expected {
		t.Errorf("Log directory permissions = %o, want %o", mode, expected)
	}
}

func TestLogger_AppendMode(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	// First logger
	logger1, err := New(logPath, LevelInfo)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	logger1.Info("first message")
	logger1.Close()

	// Second logger (should append)
	logger2, err := New(logPath, LevelInfo)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	logger2.Info("second message")
	logger2.Close()

	// Read log file
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Both messages should be present
	if !strings.Contains(logContent, "first message") {
		t.Error("First message not found")
	}
	if !strings.Contains(logContent, "second message") {
		t.Error("Second message not found")
	}
}

func TestLogger_Disabled(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	logger, err := New(logPath, LevelInfo)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("enabled message")
	logger.SetEnabled(false)
	logger.Info("disabled message")
	logger.SetEnabled(true)
	logger.Info("enabled again")

	logger.Close()

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	if !strings.Contains(logContent, "enabled message") {
		t.Error("First message not found")
	}
	if strings.Contains(logContent, "disabled message") {
		t.Error("Disabled message should not be present")
	}
	if !strings.Contains(logContent, "enabled again") {
		t.Error("Third message not found")
	}
}
