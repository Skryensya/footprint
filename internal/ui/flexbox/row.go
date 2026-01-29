package flexbox

import (
	"github.com/charmbracelet/lipgloss"
)

// Row represents a horizontal row containing cells.
// Rows have a vertical ratio that determines their height
// relative to other rows in the flexbox.
type Row struct {
	ratioY    int
	minHeight int
	maxHeight int
	cells     []*Cell
	style     lipgloss.Style
	hasStyle  bool
}

// NewRow creates a new row with the given vertical ratio.
func NewRow(ratioY int) *Row {
	return &Row{
		ratioY:    ratioY,
		minHeight: 1,
	}
}

// SetRatio sets the vertical ratio for this row.
func (r *Row) SetRatio(ratio int) *Row {
	r.ratioY = ratio
	return r
}

// SetMinHeight sets the minimum height for this row.
func (r *Row) SetMinHeight(min int) *Row {
	r.minHeight = min
	return r
}

// SetMaxHeight sets the maximum height for this row.
func (r *Row) SetMaxHeight(max int) *Row {
	r.maxHeight = max
	return r
}

// SetStyle sets the lipgloss style for this row.
func (r *Row) SetStyle(style lipgloss.Style) *Row {
	r.style = style
	r.hasStyle = true
	return r
}

// AddCell creates and adds a new cell to this row.
// Returns the new cell for chaining configuration.
func (r *Row) AddCell(ratioX int) *Cell {
	cell := NewCell(ratioX)
	r.cells = append(r.cells, cell)
	return cell
}

// AddCells creates and adds multiple cells with the given ratios.
// Returns the row for chaining.
func (r *Row) AddCells(ratios ...int) *Row {
	for _, ratio := range ratios {
		r.AddCell(ratio)
	}
	return r
}

// Cell returns the cell at the given index.
func (r *Row) Cell(index int) *Cell {
	if index < 0 || index >= len(r.cells) {
		return nil
	}
	return r.cells[index]
}

// CellCount returns the number of cells in this row.
func (r *Row) CellCount() int {
	return len(r.cells)
}

// Render renders the row to fit the given dimensions.
func (r *Row) Render(width, height int) string {
	if len(r.cells) == 0 {
		return ""
	}

	// Calculate cell widths based on ratios
	widths := r.calculateCellWidths(width)

	// Render each cell
	var renderedCells []string
	for i, cell := range r.cells {
		cellWidth := widths[i]
		rendered := cell.Render(cellWidth, height)
		renderedCells = append(renderedCells, rendered)
	}

	// Join cells horizontally
	result := lipgloss.JoinHorizontal(lipgloss.Top, renderedCells...)

	// Apply row style if set
	if r.hasStyle {
		result = r.style.Render(result)
	}

	return result
}

// calculateCellWidths distributes width among cells based on their ratios.
func (r *Row) calculateCellWidths(totalWidth int) []int {
	if len(r.cells) == 0 {
		return nil
	}

	// Calculate total ratio
	totalRatio := 0
	for _, cell := range r.cells {
		totalRatio += cell.ratioX
	}

	if totalRatio == 0 {
		// Equal distribution if no ratios set
		width := totalWidth / len(r.cells)
		widths := make([]int, len(r.cells))
		for i := range widths {
			widths[i] = width
		}
		// Give remaining pixels to last cell
		widths[len(widths)-1] += totalWidth - (width * len(r.cells))
		return widths
	}

	// Calculate widths proportionally
	widths := make([]int, len(r.cells))
	usedWidth := 0

	for i, cell := range r.cells {
		// Calculate proportional width
		width := (cell.ratioX * totalWidth) / totalRatio

		// Apply min/max constraints
		if cell.minWidth > 0 && width < cell.minWidth {
			width = cell.minWidth
		}
		if cell.maxWidth > 0 && width > cell.maxWidth {
			width = cell.maxWidth
		}

		widths[i] = width
		usedWidth += width
	}

	// Distribute remaining space to last cell
	if usedWidth < totalWidth {
		widths[len(widths)-1] += totalWidth - usedWidth
	}

	return widths
}

// GetRatio returns the vertical ratio.
func (r *Row) GetRatio() int {
	return r.ratioY
}

// GetMinHeight returns the minimum height.
func (r *Row) GetMinHeight() int {
	return r.minHeight
}

// GetMaxHeight returns the maximum height.
func (r *Row) GetMaxHeight() int {
	return r.maxHeight
}
