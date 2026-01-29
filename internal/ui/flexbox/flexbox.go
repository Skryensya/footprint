// Package flexbox provides a ratio-based layout system inspired by CSS flexbox.
//
// FlexBox allows you to create flexible layouts where rows and cells
// automatically scale based on their ratio values. This is useful for
// building responsive TUI interfaces.
//
// Example usage:
//
//	fb := flexbox.New(80, 24)
//
//	// Add a header row (ratio 1)
//	header := fb.AddRow(1)
//	header.AddCell(1).SetContent("Header")
//
//	// Add main content row (ratio 4)
//	main := fb.AddRow(4)
//	main.AddCell(1).SetContent("Sidebar")
//	main.AddCell(3).SetGenerator(func(w, h int) string {
//	    return renderContent(w, h)
//	})
//
//	// Render
//	output := fb.Render()
package flexbox

import (
	"github.com/charmbracelet/lipgloss"
)

// FlexBox is the main container that holds rows and manages layout.
type FlexBox struct {
	width    int
	height   int
	rows     []*Row
	style    lipgloss.Style
	hasStyle bool
}

// New creates a new FlexBox with the given dimensions.
func New(width, height int) *FlexBox {
	return &FlexBox{
		width:  width,
		height: height,
	}
}

// SetSize updates the flexbox dimensions.
func (f *FlexBox) SetSize(width, height int) {
	f.width = width
	f.height = height
}

// SetWidth sets the flexbox width.
func (f *FlexBox) SetWidth(width int) {
	f.width = width
}

// SetHeight sets the flexbox height.
func (f *FlexBox) SetHeight(height int) {
	f.height = height
}

// SetStyle sets the lipgloss style for the container.
func (f *FlexBox) SetStyle(style lipgloss.Style) *FlexBox {
	f.style = style
	f.hasStyle = true
	return f
}

// AddRow creates and adds a new row with the given vertical ratio.
// Returns the new row for chaining configuration.
func (f *FlexBox) AddRow(ratioY int) *Row {
	row := NewRow(ratioY)
	f.rows = append(f.rows, row)
	return row
}

// Row returns the row at the given index.
func (f *FlexBox) Row(index int) *Row {
	if index < 0 || index >= len(f.rows) {
		return nil
	}
	return f.rows[index]
}

// RowCount returns the number of rows.
func (f *FlexBox) RowCount() int {
	return len(f.rows)
}

// Clear removes all rows.
func (f *FlexBox) Clear() {
	f.rows = nil
}

// Render renders the flexbox to a string.
func (f *FlexBox) Render() string {
	if len(f.rows) == 0 {
		return ""
	}

	// Calculate row heights based on ratios
	heights := f.calculateRowHeights()

	// Render each row
	var renderedRows []string
	for i, row := range f.rows {
		rowHeight := heights[i]
		rendered := row.Render(f.width, rowHeight)
		renderedRows = append(renderedRows, rendered)
	}

	// Join rows vertically
	result := lipgloss.JoinVertical(lipgloss.Left, renderedRows...)

	// Apply container style if set
	if f.hasStyle {
		result = f.style.Render(result)
	}

	return result
}

// calculateRowHeights distributes height among rows based on their ratios.
func (f *FlexBox) calculateRowHeights() []int {
	if len(f.rows) == 0 {
		return nil
	}

	// Calculate total ratio
	totalRatio := 0
	for _, row := range f.rows {
		totalRatio += row.ratioY
	}

	if totalRatio == 0 {
		// Equal distribution if no ratios set
		height := f.height / len(f.rows)
		heights := make([]int, len(f.rows))
		for i := range heights {
			heights[i] = height
		}
		// Give remaining pixels to last row
		heights[len(heights)-1] += f.height - (height * len(f.rows))
		return heights
	}

	// Calculate heights proportionally
	heights := make([]int, len(f.rows))
	usedHeight := 0

	for i, row := range f.rows {
		// Calculate proportional height
		height := (row.ratioY * f.height) / totalRatio

		// Apply min/max constraints
		if row.minHeight > 0 && height < row.minHeight {
			height = row.minHeight
		}
		if row.maxHeight > 0 && height > row.maxHeight {
			height = row.maxHeight
		}

		heights[i] = height
		usedHeight += height
	}

	// Distribute remaining space to last row
	if usedHeight < f.height {
		heights[len(heights)-1] += f.height - usedHeight
	}

	return heights
}

// GetWidth returns the current width.
func (f *FlexBox) GetWidth() int {
	return f.width
}

// GetHeight returns the current height.
func (f *FlexBox) GetHeight() int {
	return f.height
}

// GetRowHeights returns the calculated heights for each row.
// Useful for debugging or when you need to know actual dimensions.
func (f *FlexBox) GetRowHeights() []int {
	return f.calculateRowHeights()
}

// GetCellWidths returns the calculated widths for cells in a specific row.
func (f *FlexBox) GetCellWidths(rowIndex int) []int {
	if rowIndex < 0 || rowIndex >= len(f.rows) {
		return nil
	}
	return f.rows[rowIndex].calculateCellWidths(f.width)
}
