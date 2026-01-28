#!/bin/bash
# event-generator.sh - Continuously generates git events for testing watch -i
#
# Usage: ./scripts/event-generator.sh [num_repos]
# Default: 3 repos

set -e

NUM_REPOS=${1:-3}
BASE_DIR="/tmp/fp-event-gen-$$"
REPOS=()

# Get absolute path to fp binary
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
FP_BIN="$(cd "$SCRIPT_DIR/.." && pwd)/fp"

if [ ! -f "$FP_BIN" ]; then
    echo "Error: fp binary not found at $FP_BIN"
    echo "Run 'make build' first."
    exit 1
fi

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MUTED='\033[0;90m'
NC='\033[0m'

# Cleanup on exit
cleanup() {
    echo -e "\n${YELLOW}Cleaning up...${NC}"
    for repo in "${REPOS[@]}"; do
        "$FP_BIN" teardown "$repo" 2>/dev/null || true
    done
    rm -rf "$BASE_DIR"
    echo -e "${GREEN}Done.${NC}"
}
trap cleanup EXIT INT TERM

# Random number between min and max
rand_between() {
    local min=$1
    local max=$2
    echo $((min + RANDOM % (max - min + 1)))
}

# Random sleep with variance (in seconds with decimal)
random_sleep() {
    local min_ms=$1
    local max_ms=$2
    local sleep_ms=$(rand_between $min_ms $max_ms)
    local sleep_sec=$(echo "scale=3; $sleep_ms / 1000" | bc)
    sleep "$sleep_sec"
}

# Generate random file content
random_content() {
    local lines=$(rand_between 5 50)
    for i in $(seq 1 $lines); do
        echo "Line $i: $(date +%s%N) - $RANDOM $RANDOM $RANDOM"
    done
}

# Random file name
random_filename() {
    local exts=("js" "ts" "go" "py" "rb" "rs" "java" "c" "cpp" "h" "md" "json" "yaml" "toml")
    local dirs=("src" "lib" "pkg" "internal" "cmd" "api" "utils" "helpers" "models" "services")
    local names=("main" "utils" "helpers" "config" "server" "client" "handler" "router" "model" "service" "controller" "middleware" "auth" "database" "cache" "queue" "worker" "scheduler" "monitor" "logger")

    local ext=${exts[$RANDOM % ${#exts[@]}]}
    local dir=${dirs[$RANDOM % ${#dirs[@]}]}
    local name=${names[$RANDOM % ${#names[@]}]}

    echo "$dir/${name}_$RANDOM.$ext"
}

# Random commit message
random_message() {
    local verbs=("Add" "Update" "Fix" "Refactor" "Remove" "Implement" "Improve" "Optimize" "Clean" "Document")
    local nouns=("feature" "bug" "tests" "docs" "config" "API" "endpoint" "handler" "middleware" "service" "model" "schema" "migration" "dependency" "build" "CI" "deployment")
    local details=("for better performance" "to fix edge case" "as requested" "for new requirements" "to improve UX" "for security" "to reduce complexity" "for maintainability" "based on feedback" "for consistency")

    local verb=${verbs[$RANDOM % ${#verbs[@]}]}
    local noun=${nouns[$RANDOM % ${#nouns[@]}]}
    local detail=${details[$RANDOM % ${#details[@]}]}

    echo "$verb $noun $detail"
}

# Random branch name
random_branch() {
    local prefixes=("feature" "fix" "hotfix" "refactor" "chore" "docs" "test")
    local names=("auth" "api" "ui" "database" "cache" "logging" "metrics" "config" "deploy" "ci" "security" "performance")

    local prefix=${prefixes[$RANDOM % ${#prefixes[@]}]}
    local name=${names[$RANDOM % ${#names[@]}]}

    echo "$prefix/$name-$RANDOM"
}

# Install hook in a repo
install_hook() {
    local repo=$1
    local hook_name=$2
    local hook_path="$repo/.git/hooks/$hook_name"

    cat > "$hook_path" << EOF
#!/bin/sh
FP_SOURCE='$hook_name' '$FP_BIN' record >/dev/null 2>&1 || true
EOF
    chmod +x "$hook_path"
}

# Create a repository with hooks
create_repo() {
    local name=$1
    local path="$BASE_DIR/$name"

    mkdir -p "$path"
    cd "$path"
    git init -q
    git config user.email "test@example.com"
    git config user.name "Test User"

    # Install all hooks
    install_hook "$path" "post-commit"
    install_hook "$path" "post-merge"
    install_hook "$path" "post-checkout"
    install_hook "$path" "post-rewrite"

    # Initial commit (before tracking so it doesn't get recorded)
    echo "# $name" > README.md
    echo "" >> README.md
    echo "Test repository for footprint event generation." >> README.md
    git add README.md
    git commit -q -m "Initial commit"

    echo "$path"
}

# Event: Simple commit
do_commit() {
    local repo=$1
    cd "$repo"

    local file=$(random_filename)
    mkdir -p "$(dirname "$file")"
    random_content > "$file"
    git add "$file"
    git commit -q -m "$(random_message)"

    echo -e "${GREEN}COMMIT${NC}   $(basename "$repo"): $(basename "$file")"
}

# Event: Modify existing file
do_modify() {
    local repo=$1
    cd "$repo"

    # Find a random existing file
    local files=($(find . -type f \( -name "*.go" -o -name "*.js" -o -name "*.py" -o -name "*.ts" \) 2>/dev/null | grep -v .git | head -20))
    if [ ${#files[@]} -eq 0 ]; then
        do_commit "$repo"
        return
    fi

    local file=${files[$RANDOM % ${#files[@]}]}
    echo "" >> "$file"
    echo "// Modified at $(date)" >> "$file"
    random_content >> "$file"
    git add "$file"
    git commit -q -m "$(random_message)"

    echo -e "${BLUE}MODIFY${NC}   $(basename "$repo"): $(basename "$file")"
}

# Event: Create and merge branch
do_merge() {
    local repo=$1
    cd "$repo"

    local branch=$(random_branch)
    git checkout -q -b "$branch"

    # Make some commits on branch
    local num_commits=$(rand_between 1 3)
    for i in $(seq 1 $num_commits); do
        local file=$(random_filename)
        mkdir -p "$(dirname "$file")"
        random_content > "$file"
        git add "$file"
        git commit -q -m "$(random_message)"
    done

    git checkout -q main 2>/dev/null || git checkout -q master
    git merge -q --no-ff -m "Merge $branch" "$branch" 2>/dev/null || true
    git branch -q -d "$branch" 2>/dev/null || true

    echo -e "${YELLOW}MERGE${NC}    $(basename "$repo"): $branch ($num_commits commits)"
}

# Event: Amend commit (rewrite)
do_amend() {
    local repo=$1
    cd "$repo"

    # First make a commit to amend
    local file=$(random_filename)
    mkdir -p "$(dirname "$file")"
    random_content > "$file"
    git add "$file"
    git commit -q -m "$(random_message)"

    # Now amend it
    echo "" >> "$file"
    echo "// Amended at $(date)" >> "$file"
    git add "$file"
    git commit -q --amend --no-edit

    echo -e "${RED}REWRITE${NC}  $(basename "$repo"): $(basename "$file")"
}

# Event: Checkout branch and back
do_checkout() {
    local repo=$1
    cd "$repo"

    local branch=$(random_branch)
    git checkout -q -b "$branch"

    # Make a commit
    local file=$(random_filename)
    mkdir -p "$(dirname "$file")"
    random_content > "$file"
    git add "$file"
    git commit -q -m "$(random_message)"

    # Switch back
    git checkout -q main 2>/dev/null || git checkout -q master

    echo -e "${CYAN}CHECKOUT${NC} $(basename "$repo"): $branch"
}

# Event: Multiple file commit
do_multi_commit() {
    local repo=$1
    cd "$repo"

    local num_files=$(rand_between 3 8)

    for i in $(seq 1 $num_files); do
        local file=$(random_filename)
        mkdir -p "$(dirname "$file")"
        random_content > "$file"
        git add "$file"
    done

    git commit -q -m "$(random_message)"

    echo -e "${GREEN}MULTI${NC}    $(basename "$repo"): $num_files files"
}

# Event: Delete files
do_delete() {
    local repo=$1
    cd "$repo"

    local files=($(find . -type f \( -name "*.go" -o -name "*.js" -o -name "*.py" \) 2>/dev/null | grep -v .git | head -10))
    if [ ${#files[@]} -lt 2 ]; then
        do_commit "$repo"
        return
    fi

    local file=${files[$RANDOM % ${#files[@]}]}
    git rm -q "$file"
    git commit -q -m "Remove $(basename "$file")"

    echo -e "${RED}DELETE${NC}   $(basename "$repo"): $(basename "$file")"
}

# Main
main() {
    echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║${NC}   ${YELLOW}Footprint Event Generator${NC}            ${BLUE}║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${MUTED}Binary:${NC} $FP_BIN"
    echo -e "${MUTED}Repos:${NC}  $NUM_REPOS in $BASE_DIR"
    echo ""

    mkdir -p "$BASE_DIR"

    # Create and track repos
    local repo_names=("webapp" "api-server" "shared-lib" "infra" "mobile-app" "cli-tools" "data-pipeline" "ml-models")
    for i in $(seq 1 $NUM_REPOS); do
        local name=${repo_names[$((i-1)) % ${#repo_names[@]}]}
        local path=$(create_repo "$name")
        REPOS+=("$path")
        "$FP_BIN" setup "$path" --force 2>/dev/null
        echo -e "${GREEN}✓${NC} Created and tracked: $name"
    done

    echo ""
    echo -e "${YELLOW}Generating events...${NC} (Ctrl+C to stop)"
    echo -e "${MUTED}─────────────────────────────────────────${NC}"
    echo ""

    # Event weights (more common events have higher weights)
    local events=("commit" "commit" "commit" "commit" "modify" "modify" "modify" "multi" "multi" "merge" "amend" "checkout" "delete")
    local event_count=0

    while true; do
        # Pick random repo
        local repo=${REPOS[$RANDOM % ${#REPOS[@]}]}

        # Pick random event type
        local event=${events[$RANDOM % ${#events[@]}]}

        case $event in
            commit)
                do_commit "$repo"
                ;;
            modify)
                do_modify "$repo"
                ;;
            merge)
                do_merge "$repo"
                ;;
            amend)
                do_amend "$repo"
                ;;
            checkout)
                do_checkout "$repo"
                ;;
            multi)
                do_multi_commit "$repo"
                ;;
            delete)
                do_delete "$repo"
                ;;
        esac

        event_count=$((event_count + 1))

        # Random sleep with bursts and slow periods
        local rand_mode=$((RANDOM % 100))
        if [ $rand_mode -lt 10 ]; then
            # Burst mode (10%) - very fast: 50-150ms
            random_sleep 50 150
        elif [ $rand_mode -lt 30 ]; then
            # Slow period (20%) - slow: 1000-2500ms
            random_sleep 1000 2500
        else
            # Normal (70%) - medium: 300-800ms
            random_sleep 300 800
        fi

        # Status every 25 events
        if [ $((event_count % 25)) -eq 0 ]; then
            echo -e "${MUTED}── $event_count events ──${NC}"
        fi
    done
}

main
