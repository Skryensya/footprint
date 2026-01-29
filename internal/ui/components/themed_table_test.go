package components

import (
	"testing"

	"github.com/charmbracelet/bubbles/table"
)

func TestThemedTableSorting(t *testing.T) {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Count", Width: 10},
	}
	rows := []table.Row{
		{"banana", "5"},
		{"apple", "10"},
		{"cherry", "3"},
	}

	tbl := NewThemedTable(columns, rows, 40, 10)

	// Sort by name ascending
	tbl.SortBy(0, true)

	sorted := tbl.Rows()
	if sorted[0][0] != "apple" {
		t.Errorf("Expected 'apple' first, got %q", sorted[0][0])
	}
	if sorted[1][0] != "banana" {
		t.Errorf("Expected 'banana' second, got %q", sorted[1][0])
	}
	if sorted[2][0] != "cherry" {
		t.Errorf("Expected 'cherry' third, got %q", sorted[2][0])
	}

	// Sort by name descending
	tbl.SortBy(0, false)
	sorted = tbl.Rows()
	if sorted[0][0] != "cherry" {
		t.Errorf("Expected 'cherry' first descending, got %q", sorted[0][0])
	}
}

func TestThemedTableSortingInt(t *testing.T) {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Count", Width: 10},
	}
	rows := []table.Row{
		{"banana", "5"},
		{"apple", "10"},
		{"cherry", "3"},
	}

	tbl := NewThemedTable(columns, rows, 40, 10)
	tbl.SetColumnTypes([]ColumnType{ColumnTypeString, ColumnTypeInt})

	// Sort by count (int) ascending
	tbl.SortBy(1, true)

	sorted := tbl.Rows()
	if sorted[0][1] != "3" {
		t.Errorf("Expected '3' first, got %q", sorted[0][1])
	}
	if sorted[1][1] != "5" {
		t.Errorf("Expected '5' second, got %q", sorted[1][1])
	}
	if sorted[2][1] != "10" {
		t.Errorf("Expected '10' third, got %q", sorted[2][1])
	}
}

func TestThemedTableToggleSort(t *testing.T) {
	columns := []table.Column{
		{Title: "Name", Width: 20},
	}
	rows := []table.Row{
		{"banana"},
		{"apple"},
		{"cherry"},
	}

	tbl := NewThemedTable(columns, rows, 40, 10)

	// First toggle: ascending
	tbl.ToggleSort(0)
	if !tbl.SortAscending() {
		t.Error("Expected ascending after first toggle")
	}
	if tbl.Rows()[0][0] != "apple" {
		t.Errorf("Expected 'apple' first, got %q", tbl.Rows()[0][0])
	}

	// Second toggle: descending
	tbl.ToggleSort(0)
	if tbl.SortAscending() {
		t.Error("Expected descending after second toggle")
	}
	if tbl.Rows()[0][0] != "cherry" {
		t.Errorf("Expected 'cherry' first descending, got %q", tbl.Rows()[0][0])
	}
}

func TestThemedTableSortIndicator(t *testing.T) {
	columns := []table.Column{
		{Title: "A", Width: 10},
		{Title: "B", Width: 10},
	}
	rows := []table.Row{{"1", "2"}}

	tbl := NewThemedTable(columns, rows, 40, 10)

	// No sort initially
	if ind := tbl.GetSortIndicator(0); ind != "" {
		t.Errorf("Expected empty indicator, got %q", ind)
	}

	// Sort ascending
	tbl.SortBy(0, true)
	if ind := tbl.GetSortIndicator(0); ind != SortIndicatorAsc {
		t.Errorf("Expected %q, got %q", SortIndicatorAsc, ind)
	}
	if ind := tbl.GetSortIndicator(1); ind != "" {
		t.Errorf("Expected empty for non-sorted column, got %q", ind)
	}

	// Sort descending
	tbl.SortBy(0, false)
	if ind := tbl.GetSortIndicator(0); ind != SortIndicatorDesc {
		t.Errorf("Expected %q, got %q", SortIndicatorDesc, ind)
	}
}

func TestThemedTableFiltering(t *testing.T) {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Type", Width: 10},
	}
	rows := []table.Row{
		{"apple", "fruit"},
		{"banana", "fruit"},
		{"carrot", "vegetable"},
		{"daikon", "vegetable"},
	}

	tbl := NewThemedTable(columns, rows, 40, 10)

	// Filter by type
	tbl.SetFilter(1, "fruit")

	if tbl.RowCount() != 2 {
		t.Errorf("Expected 2 filtered rows, got %d", tbl.RowCount())
	}
	if tbl.TotalRowCount() != 4 {
		t.Errorf("Expected 4 total rows, got %d", tbl.TotalRowCount())
	}

	// Clear filter
	tbl.ClearFilter()
	if tbl.RowCount() != 4 {
		t.Errorf("Expected 4 rows after clear, got %d", tbl.RowCount())
	}
}

func TestThemedTableFilterCaseInsensitive(t *testing.T) {
	columns := []table.Column{{Title: "Name", Width: 20}}
	rows := []table.Row{
		{"Apple"},
		{"BANANA"},
		{"cherry"},
	}

	tbl := NewThemedTable(columns, rows, 40, 10)

	// Case-insensitive filter
	tbl.SetFilter(0, "APPLE")
	if tbl.RowCount() != 1 {
		t.Errorf("Expected 1 row for case-insensitive match, got %d", tbl.RowCount())
	}

	tbl.ClearFilter()
	tbl.SetFilter(0, "an")
	if tbl.RowCount() != 1 { // Only BANANA contains "an"
		t.Errorf("Expected 1 row for partial match, got %d", tbl.RowCount())
	}
}

func TestThemedTableSortAndFilter(t *testing.T) {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Count", Width: 10},
	}
	rows := []table.Row{
		{"apple", "10"},
		{"apricot", "5"},
		{"banana", "3"},
		{"avocado", "7"},
	}

	tbl := NewThemedTable(columns, rows, 40, 10)
	tbl.SetColumnTypes([]ColumnType{ColumnTypeString, ColumnTypeInt})

	// Filter first
	tbl.SetFilter(0, "a")

	// Sort by count
	tbl.SortBy(1, true)

	// Should have filtered rows sorted by count
	filtered := tbl.Rows()
	if len(filtered) != 4 { // All contain 'a'
		t.Errorf("Expected 4 filtered rows, got %d", len(filtered))
	}
	if filtered[0][1] != "3" {
		t.Errorf("Expected '3' first after sort, got %q", filtered[0][1])
	}
}

func TestThemedTableSetRowsWithReset(t *testing.T) {
	columns := []table.Column{{Title: "Name", Width: 20}}
	rows := []table.Row{{"apple"}, {"banana"}}

	tbl := NewThemedTable(columns, rows, 40, 10)

	// Apply sort and filter
	tbl.SortBy(0, false)
	tbl.SetFilter(0, "a")

	// Reset with new rows
	newRows := []table.Row{{"cherry"}, {"date"}}
	tbl.SetRowsWithReset(newRows)

	if tbl.RowCount() != 2 {
		t.Errorf("Expected 2 rows after reset, got %d", tbl.RowCount())
	}
	if tbl.FilterQuery() != "" {
		t.Errorf("Expected empty filter after reset, got %q", tbl.FilterQuery())
	}
}
