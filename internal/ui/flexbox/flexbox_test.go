package flexbox

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestFlexBoxBasic(t *testing.T) {
	fb := New(80, 24)

	if fb.GetWidth() != 80 {
		t.Errorf("Width = %d, want 80", fb.GetWidth())
	}
	if fb.GetHeight() != 24 {
		t.Errorf("Height = %d, want 24", fb.GetHeight())
	}
	if fb.RowCount() != 0 {
		t.Errorf("RowCount = %d, want 0", fb.RowCount())
	}
}

func TestFlexBoxAddRow(t *testing.T) {
	fb := New(80, 24)

	row := fb.AddRow(1)
	if row == nil {
		t.Fatal("AddRow returned nil")
	}
	if fb.RowCount() != 1 {
		t.Errorf("RowCount = %d, want 1", fb.RowCount())
	}

	// Check row retrieval
	if fb.Row(0) != row {
		t.Error("Row(0) does not match added row")
	}
	if fb.Row(-1) != nil {
		t.Error("Row(-1) should be nil")
	}
	if fb.Row(100) != nil {
		t.Error("Row(100) should be nil")
	}
}

func TestFlexBoxRowHeights(t *testing.T) {
	fb := New(80, 100)

	fb.AddRow(1) // 25%
	fb.AddRow(3) // 75%

	heights := fb.GetRowHeights()

	if len(heights) != 2 {
		t.Fatalf("Expected 2 heights, got %d", len(heights))
	}

	// Row 1 should be ~25 (1/4 of 100)
	if heights[0] != 25 {
		t.Errorf("First row height = %d, want 25", heights[0])
	}

	// Row 2 gets the remainder
	if heights[1] != 75 {
		t.Errorf("Second row height = %d, want 75", heights[1])
	}
}

func TestFlexBoxCellWidths(t *testing.T) {
	fb := New(100, 24)

	row := fb.AddRow(1)
	row.AddCell(1) // 25%
	row.AddCell(3) // 75%

	widths := fb.GetCellWidths(0)

	if len(widths) != 2 {
		t.Fatalf("Expected 2 widths, got %d", len(widths))
	}

	// Cell 1 should be ~25 (1/4 of 100)
	if widths[0] != 25 {
		t.Errorf("First cell width = %d, want 25", widths[0])
	}

	// Cell 2 gets the remainder
	if widths[1] != 75 {
		t.Errorf("Second cell width = %d, want 75", widths[1])
	}
}

func TestFlexBoxRender(t *testing.T) {
	fb := New(20, 4)

	row1 := fb.AddRow(1)
	row1.AddCell(1).SetContent("Header")

	row2 := fb.AddRow(3)
	row2.AddCell(1).SetContent("Content")

	output := fb.Render()

	if output == "" {
		t.Fatal("Render returned empty string")
	}

	lines := strings.Split(output, "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}
}

func TestFlexBoxGenerator(t *testing.T) {
	fb := New(20, 4)

	row := fb.AddRow(1)
	row.AddCell(1).SetGenerator(func(w, h int) string {
		return strings.Repeat("X", w)
	})

	output := fb.Render()

	if !strings.Contains(output, "XXXX") {
		t.Errorf("Generator output not found in render: %q", output)
	}
}

func TestFlexBoxMinMaxConstraints(t *testing.T) {
	fb := New(100, 100)

	// Row with minHeight
	row := fb.AddRow(1)
	row.SetMinHeight(30)
	fb.AddRow(1)

	heights := fb.GetRowHeights()

	// First row should be at least 30
	if heights[0] < 30 {
		t.Errorf("Row with minHeight=%d has height=%d", 30, heights[0])
	}
}

func TestFlexBoxClear(t *testing.T) {
	fb := New(80, 24)

	fb.AddRow(1)
	fb.AddRow(2)

	if fb.RowCount() != 2 {
		t.Errorf("RowCount before clear = %d, want 2", fb.RowCount())
	}

	fb.Clear()

	if fb.RowCount() != 0 {
		t.Errorf("RowCount after clear = %d, want 0", fb.RowCount())
	}
}

func TestRowAddCells(t *testing.T) {
	row := NewRow(1)
	row.AddCells(1, 2, 3)

	if row.CellCount() != 3 {
		t.Errorf("CellCount = %d, want 3", row.CellCount())
	}

	if row.Cell(0).GetRatio() != 1 {
		t.Errorf("Cell 0 ratio = %d, want 1", row.Cell(0).GetRatio())
	}
	if row.Cell(1).GetRatio() != 2 {
		t.Errorf("Cell 1 ratio = %d, want 2", row.Cell(1).GetRatio())
	}
	if row.Cell(2).GetRatio() != 3 {
		t.Errorf("Cell 2 ratio = %d, want 3", row.Cell(2).GetRatio())
	}
}

func TestCellChaining(t *testing.T) {
	cell := NewCell(1).
		SetRatio(2).
		SetMinWidth(10).
		SetMaxWidth(50).
		SetContent("test")

	if cell.GetRatio() != 2 {
		t.Errorf("Ratio = %d, want 2", cell.GetRatio())
	}
	if cell.GetMinWidth() != 10 {
		t.Errorf("MinWidth = %d, want 10", cell.GetMinWidth())
	}
	if cell.GetMaxWidth() != 50 {
		t.Errorf("MaxWidth = %d, want 50", cell.GetMaxWidth())
	}
}

func TestCellStyle(t *testing.T) {
	// Style application is tested by verifying the cell stores the style
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	cell := NewCell(1).SetStyle(style).SetContent("red")

	// Verify cell has style set
	if !cell.hasStyle {
		t.Error("Expected hasStyle to be true after SetStyle")
	}

	// Render should work without errors
	output := cell.Render(10, 1)
	if output == "" {
		t.Error("Render returned empty string")
	}
	// Note: ANSI codes may not be present when not running in a TTY
}

func TestFlexBoxSetSize(t *testing.T) {
	fb := New(80, 24)

	fb.SetSize(100, 50)
	if fb.GetWidth() != 100 {
		t.Errorf("Width after SetSize = %d, want 100", fb.GetWidth())
	}
	if fb.GetHeight() != 50 {
		t.Errorf("Height after SetSize = %d, want 50", fb.GetHeight())
	}

	fb.SetWidth(120)
	if fb.GetWidth() != 120 {
		t.Errorf("Width after SetWidth = %d, want 120", fb.GetWidth())
	}

	fb.SetHeight(60)
	if fb.GetHeight() != 60 {
		t.Errorf("Height after SetHeight = %d, want 60", fb.GetHeight())
	}
}

func TestCellTruncation(t *testing.T) {
	cell := NewCell(1).SetContent("This is a very long line that should be truncated")

	output := cell.Render(10, 1)
	width := lipgloss.Width(output)

	if width != 10 {
		t.Errorf("Output width = %d, want 10", width)
	}
}

func TestCellMultiline(t *testing.T) {
	cell := NewCell(1).SetContent("Line 1\nLine 2\nLine 3")

	output := cell.Render(10, 5)
	lines := strings.Split(output, "\n")

	if len(lines) != 5 {
		t.Errorf("Expected 5 lines, got %d", len(lines))
	}
}
