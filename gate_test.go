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
