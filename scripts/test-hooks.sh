#!/bin/bash
# Test script that creates a test repo, triggers all hooks, and cleans up
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$SCRIPT_DIR/.."
FP_BIN="$PROJECT_DIR/fp"
TEST_REPO="$PROJECT_DIR/test-fp-repo"

# Verify fp binary exists
if [ ! -x "$FP_BIN" ]; then
    echo "ERROR: fp binary not found at $FP_BIN"
    echo "Run 'make build' first"
    exit 1
fi

# Clean up any existing test repo
rm -rf "$TEST_REPO"
mkdir -p "$TEST_REPO"

cleanup() {
    "$FP_BIN" untrack "$TEST_REPO" 2>/dev/null || true
    rm -rf "$TEST_REPO"
    echo "Cleaned up: $TEST_REPO"
}
trap cleanup EXIT

echo "=== Creating test repo: $TEST_REPO ==="
cd "$TEST_REPO"
git init -q
# Use global git config if available, otherwise set test user
if ! git config user.email &>/dev/null; then
    git config user.email "test@example.com"
    git config user.name "Test User"
fi

echo "=== Setting up fp ==="
"$FP_BIN" setup --force
"$FP_BIN" track

# Verify setup
echo ""
echo "=== Verifying setup ==="
"$FP_BIN" status
"$FP_BIN" check
echo ""

# Verify hooks exist
if [ ! -f ".git/hooks/post-commit" ]; then
    echo "ERROR: post-commit hook not installed"
    exit 1
fi

echo "=== Triggering hooks ==="

# post-commit: 5 commits
echo "[post-commit] 5 commits..."
for i in 1 2 3 4 5; do
    echo "content $i" > "file$i.txt"
    git add .
    git commit -q -m "commit $i"
done

# post-checkout: create and switch branches
echo "[post-checkout] branch switches..."
git checkout -q -b feature-a
git checkout -q -b feature-b
git checkout -q main

# post-commit: commits on branches
echo "[post-commit] commits on branches..."
git checkout -q feature-a
echo "feature a work" > feature-a.txt
git add . && git commit -q -m "feature a work"

git checkout -q feature-b
echo "feature b work" > feature-b.txt
git add . && git commit -q -m "feature b work"

# post-merge: merge branches
echo "[post-merge] merging branches..."
git checkout -q main
git merge -q feature-a -m "merge feature-a"
git merge -q feature-b -m "merge feature-b"

# post-rewrite: amend commits
echo "[post-rewrite] amending commits..."
echo "will be amended" > amended.txt
git add . && git commit -q -m "before amend"
echo "after amend" >> amended.txt
git add . && git commit -q --amend -m "after amend"

# post-rewrite: rebase
echo "[post-rewrite] rebasing..."
git checkout -q -b rebase-branch
echo "rebase 1" > rb1.txt && git add . && git commit -q -m "rebase commit 1"
echo "rebase 2" > rb2.txt && git add . && git commit -q -m "rebase commit 2"
git checkout -q main
echo "main progress" > main-progress.txt && git add . && git commit -q -m "main progress"
git checkout -q rebase-branch
git rebase -q main

# More commits and checkouts to reach ~20 events
echo "[post-commit] more commits..."
git checkout -q main
for i in 6 7 8; do
    echo "extra $i" > "extra$i.txt"
    git add .
    git commit -q -m "extra commit $i"
done

echo "[post-checkout] final switches..."
git checkout -q feature-a
git checkout -q feature-b
git checkout -q main

echo ""
echo "=== Results ==="
EVENT_COUNT=$("$FP_BIN" activity --oneline 2>/dev/null | wc -l | tr -d ' ')
echo "Total events recorded: $EVENT_COUNT"
echo ""
echo "Last 10 events:"
"$FP_BIN" activity --oneline 2>/dev/null | head -10

echo ""
echo "=== Done ==="
