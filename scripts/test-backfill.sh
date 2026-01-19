#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
FP="$PROJECT_DIR/fp"

# Create temp directory
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

echo -e "${YELLOW}=== Test: fp backfill ===${NC}"
echo "Temp directory: $TEMP_DIR"

# Create a test repository
TEST_REPO="$TEMP_DIR/test-repo"
mkdir -p "$TEST_REPO"
cd "$TEST_REPO"

echo -e "\n${YELLOW}1. Creating test repository with commits...${NC}"
git init -q
git config user.email "test@example.com"
git config user.name "Test User"

# Create commits with different dates
echo "file1" > file1.txt
git add file1.txt
GIT_AUTHOR_DATE="2025-01-01T10:00:00" GIT_COMMITTER_DATE="2025-01-01T10:00:00" \
  git commit -q -m "First commit"

echo "file2" > file2.txt
git add file2.txt
GIT_AUTHOR_DATE="2025-01-02T10:00:00" GIT_COMMITTER_DATE="2025-01-02T10:00:00" \
  git commit -q -m "Second commit"

echo "file3" > file3.txt
git add file3.txt
GIT_AUTHOR_DATE="2025-01-03T10:00:00" GIT_COMMITTER_DATE="2025-01-03T10:00:00" \
  git commit -q -m "Third commit"

echo "file4" > file4.txt
git add file4.txt
GIT_AUTHOR_DATE="2025-01-04T10:00:00" GIT_COMMITTER_DATE="2025-01-04T10:00:00" \
  git commit -q -m "Fourth commit"

echo "file5" > file5.txt
git add file5.txt
GIT_AUTHOR_DATE="2025-01-05T10:00:00" GIT_COMMITTER_DATE="2025-01-05T10:00:00" \
  git commit -q -m "Fifth commit"

COMMIT_COUNT=$(git rev-list --count HEAD)
echo "Created $COMMIT_COUNT commits"

echo -e "\n${YELLOW}2. Tracking repository...${NC}"
$FP track .
echo "Repository tracked"

echo -e "\n${YELLOW}3. Testing --dry-run...${NC}"
DRY_RUN_OUTPUT=$($FP backfill --dry-run)
echo "$DRY_RUN_OUTPUT"

if echo "$DRY_RUN_OUTPUT" | grep -q "Found 5 commits"; then
  echo -e "${GREEN}✓ Dry run found all 5 commits${NC}"
else
  echo -e "${RED}✗ Dry run did not find expected commits${NC}"
  exit 1
fi

echo -e "\n${YELLOW}4. Testing --limit flag...${NC}"
LIMIT_OUTPUT=$($FP backfill --dry-run --limit=2)
if echo "$LIMIT_OUTPUT" | grep -q "Found 2 commits"; then
  echo -e "${GREEN}✓ Limit flag works correctly${NC}"
else
  echo -e "${RED}✗ Limit flag did not work${NC}"
  exit 1
fi

echo -e "\n${YELLOW}5. Testing --since flag...${NC}"
SINCE_OUTPUT=$($FP backfill --dry-run --since=2025-01-04)
SINCE_COUNT=$(echo "$SINCE_OUTPUT" | grep -c "^  " || true)
# --since filters commits, so we should get fewer than 5
if [ "$SINCE_COUNT" -lt 5 ] && [ "$SINCE_COUNT" -gt 0 ]; then
  echo -e "${GREEN}✓ Since flag works correctly (found $SINCE_COUNT commits)${NC}"
else
  echo -e "${RED}✗ Since flag did not work (expected 1-4, got $SINCE_COUNT)${NC}"
  exit 1
fi

echo -e "\n${YELLOW}6. Running actual backfill...${NC}"
$FP backfill --background
sleep 1

echo -e "\n${YELLOW}7. Verifying events were inserted...${NC}"
# Filter by the test repo path to avoid counting events from other repos
ACTIVITY_OUTPUT=$($FP activity --source=backfill --oneline | grep "$TEST_REPO" || true)
echo "$ACTIVITY_OUTPUT"

BACKFILL_COUNT=$(echo "$ACTIVITY_OUTPUT" | grep -c "BACKFILL" || true)
if [ "$BACKFILL_COUNT" -eq 5 ]; then
  echo -e "${GREEN}✓ All 5 commits were backfilled${NC}"
else
  echo -e "${RED}✗ Expected 5 backfilled events, got $BACKFILL_COUNT${NC}"
  exit 1
fi

echo -e "\n${YELLOW}8. Testing idempotency (run backfill again)...${NC}"
$FP backfill --background
sleep 1

ACTIVITY_OUTPUT_2=$($FP activity --source=backfill --oneline | grep "$TEST_REPO" || true)
BACKFILL_COUNT_2=$(echo "$ACTIVITY_OUTPUT_2" | grep -c "BACKFILL" || true)
if [ "$BACKFILL_COUNT_2" -eq 5 ]; then
  echo -e "${GREEN}✓ Idempotent: still 5 events (no duplicates)${NC}"
else
  echo -e "${RED}✗ Duplicates created: expected 5, got $BACKFILL_COUNT_2${NC}"
  exit 1
fi

echo -e "\n${YELLOW}9. Cleaning up tracking...${NC}"
$FP untrack .

echo -e "\n${GREEN}=== All backfill tests passed! ===${NC}"
