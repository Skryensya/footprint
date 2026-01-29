package flexbox

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/footprint-tools/cli/internal/ui/text"
)

// Cell represents a single cell within a row.
// Cells can contain static content or use a generator function
// that receives the available dimensions.
type Cell struct {
	ratioX    int
	minWidth  int
	maxWidth  int
	content   string
	generator func(width, height int) string
	style     lipgloss.Style
	hasStyle  bool
}

// NewCell creates a new cell with the given horizontal ratio.
// The ratio determines how much horizontal space this cell gets
// relative to other cells in the same row.
func NewCell(ratioX int) *Cell {
	return &Cell{
		ratioX:   ratioX,
		minWidth: 1,
	}
}

// SetRatio sets the horizontal ratio for this cell.
func (c *Cell) SetRatio(ratio int) *Cell {
	c.ratioX = ratio
	return c
}

// SetMinWidth sets the minimum width for this cell.
func (c *Cell) SetMinWidth(min int) *Cell {
	c.minWidth = min
	return c
}

// SetMaxWidth sets the maximum width for this cell.
func (c *Cell) SetMaxWidth(max int) *Cell {
	c.maxWidth = max
	return c
}

// SetContent sets static content for the cell.
// When set, the generator function is ignored.
func (c *Cell) SetContent(content string) *Cell {
	c.content = content
	c.generator = nil
	return c
}

// SetGenerator sets a function that generates content based on available dimensions.
// The function receives the cell's computed width and height.
func (c *Cell) SetGenerator(fn func(width, height int) string) *Cell {
	c.generator = fn
	c.content = ""
	return c
}

// SetStyle sets the lipgloss style for this cell.
func (c *Cell) SetStyle(style lipgloss.Style) *Cell {
	c.style = style
	c.hasStyle = true
	return c
}

// Render renders the cell content to fit the given dimensions.
func (c *Cell) Render(width, height int) string {
	var content string

	if c.generator != nil {
		content = c.generator(width, height)
	} else {
		content = c.content
	}

	// Split content into lines and fit to dimensions
	lines := strings.Split(content, "\n")

	// Ensure we have exactly 'height' lines
	result := make([]string, height)
	for i := 0; i < height; i++ {
		if i < len(lines) {
			line := lines[i]
			// Truncate or pad to width
			lineWidth := lipgloss.Width(line)
			if lineWidth > width {
				line = text.TruncateWithEllipsis(line, width)
			} else if lineWidth < width {
				line = text.PadRight(line, width)
			}
			result[i] = line
		} else {
			result[i] = strings.Repeat(" ", width)
		}
	}

	rendered := strings.Join(result, "\n")

	// Apply style if set
	if c.hasStyle {
		rendered = c.style.Render(rendered)
	}

	return rendered
}

// GetRatio returns the horizontal ratio.
func (c *Cell) GetRatio() int {
	return c.ratioX
}

// GetMinWidth returns the minimum width.
func (c *Cell) GetMinWidth() int {
	return c.minWidth
}

// GetMaxWidth returns the maximum width.
func (c *Cell) GetMaxWidth() int {
	return c.maxWidth
}
