package components

import (
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/footprint-tools/cli/internal/ui/style"
)

// ColumnType defines how a column's values should be sorted and compared.
type ColumnType int

const (
	// ColumnTypeString sorts values alphabetically (default).
	ColumnTypeString ColumnType = iota
	// ColumnTypeInt sorts values as integers.
	ColumnTypeInt
	// ColumnTypeFloat sorts values as floating-point numbers.
	ColumnTypeFloat
)

// SortIndicatorAsc is shown next to ascending sort columns.
const SortIndicatorAsc = "▲"

// SortIndicatorDesc is shown next to descending sort columns.
const SortIndicatorDesc = "▼"

// ThemedTable wraps bubbles/table with theme-aware styling,
// plus sorting, filtering, and enhanced navigation capabilities.
type ThemedTable struct {
	Model table.Model

	// Sorting state
	sortColumn    int
	sortAscending bool
	columnTypes   []ColumnType
	originalRows  []table.Row // Rows before any filtering

	// Filtering state
	filterColumn int
	filterQuery  string
}

// NewThemedTable creates a table styled with the current theme.
func NewThemedTable(columns []table.Column, rows []table.Row, width, height int) ThemedTable {
	colors := style.GetColors()

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithWidth(width),
		table.WithHeight(height),
	)

	// Apply theme-based styles
	s := table.DefaultStyles()

	// Header style - bold with background
	s.Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color(colors.Info)).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color(colors.Border))

	// Cell style - normal rows
	s.Cell = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color(colors.Muted))

	// Selected row style - highlighted background
	s.Selected = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color(colors.UIActive)).
		Padding(0, 1)

	t.SetStyles(s)

	return ThemedTable{Model: t, sortColumn: -1}
}

// NewThemedTableSimple creates a table with minimal configuration.
func NewThemedTableSimple(columns []table.Column, width, height int) ThemedTable {
	return NewThemedTable(columns, nil, width, height)
}

// SetRows sets the table rows.
func (t *ThemedTable) SetRows(rows []table.Row) {
	t.Model.SetRows(rows)
}

// SetColumns sets the table columns.
func (t *ThemedTable) SetColumns(columns []table.Column) {
	t.Model.SetColumns(columns)
}

// SetWidth sets the table width.
func (t *ThemedTable) SetWidth(w int) {
	t.Model.SetWidth(w)
}

// SetHeight sets the table height.
func (t *ThemedTable) SetHeight(h int) {
	t.Model.SetHeight(h)
}

// Focus sets the table to focused state.
func (t *ThemedTable) Focus() {
	t.Model.Focus()
}

// Blur removes focus from the table.
func (t *ThemedTable) Blur() {
	t.Model.Blur()
}

// Focused returns whether the table is focused.
func (t ThemedTable) Focused() bool {
	return t.Model.Focused()
}

// SelectedRow returns the currently selected row.
func (t ThemedTable) SelectedRow() table.Row {
	return t.Model.SelectedRow()
}

// Cursor returns the current cursor position.
func (t ThemedTable) Cursor() int {
	return t.Model.Cursor()
}

// SetCursor sets the cursor position.
func (t *ThemedTable) SetCursor(n int) {
	t.Model.SetCursor(n)
}

// MoveUp moves the cursor up by n rows.
func (t *ThemedTable) MoveUp(n int) {
	t.Model.MoveUp(n)
}

// MoveDown moves the cursor down by n rows.
func (t *ThemedTable) MoveDown(n int) {
	t.Model.MoveDown(n)
}

// GotoTop moves the cursor to the first row.
func (t *ThemedTable) GotoTop() {
	t.Model.GotoTop()
}

// GotoBottom moves the cursor to the last row.
func (t *ThemedTable) GotoBottom() {
	t.Model.GotoBottom()
}

// Rows returns all rows.
func (t ThemedTable) Rows() []table.Row {
	return t.Model.Rows()
}

// Update handles a tea.Msg and returns updated model and command.
func (t ThemedTable) Update(msg tea.Msg) (ThemedTable, tea.Cmd) {
	var cmd tea.Cmd
	t.Model, cmd = t.Model.Update(msg)
	return t, cmd
}

// View renders the table.
func (t ThemedTable) View() string {
	return t.Model.View()
}

// RefreshColors updates the table styling from the current theme.
func (t *ThemedTable) RefreshColors() {
	colors := style.GetColors()
	s := table.DefaultStyles()

	s.Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color(colors.Info)).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color(colors.Border))

	s.Cell = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color(colors.Muted))

	s.Selected = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color(colors.UIActive)).
		Padding(0, 1)

	t.Model.SetStyles(s)
}

// SetStyles sets custom styles for the table.
func (t *ThemedTable) SetStyles(s table.Styles) {
	t.Model.SetStyles(s)
}

// --- Sorting ---

// SetColumnTypes sets the data type for each column, which affects sorting behavior.
// If not set, all columns default to ColumnTypeString.
func (t *ThemedTable) SetColumnTypes(types []ColumnType) {
	t.columnTypes = types
}

// SortBy sorts the table rows by the specified column.
// Column indices are 0-based. If the same column is sorted again,
// the sort direction is toggled.
func (t *ThemedTable) SortBy(column int, ascending bool) {
	t.sortColumn = column
	t.sortAscending = ascending
	t.applySort()
}

// ToggleSort toggles sort direction if the column is already sorted,
// otherwise sorts ascending by the new column.
func (t *ThemedTable) ToggleSort(column int) {
	if t.sortColumn == column {
		t.sortAscending = !t.sortAscending
	} else {
		t.sortColumn = column
		t.sortAscending = true
	}
	t.applySort()
}

// GetSortIndicator returns the sort indicator for a column header.
// Returns SortIndicatorAsc, SortIndicatorDesc, or empty string.
func (t *ThemedTable) GetSortIndicator(column int) string {
	if t.sortColumn != column {
		return ""
	}
	if t.sortAscending {
		return SortIndicatorAsc
	}
	return SortIndicatorDesc
}

// SortColumn returns the currently sorted column index.
func (t *ThemedTable) SortColumn() int {
	return t.sortColumn
}

// SortAscending returns true if the current sort is ascending.
func (t *ThemedTable) SortAscending() bool {
	return t.sortAscending
}

// applySort sorts the current rows based on sortColumn and sortAscending.
func (t *ThemedTable) applySort() {
	rows := t.getSourceRows()
	if len(rows) == 0 {
		return
	}

	colType := ColumnTypeString
	if t.sortColumn < len(t.columnTypes) {
		colType = t.columnTypes[t.sortColumn]
	}

	sorted := make([]table.Row, len(rows))
	copy(sorted, rows)

	sort.SliceStable(sorted, func(i, j int) bool {
		a := ""
		b := ""
		if t.sortColumn < len(sorted[i]) {
			a = sorted[i][t.sortColumn]
		}
		if t.sortColumn < len(sorted[j]) {
			b = sorted[j][t.sortColumn]
		}

		var less bool
		switch colType {
		case ColumnTypeInt:
			ai, _ := strconv.ParseInt(strings.TrimSpace(a), 10, 64)
			bi, _ := strconv.ParseInt(strings.TrimSpace(b), 10, 64)
			less = ai < bi
		case ColumnTypeFloat:
			af, _ := strconv.ParseFloat(strings.TrimSpace(a), 64)
			bf, _ := strconv.ParseFloat(strings.TrimSpace(b), 64)
			less = af < bf
		default:
			less = strings.ToLower(a) < strings.ToLower(b)
		}

		if !t.sortAscending {
			less = !less
		}
		return less
	})

	// Apply filter if active, then set rows
	if t.filterQuery != "" {
		sorted = t.filterRows(sorted)
	}
	t.Model.SetRows(sorted)
}

// --- Filtering ---

// SetFilter filters rows by checking if the specified column contains the query.
// The filter is case-insensitive.
func (t *ThemedTable) SetFilter(column int, query string) {
	// Store original rows on first filter
	if t.filterQuery == "" && query != "" {
		t.originalRows = t.Model.Rows()
	}

	t.filterColumn = column
	t.filterQuery = strings.ToLower(query)
	t.applyFilter()
}

// ClearFilter removes any active filter.
func (t *ThemedTable) ClearFilter() {
	t.filterQuery = ""
	t.filterColumn = 0

	// Restore original rows and re-apply sort
	if t.originalRows != nil {
		t.Model.SetRows(t.originalRows)
		t.originalRows = nil
		// Re-apply sort if active
		if t.sortColumn >= 0 {
			t.applySort()
		}
	}
}

// FilterQuery returns the current filter query.
func (t *ThemedTable) FilterQuery() string {
	return t.filterQuery
}

// FilterColumn returns the current filter column.
func (t *ThemedTable) FilterColumn() int {
	return t.filterColumn
}

// applyFilter filters rows based on filterColumn and filterQuery.
func (t *ThemedTable) applyFilter() {
	source := t.getSourceRows()
	filtered := t.filterRows(source)
	t.Model.SetRows(filtered)
}

// filterRows filters the given rows based on current filter settings.
func (t *ThemedTable) filterRows(rows []table.Row) []table.Row {
	if t.filterQuery == "" {
		return rows
	}

	var filtered []table.Row
	for _, row := range rows {
		if t.filterColumn < len(row) {
			cellValue := strings.ToLower(row[t.filterColumn])
			if strings.Contains(cellValue, t.filterQuery) {
				filtered = append(filtered, row)
			}
		}
	}
	return filtered
}

// getSourceRows returns the original rows (before filtering) if available,
// otherwise returns the current rows.
func (t *ThemedTable) getSourceRows() []table.Row {
	if t.originalRows != nil {
		return t.originalRows
	}
	return t.Model.Rows()
}

// --- Row Management ---

// SetRowsWithReset sets rows and clears any active sort/filter state.
func (t *ThemedTable) SetRowsWithReset(rows []table.Row) {
	t.originalRows = nil
	t.filterQuery = ""
	t.sortColumn = -1
	t.Model.SetRows(rows)
}

// RowCount returns the number of visible rows (after filtering).
func (t *ThemedTable) RowCount() int {
	return len(t.Model.Rows())
}

// TotalRowCount returns the total number of rows (before filtering).
func (t *ThemedTable) TotalRowCount() int {
	if t.originalRows != nil {
		return len(t.originalRows)
	}
	return len(t.Model.Rows())
}
