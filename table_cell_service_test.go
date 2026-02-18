package main

import (
	"image/color"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a gate with specific time offsets from now
func createTestGate(name string, startOffsetHours int, endOffsetHours int) *Gate {
	now := time.Now()
	start := now.Add(time.Hour * time.Duration(startOffsetHours))
	end := now.Add(time.Hour * time.Duration(endOffsetHours))

	return &Gate{
		GateName: name,
		Start:    start.Format("2006-01-02 15:04:05"),
		End:      end.Format("2006-01-02 15:04:05"),
	}
}

func TestGetCellText(t *testing.T) {
	type GetCellTextTest struct {
		columnID     int
		gate         *Gate
		expectedText string
		desc         string
	}

	table := []GetCellTextTest{
		{
			columnID:     0,
			gate:         createTestGate("TestGate", -24, 24),
			expectedText: "TestGate",
			desc:         "Column 0 should return gate name",
		},
		{
			columnID:     1,
			gate:         createTestGate("PastGate", -72, -24),
			expectedText: "Past",
			desc:         "Column 1 with past gate should return 'Past'",
		},
		{
			columnID:     1,
			gate:         createTestGate("OpenGate", -12, 12),
			expectedText: "Open",
			desc:         "Column 1 with open gate should return 'Open'",
		},
		{
			columnID:     1,
			gate:         createTestGate("FutureGate", 24, 72),
			expectedText: "Future",
			desc:         "Column 1 with future gate should return 'Future'",
		},
		{
			columnID:     99,
			gate:         createTestGate("TestGate", -24, 24),
			expectedText: "",
			desc:         "Invalid column should return empty string",
		},
	}

	for _, v := range table {
		result := getCellText(v.columnID, v.gate)
		assert.Equal(t, v.expectedText, result, v.desc)
	}
}

func TestGetCellColor(t *testing.T) {
	type GetCellColorTest struct {
		gate             *Gate
		isTimelineLinear bool
		expectedColor    color.Color
		desc             string
	}

	table := []GetCellColorTest{
		{
			gate:             createTestGate("TestGate", -24, 24),
			isTimelineLinear: false,
			expectedColor:    color.RGBA{R: 255, G: 255, B: 0, A: 128},
			desc:             "Non-linear timeline should return yellow",
		},
		{
			gate:             createTestGate("OpenGate", -12, 12),
			isTimelineLinear: true,
			expectedColor:    color.RGBA{R: 0, G: 255, B: 0, A: 128},
			desc:             "Linear timeline with open gate should return green",
		},
		{
			gate:             createTestGate("ClosedGate", -72, -24),
			isTimelineLinear: true,
			expectedColor:    color.RGBA{R: 255, G: 0, B: 0, A: 128},
			desc:             "Linear timeline with closed gate should return red",
		},
		{
			gate:             createTestGate("FutureGate", 24, 72),
			isTimelineLinear: true,
			expectedColor:    color.RGBA{R: 255, G: 0, B: 0, A: 128},
			desc:             "Linear timeline with future gate should return red",
		},
	}

	for _, v := range table {
		result := getCellColor(v.gate, v.isTimelineLinear)
		assert.Equal(t, v.expectedColor, result, v.desc)
	}
}

func TestGetCellStyle(t *testing.T) {
	type GetCellStyleTest struct {
		columnID         int
		gate             *Gate
		isTimelineLinear bool
		expectedText     string
		expectedColor    color.Color
		desc             string
	}

	table := []GetCellStyleTest{
		{
			columnID:         0,
			gate:             createTestGate("IntegrationGate", -24, 24),
			isTimelineLinear: false,
			expectedText:     "IntegrationGate",
			expectedColor:    color.RGBA{R: 255, G: 255, B: 0, A: 128},
			desc:             "Column 0 with non-linear timeline",
		},
		{
			columnID:         1,
			gate:             createTestGate("OpenGate", -12, 12),
			isTimelineLinear: true,
			expectedText:     "Open",
			expectedColor:    color.RGBA{R: 0, G: 255, B: 0, A: 128},
			desc:             "Column 1 with open gate and linear timeline",
		},
		{
			columnID:         1,
			gate:             createTestGate("PastGate", -72, -24),
			isTimelineLinear: true,
			expectedText:     "Past",
			expectedColor:    color.RGBA{R: 255, G: 0, B: 0, A: 128},
			desc:             "Column 1 with past gate and linear timeline",
		},
		{
			columnID:         1,
			gate:             createTestGate("FutureGate", 24, 72),
			isTimelineLinear: true,
			expectedText:     "Future",
			expectedColor:    color.RGBA{R: 255, G: 0, B: 0, A: 128},
			desc:             "Column 1 with future gate and linear timeline",
		},
	}

	for _, v := range table {
		result := GetCellStyle(v.columnID, v.gate, v.isTimelineLinear)
		assert.Equal(t, v.expectedText, result.Text, v.desc+" - text")
		assert.Equal(t, v.expectedColor, result.Color, v.desc+" - color")
	}
}
