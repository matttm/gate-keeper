package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const year = "2026"

var config = &GateConfig{
	Dbname:      "DB",
	TableName:   "G",
	GateNameKey: "P",
	GateYearKey: "Y",
	StartKey:    "S",
	EndKey:      "E",
}

func TestGate_ShouldCreateQueryString(t *testing.T) {
	type QueryStringTest struct {
		config        *GateConfig
		Date          string
		year          int
		pastGate      string
		magnitude     int
		pos           RelativePosition
		expectedQuery string
	}
	table := []QueryStringTest{
		{
			config:        config,
			Date:          "2025-06-15 12:00:00",
			year:          2026,
			pastGate:      "P2",
			magnitude:     0,
			pos:           INSIDE,
			expectedQuery: "UPDATE DB.G SET S = '2025-06-15 00:00:00', E = '2025-06-16 00:00:00' WHERE P = 'P2' AND Y = 2026;",
		},
		{
			config:        config,
			Date:          "2025-06-15 12:00:00",
			year:          2026,
			pastGate:      "P2",
			magnitude:     -1, // indicates this is the first prev gate
			pos:           INSIDE,
			expectedQuery: "UPDATE DB.G SET S = '2025-06-12 00:00:00', E = '2025-06-13 00:00:00' WHERE P = 'P2' AND Y = 2026;",
		},
		{
			config:        config,
			Date:          "2025-06-15 12:00:00",
			year:          2026,
			pastGate:      "P2",
			magnitude:     1, // indicates this is the first future
			pos:           INSIDE,
			expectedQuery: "UPDATE DB.G SET S = '2025-06-18 00:00:00', E = '2025-06-19 00:00:00' WHERE P = 'P2' AND Y = 2026;",
		},
		{
			config:        config,
			Date:          "2025-06-15 12:00:00",
			year:          2026,
			pastGate:      "P2",
			magnitude:     2, // indicates this is the second future
			pos:           INSIDE,
			expectedQuery: "UPDATE DB.G SET S = '2025-06-21 00:00:00', E = '2025-06-22 00:00:00' WHERE P = 'P2' AND Y = 2026;",
		},
	}
	for _, v := range table {
		d, _ := time.Parse(createdFormat, v.Date)
		q := _createQueryString(v.config, d, v.year, v.pastGate, v.magnitude, v.pos)
		assert.Equal(t, q, v.expectedQuery)
	}
}
func TestGate_ShouldCreateQueryStrings(t *testing.T) {
	type QueryStringsTest struct {
		config          *GateConfig
		gates           []*Gate
		Date            string
		year            int
		pastGate        string
		pos             RelativePosition
		expectedQueries []string
		desc            string
	}
	gates := []*Gate{
		{
			GateName: "A",
			GateYear: year,
			Start:    "2020-01-01 00:00:00",
			End:      "2020-01-01 00:00:00",
		},
		{
			GateName: "B",
			GateYear: year,
			Start:    "2020-01-01 00:00:00",
			End:      "2020-01-01 00:00:00",
		},
		{
			GateName: "C",
			GateYear: year,
			Start:    "2020-01-01 00:00:00",
			End:      "2020-01-01 00:00:00",
		},
		{
			GateName: "D",
			GateYear: year,
			Start:    "2020-01-01 00:00:00",
			End:      "2020-01-01 00:00:00",
		},
	}
	table := []QueryStringsTest{
		{
			config:   config,
			gates:    gates,
			Date:     "2025-06-15 12:00:00", // this is time.Now()
			year:     2026,
			pastGate: "B",
			pos:      INSIDE,
			expectedQueries: []string{
				"UPDATE DB.G SET S = '2025-06-12 00:00:00', E = '2025-06-13 00:00:00' WHERE P = 'A' AND Y = 2026;",
				"UPDATE DB.G SET S = '2025-06-15 00:00:00', E = '2025-06-16 00:00:00' WHERE P = 'B' AND Y = 2026;",
				"UPDATE DB.G SET S = '2025-06-18 00:00:00', E = '2025-06-19 00:00:00' WHERE P = 'C' AND Y = 2026;",
				"UPDATE DB.G SET S = '2025-06-21 00:00:00', E = '2025-06-22 00:00:00' WHERE P = 'D' AND Y = 2026;",
			},
			desc: "",
		},
		{
			config:   config,
			gates:    gates,
			Date:     "2025-06-15 12:00:00", // this is time.Now()
			year:     2026,
			pastGate: "B",
			pos:      BEFORE,
			expectedQueries: []string{
				"UPDATE DB.G SET S = '2025-06-12 00:00:00', E = '2025-06-13 00:00:00' WHERE P = 'A' AND Y = 2026;",
				"UPDATE DB.G SET S = '2025-06-15 18:00:00', E = '2025-06-16 00:00:00' WHERE P = 'B' AND Y = 2026;",
				"UPDATE DB.G SET S = '2025-06-18 00:00:00', E = '2025-06-19 00:00:00' WHERE P = 'C' AND Y = 2026;",
				"UPDATE DB.G SET S = '2025-06-21 00:00:00', E = '2025-06-22 00:00:00' WHERE P = 'D' AND Y = 2026;",
			},
			desc: "",
		},
		{
			config:   config,
			gates:    gates,
			Date:     "2025-06-15 12:00:00", // this is time.Now()
			year:     2026,
			pastGate: "B",
			pos:      AFTER,
			expectedQueries: []string{
				"UPDATE DB.G SET S = '2025-06-12 00:00:00', E = '2025-06-13 00:00:00' WHERE P = 'A' AND Y = 2026;",
				"UPDATE DB.G SET S = '2025-06-15 00:00:00', E = '2025-06-15 06:00:00' WHERE P = 'B' AND Y = 2026;",
				"UPDATE DB.G SET S = '2025-06-18 00:00:00', E = '2025-06-19 00:00:00' WHERE P = 'C' AND Y = 2026;",
				"UPDATE DB.G SET S = '2025-06-21 00:00:00', E = '2025-06-22 00:00:00' WHERE P = 'D' AND Y = 2026;",
			},
			desc: "",
		},
	}
	for _, v := range table {
		d, _ := time.Parse(createdFormat, v.Date)
		q := _setGatesRelativeTo(v.config, v.gates, d, v.year, v.pastGate, v.pos)
		assert.Equal(t, q, v.expectedQueries)
	}
}

func TestIsTimelineLinear(t *testing.T) {
	type IsTimelineLinearTest struct {
		gates    []*Gate
		expected bool
		desc     string
	}

	table := []IsTimelineLinearTest{
		{
			gates: []*Gate{
				{GateName: "A", Start: "2025-01-01 00:00:00", End: "2025-01-05 00:00:00"},
				{GateName: "B", Start: "2025-01-06 00:00:00", End: "2025-01-10 00:00:00"},
				{GateName: "C", Start: "2025-01-11 00:00:00", End: "2025-01-15 00:00:00"},
			},
			expected: true,
			desc:     "Linear timeline with no overlaps",
		},
		{
			gates: []*Gate{
				{GateName: "A", Start: "2025-01-01 00:00:00", End: "2025-01-05 00:00:00"},
				{GateName: "B", Start: "2025-01-05 00:00:00", End: "2025-01-10 00:00:00"},
				{GateName: "C", Start: "2025-01-11 00:00:00", End: "2025-01-15 00:00:00"},
			},
			expected: false,
			desc:     "Non-linear timeline with gates touching (B starts when A ends)",
		},
		{
			gates: []*Gate{
				{GateName: "A", Start: "2025-01-01 00:00:00", End: "2025-01-10 00:00:00"},
				{GateName: "B", Start: "2025-01-05 00:00:00", End: "2025-01-15 00:00:00"},
				{GateName: "C", Start: "2025-01-16 00:00:00", End: "2025-01-20 00:00:00"},
			},
			expected: false,
			desc:     "Non-linear timeline with overlapping gates",
		},
		{
			gates: []*Gate{
				{GateName: "A", Start: "2025-01-01 00:00:00", End: "2025-01-05 00:00:00"},
				{GateName: "B", Start: "2025-01-10 00:00:00", End: "2025-01-15 00:00:00"},
				{GateName: "C", Start: "2025-01-06 00:00:00", End: "2025-01-09 00:00:00"},
			},
			expected: false,
			desc:     "Non-linear timeline with out-of-order gates",
		},
		{
			gates: []*Gate{
				{GateName: "A", Start: "2025-01-01 00:00:00", End: "2025-01-05 00:00:00"},
			},
			expected: true,
			desc:     "Single gate is linear",
		},
		{
			gates:    []*Gate{},
			expected: true,
			desc:     "Empty gates list is linear",
		},
	}

	for _, v := range table {
		result := isTimelineLinear(v.gates)
		assert.Equal(t, v.expected, result, v.desc)
	}
}

func TestIsGateOpen(t *testing.T) {
	type IsGateOpenTest struct {
		gate     *Gate
		expected bool
		desc     string
	}

	now := time.Now()
	table := []IsGateOpenTest{
		{
			gate: &Gate{
				GateName: "OpenGate",
				Start:    now.Add(-12 * time.Hour).Format(createdFormat),
				End:      now.Add(12 * time.Hour).Format(createdFormat),
			},
			expected: true,
			desc:     "Gate currently open (now is between start and end)",
		},
		{
			gate: &Gate{
				GateName: "PastGate",
				Start:    now.Add(-72 * time.Hour).Format(createdFormat),
				End:      now.Add(-24 * time.Hour).Format(createdFormat),
			},
			expected: false,
			desc:     "Gate in the past (ended before now)",
		},
		{
			gate: &Gate{
				GateName: "FutureGate",
				Start:    now.Add(24 * time.Hour).Format(createdFormat),
				End:      now.Add(72 * time.Hour).Format(createdFormat),
			},
			expected: false,
			desc:     "Gate in the future (starts after now)",
		},
		{
			gate: &Gate{
				GateName: "JustStarted",
				Start:    now.Add(-2 * time.Hour).Format(createdFormat),
				End:      now.Add(24 * time.Hour).Format(createdFormat),
			},
			expected: true,
			desc:     "Gate that started recently is open",
		},
	}

	for _, v := range table {
		result := isGateOpen(v.gate)
		assert.Equal(t, v.expected, result, v.desc)
	}
}

func TestGetGatePosition(t *testing.T) {
	type GetGatePositionTest struct {
		gate     *Gate
		expected int
		desc     string
	}

	now := time.Now()
	table := []GetGatePositionTest{
		{
			gate: &Gate{
				GateName: "OpenGate",
				Start:    now.Add(-12 * time.Hour).Format(createdFormat),
				End:      now.Add(12 * time.Hour).Format(createdFormat),
			},
			expected: 0,
			desc:     "Gate currently open returns 0",
		},
		{
			gate: &Gate{
				GateName: "PastGate",
				Start:    now.Add(-72 * time.Hour).Format(createdFormat),
				End:      now.Add(-24 * time.Hour).Format(createdFormat),
			},
			expected: -1,
			desc:     "Gate in the past returns -1",
		},
		{
			gate: &Gate{
				GateName: "FutureGate",
				Start:    now.Add(24 * time.Hour).Format(createdFormat),
				End:      now.Add(72 * time.Hour).Format(createdFormat),
			},
			expected: 1,
			desc:     "Gate in the future returns 1",
		},
	}

	for _, v := range table {
		result := getGatePosition(v.gate)
		assert.Equal(t, v.expected, result, v.desc)
	}
}
