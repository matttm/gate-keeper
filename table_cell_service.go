package main

import "image/color"

// CellStyle represents the visual styling for a table cell
type CellStyle struct {
	Text  string      // The text content to display
	Color color.Color // The background color of the cell
}

// GetCellStyle calculates the text and color for a table cell based on the column,
// gate data, and timeline linearity.
//
// Parameters:
//   - columnID: The column index (0 for gate name, 1 for position)
//   - gate: The gate data for this row
//   - isTimelineLinear: Whether the timeline is linear (gates don't overlap)
//
// Returns:
//   - CellStyle with the appropriate text and background color
func GetCellStyle(columnID int, gate *Gate, isTimelineLinear bool) CellStyle {
	text := getCellText(columnID, gate)
	cellColor := getCellColor(gate, isTimelineLinear)

	return CellStyle{
		Text:  text,
		Color: cellColor,
	}
}

// getCellText determines the text content for a cell based on its column
func getCellText(columnID int, gate *Gate) string {
	if columnID == 0 {
		return gate.GateName
	}

	if columnID == 1 {
		p := RelativePosition(getGatePosition(gate))
		positions := []string{"Past", "Open", "Future"}
		return positions[p+1]
	}

	return ""
}

// getCellColor determines the background color for a cell
// Color coding:
//   - Yellow (with alpha): Timeline is not linear (gates are out of order)
//   - Green (with alpha): Gate is currently open
//   - Red (with alpha): Gate is closed (past or future)
func getCellColor(gate *Gate, isTimelineLinear bool) color.Color {
	if !isTimelineLinear {
		return color.RGBA{R: 255, G: 255, B: 0, A: 128} // Yellow
	}

	if isGateOpen(gate) {
		return color.RGBA{R: 0, G: 255, B: 0, A: 128} // Green
	}

	return color.RGBA{R: 255, G: 0, B: 0, A: 128} // Red
}
