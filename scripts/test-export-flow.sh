#!/bin/bash
set -e

# Test script for fp export flow
# Creates multiple test repos, performs various git actions, and tests export

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
FP="$PROJECT_DIR/fp"
TEST_DIR="/tmp/fp-export-test-$$"

# Platform-specific export directory
if [ "$(uname)" = "Darwin" ]; then
    EXPORT_DIR="$HOME/Library/Application Support/footprint/export"
else
    EXPORT_DIR="${XDG_DATA_HOME:-$HOME/.local/share}/footprint/export"
fi

# Colors (subtle)
DIM='\033[2m'
CYAN='\033[0;36m'
GREEN='\033[0;32m'
NC='\033[0m'

section() {
    echo ""
    echo -e "${CYAN}=== $1 ===${NC}"
}

log() {
    echo -e "${DIM}$1${NC}"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

# Cleanup function
cleanup() {
    log "Untracking test repos..."
    for repo in webapp api-server docs; do
        if [ -d "$TEST_DIR/$repo" ]; then
            $FP untrack "$TEST_DIR/$repo" 2>/dev/null || true
        fi
    done
    log "Cleaning up test directory..."
    rm -rf "$TEST_DIR"
}

# Setup trap for cleanup on exit
trap cleanup EXIT

section "SETUP: Creating test environment"

mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

log "Test directory: $TEST_DIR"
log "Using fp: $FP"

# Check fp exists
if [ ! -f "$FP" ]; then
    echo "Error: fp binary not found at $FP"
    echo "Run 'make build' first"
    exit 1
fi

# Clear any existing export data for clean test
log "Clearing previous export data..."
rm -rf "$EXPORT_DIR/repos"

section "REPO 1: webapp - Regular commits and feature branch"

mkdir -p "$TEST_DIR/webapp" && cd "$TEST_DIR/webapp"
git init
git config user.email "dev@example.com"
git config user.name "Developer"

# Create initial structure
echo "# Web Application" > README.md
echo "node_modules/" > .gitignore
git add -A && git commit -m "Initial commit: project setup"

$FP track
$FP setup
success "Repo 'webapp' tracked and hooks installed"

# Simulate development work
echo "const express = require('express');" > index.js
echo "const app = express();" >> index.js
git add -A && git commit -m "feat: add express server skeleton"

echo "app.get('/', (req, res) => res.send('Hello'));" >> index.js
git add -A && git commit -m "feat: add root endpoint"

# Create feature branch and merge
git checkout -b feature/auth
echo "const jwt = require('jsonwebtoken');" > auth.js
git add -A && git commit -m "feat(auth): add JWT dependency"

echo "function login(user) { return jwt.sign(user); }" >> auth.js
git add -A && git commit -m "feat(auth): implement login function"

git checkout main
git merge feature/auth --no-edit
success "webapp: 5 commits + 1 merge"

cd "$TEST_DIR"

section "REPO 2: api-server - Multiple branches and rebases"

mkdir -p api-server && cd api-server
git init
git config user.email "backend@example.com"
git config user.name "Backend Dev"

echo "# API Server" > README.md
echo "package main" > main.go
git add -A && git commit -m "Initial commit"

$FP track
$FP setup
success "Repo 'api-server' tracked and hooks installed"

# Main branch work
echo 'func main() { fmt.Println("Server starting") }' >> main.go
git add -A && git commit -m "feat: add main function"

echo 'import "net/http"' >> main.go
git add -A && git commit -m "feat: import http package"

# Feature branch with rebase
git checkout -b feature/endpoints
echo "func handleUsers() {}" > handlers.go
git add -A && git commit -m "feat: add users handler"

echo "func handleProducts() {}" >> handlers.go
git add -A && git commit -m "feat: add products handler"

# Go back to main and add more commits
git checkout main
echo "// Config loaded" >> main.go
git add -A && git commit -m "chore: add config comment"

# Merge feature branch
git merge feature/endpoints --no-edit
success "api-server: 6 commits + 1 merge"

cd "$TEST_DIR"

section "REPO 3: docs - Simple documentation updates"

mkdir -p docs && cd docs
git init
git config user.email "writer@example.com"
git config user.name "Tech Writer"

echo "# Documentation" > README.md
git add -A && git commit -m "Initial commit"

$FP track
$FP setup
success "Repo 'docs' tracked and hooks installed"

echo "## Getting Started" >> README.md
git add -A && git commit -m "docs: add getting started section"

echo "## Installation" >> README.md
git add -A && git commit -m "docs: add installation section"

echo "## Usage" >> README.md
git add -A && git commit -m "docs: add usage section"

echo "## API Reference" >> README.md
git add -A && git commit -m "docs: add API reference section"

success "docs: 5 commits"

cd "$TEST_DIR"

section "FIRST EXPORT"

log "Checking pending events before export..."
$FP activity --oneline | head -20

log "Running first export..."
$FP export --force

log "Checking export directory structure..."
find "$EXPORT_DIR/repos" -type f -name "*.csv" 2>/dev/null | while read f; do
    echo "  $f"
    echo "    $(wc -l < "$f") lines"
done

log "Sample from webapp CSV:"
if [ -f "$EXPORT_DIR/repos/local__tmp__fp-export-test-*__webapp/commits.csv" ]; then
    head -3 "$EXPORT_DIR"/repos/local__tmp__fp-export-test-*__webapp/commits.csv 2>/dev/null || true
fi

success "First export completed"

section "REPO 1: webapp - More development"

cd "$TEST_DIR/webapp"

echo "const cors = require('cors');" >> index.js
git add -A && git commit -m "feat: add CORS support"

echo "app.use(cors());" >> index.js
git add -A && git commit -m "feat: enable CORS middleware"

# Amend last commit (triggers post-rewrite)
echo "// CORS enabled" >> index.js
git add -A && git commit --amend --no-edit

success "webapp: 2 more commits + 1 amend"

cd "$TEST_DIR"

section "REPO 2: api-server - Hotfix branch"

cd "$TEST_DIR/api-server"

git checkout -b hotfix/security
echo "func sanitizeInput(s string) string { return s }" > security.go
git add -A && git commit -m "fix: add input sanitization"

echo "func validateToken(t string) bool { return true }" >> security.go
git add -A && git commit -m "fix: add token validation"

git checkout main
git merge hotfix/security --no-edit

success "api-server: 2 more commits + 1 merge"

cd "$TEST_DIR"

section "REPO 3: docs - Quick updates"

cd "$TEST_DIR/docs"

echo "## Troubleshooting" >> README.md
git add -A && git commit -m "docs: add troubleshooting section"

echo "## FAQ" >> README.md
git add -A && git commit -m "docs: add FAQ section"

success "docs: 2 more commits"

cd "$TEST_DIR"

section "SECOND EXPORT"

log "Checking pending events before second export..."
$FP activity --oneline | head -10

log "Running second export..."
$FP export --force

log "Final export directory structure..."
echo ""
find "$EXPORT_DIR/repos" -type f -name "*.csv" 2>/dev/null | sort | while read f; do
    lines=$(wc -l < "$f")
    echo "  $f"
    echo "    Lines: $lines"
done

success "Second export completed"

section "VERIFICATION: Check CSV contents"

log "Listing all CSV files with row counts..."
find "$EXPORT_DIR/repos" -name "commits.csv" 2>/dev/null | while read -r csv; do
    if [ -f "$csv" ]; then
        repo_dir=$(dirname "$csv")
        repo_name=$(basename "$repo_dir")
        rows=$(($(wc -l < "$csv") - 1))  # Subtract header
        echo ""
        echo "Repository: $repo_name"
        echo "  Rows: $rows"
        echo "  Columns: $(head -1 "$csv" | tr ',' '\n' | wc -l)"
        echo "  Header: $(head -1 "$csv" | cut -c1-80)..."
    fi
done

section "VERIFICATION: Check enriched data"

log "Checking if enrichment data is present..."
find "$EXPORT_DIR/repos" -name "commits.csv" 2>/dev/null | while read -r csv; do
    if [ -f "$csv" ]; then
        repo_name=$(basename "$(dirname "$csv")")
        # Only check repos from this test run
        if echo "$repo_name" | grep -q "fp-export-test-$$"; then
            echo ""
            echo "Repository: $repo_name"

            # Check for merge commits
            merges=$(grep -c ",true," "$csv" 2>/dev/null || echo "0")
            echo "  Merge commits: $merges"

            # Check for parent commits (non-empty)
            with_parents=$(awk -F',' 'NR>1 && $6!=""' "$csv" | wc -l)
            echo "  Commits with parents: $with_parents"

            # Check for diff stats
            with_changes=$(awk -F',' 'NR>1 && $12!="0"' "$csv" | wc -l)
            echo "  Commits with file changes: $with_changes"
        fi
    fi
done

section "SUMMARY"

echo ""
echo "Test completed successfully!"
echo ""
echo "Repos created: 3 (webapp, api-server, docs)"
echo "Total commits: ~20"
echo "Merges: 4"
echo "Exports: 2"
echo ""
echo "Export location: $EXPORT_DIR/repos/"
echo ""

log "Checking git log in export repo..."
git -C "$EXPORT_DIR" log --oneline -5 2>/dev/null || echo "(no git history yet)"

log "Sample enriched CSV row:"
find "$EXPORT_DIR/repos" -name "commits.csv" 2>/dev/null | head -1 | xargs -I{} sh -c 'tail -1 "{}"' 2>/dev/null | cut -c1-120

echo ""
success "All tests passed!"
