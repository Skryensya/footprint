package components

import (
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/footprint-tools/cli/internal/ui/style"
)

// ThemedPaginator wraps bubbles/paginator with theme-aware styling.
type ThemedPaginator struct {
	Model paginator.Model
}

// NewThemedPaginator creates a paginator styled with the current theme.
// Uses dot-style pagination by default.
func NewThemedPaginator() ThemedPaginator {
	colors := style.GetColors()

	p := paginator.New()
	p.Type = paginator.Dots
	p.ActiveDot = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.UIActive)).
		Render("●")
	p.InactiveDot = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.UIDim)).
		Render("○")

	return ThemedPaginator{Model: p}
}

// NewThemedPaginatorArabic creates a paginator with arabic numerals (1/10 style).
func NewThemedPaginatorArabic() ThemedPaginator {
	colors := style.GetColors()

	p := paginator.New()
	p.Type = paginator.Arabic
	p.ArabicFormat = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Muted)).
		Render("%d/%d")

	return ThemedPaginator{Model: p}
}

// SetTotalPages sets the total number of pages.
func (t *ThemedPaginator) SetTotalPages(n int) {
	t.Model.SetTotalPages(n)
}

// TotalPages returns the total number of pages.
func (t ThemedPaginator) TotalPages() int {
	return t.Model.TotalPages
}

// Page returns the current page (0-indexed).
func (t ThemedPaginator) Page() int {
	return t.Model.Page
}

// SetPage sets the current page.
func (t *ThemedPaginator) SetPage(n int) {
	t.Model.Page = n
}

// PerPage returns items per page.
func (t ThemedPaginator) PerPage() int {
	return t.Model.PerPage
}

// SetPerPage sets items per page.
func (t *ThemedPaginator) SetPerPage(n int) {
	t.Model.PerPage = n
}

// ItemsOnPage returns the number of items on the current page.
func (t ThemedPaginator) ItemsOnPage(totalItems int) int {
	return t.Model.ItemsOnPage(totalItems)
}

// GetSliceBounds returns start and end indices for the current page.
func (t ThemedPaginator) GetSliceBounds(totalItems int) (int, int) {
	return t.Model.GetSliceBounds(totalItems)
}

// PrevPage moves to the previous page.
func (t *ThemedPaginator) PrevPage() {
	t.Model.PrevPage()
}

// NextPage moves to the next page.
func (t *ThemedPaginator) NextPage() {
	t.Model.NextPage()
}

// OnFirstPage returns true if on the first page.
func (t ThemedPaginator) OnFirstPage() bool {
	return t.Model.OnFirstPage()
}

// OnLastPage returns true if on the last page.
func (t ThemedPaginator) OnLastPage() bool {
	return t.Model.OnLastPage()
}

// Update handles a tea.Msg and returns updated model and command.
func (t ThemedPaginator) Update(msg tea.Msg) (ThemedPaginator, tea.Cmd) {
	var cmd tea.Cmd
	t.Model, cmd = t.Model.Update(msg)
	return t, cmd
}

// View renders the paginator.
func (t ThemedPaginator) View() string {
	return t.Model.View()
}

// RefreshColors updates the paginator styling from the current theme.
func (t *ThemedPaginator) RefreshColors() {
	colors := style.GetColors()

	if t.Model.Type == paginator.Dots {
		t.Model.ActiveDot = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.UIActive)).
			Render("●")
		t.Model.InactiveDot = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.UIDim)).
			Render("○")
	} else {
		t.Model.ArabicFormat = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Muted)).
			Render("%d/%d")
	}
}
