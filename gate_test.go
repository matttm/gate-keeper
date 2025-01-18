package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldCreateQueryString(t *testing.T) {
	type QueryStringTest struct {
		config        *GateConfig
		Date          string
		year          int
		pastGate      string
		magnitude     int
		expectedQuery string
	}
	config := &GateConfig{
		Dbname:      "DB",
		TableName:   "G",
		GateNameKey: "P",
		GateYearKey: "Y",
		StartKey:    "S",
		EndKey:      "E",
	}
	table := []QueryStringTest{
		{
			config:        config,
			Date:          "2025-06-15 12:00:00",
			year:          2026,
			pastGate:      "P2",
			magnitude:     0,
			expectedQuery: "UPDATE DB.G SET S = '2025-06-15 00:00:00', E = '2025-06-16 00:00:00' WHERE P = 'P2' AND Y = 2026;",
		},
		{
			config:        config,
			Date:          "2025-06-15 12:00:00",
			year:          2026,
			pastGate:      "P2",
			magnitude:     -1, // indicates this is the first prev gate
			expectedQuery: "UPDATE DB.G SET S = '2025-06-12 00:00:00', E = '2025-06-13 00:00:00' WHERE P = 'P2' AND Y = 2026;",
		},
	}
	for _, v := range table {
		d, _ := time.Parse(createdFormat, v.Date)
		q := _createQueryString(v.config, d, v.year, v.pastGate, v.magnitude)
		assert.Equal(t, q, v.expectedQuery)
	}
}
